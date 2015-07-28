package main

import (
	"runtime"
	"sync"
	"time"
)

type intWrapper struct {
	v         int
	expiredAt time.Time
}

func (w intWrapper) isExpired() bool {
	return time.Now().After(w.expiredAt)
}

type IntCache struct {
	mtx sync.RWMutex
	m   map[string]*intWrapper

	ttl         time.Duration
	cleanupTime time.Duration
	stop        chan struct{}
}

func NewACache(ttl, cleanupTime time.Duration) *IntCache {
	if cleanupTime == 0 {
		cleanupTime = 5 * time.Second
	}

	c := &IntCache{
		m:           map[string]*intWrapper{},
		ttl:         ttl,
		cleanupTime: cleanupTime,
		stop:        make(chan struct{}),
	}

	go c.cleanup()
	runtime.SetFinalizer(c, stopCleanup)

	return c
}

func (c *IntCache) cleanup() {
	for {
		select {
		case <-time.After(c.cleanupTime):
			c.deleteExpired()
		case <-c.stop:
			return
		}
	}
}

func (c *IntCache) Add(k string, v int) {
	c.mtx.Lock()
	c.m[k] = &intWrapper{
		v:         v,
		expiredAt: time.Now().Add(c.ttl),
	}
	c.mtx.Unlock()
}

func (c *IntCache) deleteExpired() {
	c.mtx.Lock()
	for k, v := range c.m {
		if v.isExpired() {
			delete(c.m, k)
		}
	}
	c.mtx.Unlock()
}

func (c *IntCache) Get(k string) (v int, ok bool) {
	c.mtx.RLock()
	wrapper, ok := c.m[k]
	c.mtx.RUnlock()
	if !ok || wrapper.isExpired() {
		return v, false
	}

	return wrapper.v, true
}

func (c *IntCache) Expire(k string) {
	c.mtx.RLock()
	wrapper, ok := c.m[k]
	c.mtx.RUnlock()
	if ok {
		wrapper.expiredAt = time.Now()
	}
}

func stopCleanup(c *IntCache) {
	close(c.stop)
}