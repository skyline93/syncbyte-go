package cache

import (
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

type item struct {
	object     interface{}
	expiration int64
}

func (i *item) expired() bool {
	if i.expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > i.expiration
}

type LRU struct {
	evictList *List
	size      int
}

func NewLRU(size int) *LRU {
	return &LRU{
		evictList: NewList(),
		size:      size,
	}
}

func (l *LRU) MoveToFront(key interface{}) {
	l.evictList.Delete(key)
	l.evictList.Insert(key)
}

func (l *LRU) removeOldest() (key interface{}) {
	return l.evictList.DeleteTail()
}

func (l *LRU) Insert(key interface{}) (evicted_data interface{}, evicted bool) {
	if l.evictList.Len >= int64(l.size) {
		evicted_data = l.removeOldest()
		evicted = true
	}

	l.evictList.Insert(key)
	return
}

func (l *LRU) Delete(key interface{}) {
	l.evictList.Delete(key)
}

func (l *LRU) Contains(key interface{}) bool {
	return l.evictList.Contains(key)
}

var DefaultDuration time.Duration = 0

type LRUCache struct {
	items map[interface{}]item
	mu    sync.RWMutex

	lru             *LRU
	defaultDuration time.Duration
}

func NewLRUCache(size int) *LRUCache {
	return &LRUCache{
		items:           make(map[interface{}]item, size),
		lru:             NewLRU(size),
		defaultDuration: DefaultDuration,
	}
}

func (c *LRUCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

func (c *LRUCache) Set(key interface{}, value interface{}, duration time.Duration) {
	c.set(key, value, duration)
}

func (c *LRUCache) SetWithUuidKey(value interface{}, duration time.Duration) (key string) {
	key = uuid.NewV4().String()
	c.set(key, value, duration)

	return key
}

func (c *LRUCache) SetDefault(key interface{}, value interface{}) {
	c.set(key, value, c.defaultDuration)
}

func (c *LRUCache) SetDefaultWithUuidKey(value interface{}) (key string) {
	key = uuid.NewV4().String()
	c.set(key, value, c.defaultDuration)

	return key
}

func (c *LRUCache) set(key interface{}, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var exp int64

	if duration == 0 {
		exp = 0
	} else {
		exp = time.Now().Add(duration).UnixNano()
	}

	if c.lru.Contains(key) {
		c.lru.MoveToFront(key)
	} else {
		evicted_data, is_evicted := c.lru.Insert(key)
		if is_evicted {
			delete(c.items, evicted_data)
		}
	}

	c.items[key] = item{object: value, expiration: exp}
}

func (c *LRUCache) Get(key interface{}) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		return nil
	}

	if item.expired() {
		return nil
	}

	c.lru.MoveToFront(key)

	return item.object
}

func (c *LRUCache) delete(key interface{}) {
	c.lru.Delete(key)
	delete(c.items, key)
}

func (c *LRUCache) Delete(key interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.delete(key)
}

func (c *LRUCache) DeleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, i := range c.items {
		if i.expired() {
			c.delete(k)
		}
	}
}

func (c *LRUCache) Range() []interface{} {
	var ls []interface{}

	for _, i := range c.items {
		ls = append(ls, i.object)
	}

	return ls
}

type gc struct {
	Interval time.Duration
	stop     chan interface{}
}

func (g *gc) Run(c *LRUCache) {
	ticker := time.NewTicker(g.Interval)

	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-g.stop:
			ticker.Stop()
			return
		}
	}
}

func (g *gc) Stop() {
	g.stop <- struct{}{}
}

func newLRUCacheWithGC(size int, cleanupInterval time.Duration) *LRUCache {
	c := NewLRUCache(size)

	g := gc{
		Interval: cleanupInterval,
		stop:     make(chan interface{}),
	}

	go g.Run(c)

	return c
}

type Cache struct {
	*LRUCache
}

func New(size int, cleanupInterval time.Duration) *Cache {
	c := newLRUCacheWithGC(size, cleanupInterval)
	return &Cache{c}
}
