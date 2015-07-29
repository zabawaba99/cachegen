// Do not modify this file.
// This file was automatically generate using github.com/zabawaba99/cachegen.

package cache

import (
	"runtime"
	"sync"
	"time"
)

type ReplaceKey string
type ReplaceValue string

type aWrapper struct {
	v         ReplaceValue
	expiredAt time.Time
}

func (w aWrapper) isExpired() bool {
	return time.Now().After(w.expiredAt)
}

type aCache struct {
	mtx sync.RWMutex
	m   map[ReplaceKey]*aWrapper

	ttl         time.Duration
	cleanupTime time.Duration
	stop        chan struct{}
}

type ACache struct {
	*aCache
}

func NewACache(ttl, cleanupTime time.Duration) *ACache {
	if cleanupTime == 0 {
		cleanupTime = 5 * time.Second
	}

	c := &ACache{
		aCache: &aCache{
			m:           map[ReplaceKey]*aWrapper{},
			ttl:         ttl,
			cleanupTime: cleanupTime,
			stop:        make(chan struct{}),
		},
	}

	go c.cleanup()
	runtime.SetFinalizer(c, stopACacheCleanup)

	return c
}

func (c *aCache) cleanup() {
	for {
		select {
		case <-time.After(c.cleanupTime):
			c.deleteExpired()
		case <-c.stop:
			return
		}
	}
}

func (c *aCache) Add(k ReplaceKey, v ReplaceValue) {
	c.mtx.Lock()
	c.m[k] = &aWrapper{
		v:         v,
		expiredAt: time.Now().Add(c.ttl),
	}
	c.mtx.Unlock()
}

func (c *aCache) deleteExpired() {
	c.mtx.Lock()
	for k, v := range c.m {
		if v.isExpired() {
			delete(c.m, k)
		}
	}
	c.mtx.Unlock()
}

func (c *aCache) Get(k ReplaceKey) (v ReplaceValue, ok bool) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	wrapper, ok := c.m[k]
	if !ok || wrapper.isExpired() {
		return v, false
	}

	return wrapper.v, true
}

func (c *aCache) Expire(k ReplaceKey) {
	c.mtx.RLock()
	wrapper, ok := c.m[k]
	if ok {
		wrapper.expiredAt = time.Now()
	}
	c.mtx.RUnlock()
}

func stopACacheCleanup(c *ACache) {
	close(c.stop)
}
