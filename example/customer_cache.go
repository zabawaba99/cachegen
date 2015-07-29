// Do not modify this file.
// This file was automatically generate using github.com/zabawaba99/cachegen.

package main

import (
	"runtime"
	"sync"
	"time"
)

type customerWrapper struct {
	v         Customer
	expiredAt time.Time
}

func (w customerWrapper) isExpired() bool {
	return time.Now().After(w.expiredAt)
}

type customerCache struct {
	mtx sync.RWMutex
	m   map[int]*customerWrapper

	ttl         time.Duration
	cleanupTime time.Duration
	stop        chan struct{}
}

type CustomerCache struct {
	*customerCache
}

func NewCustomerCache(ttl, cleanupTime time.Duration) *CustomerCache {
	if cleanupTime == 0 {
		cleanupTime = 5 * time.Second
	}

	c := &CustomerCache{
		customerCache: &customerCache{
			m:           map[int]*customerWrapper{},
			ttl:         ttl,
			cleanupTime: cleanupTime,
			stop:        make(chan struct{}),
		},
	}

	go c.cleanup()
	runtime.SetFinalizer(c, stopCustomerCacheCleanup)

	return c
}

func (c *customerCache) cleanup() {
	for {
		select {
		case <-time.After(c.cleanupTime):
			c.deleteExpired()
		case <-c.stop:
			return
		}
	}
}

func (c *customerCache) Add(k int, v Customer) {
	c.mtx.Lock()
	c.m[k] = &customerWrapper{
		v:         v,
		expiredAt: time.Now().Add(c.ttl),
	}
	c.mtx.Unlock()
}

func (c *customerCache) deleteExpired() {
	c.mtx.Lock()
	for k, v := range c.m {
		if v.isExpired() {
			delete(c.m, k)
		}
	}
	c.mtx.Unlock()
}

func (c *customerCache) Get(k int) (v Customer, ok bool) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	wrapper, ok := c.m[k]
	if !ok || wrapper.isExpired() {
		return v, false
	}

	return wrapper.v, true
}

func (c *customerCache) Expire(k int) {
	c.mtx.RLock()
	wrapper, ok := c.m[k]
	if ok {
		wrapper.expiredAt = time.Now()
	}
	c.mtx.RUnlock()
}

func stopCustomerCacheCleanup(c *CustomerCache) {
	close(c.stop)
}
