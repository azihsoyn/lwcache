package lwcache

import (
	"sync"
	"time"
)

const (
	NoExpire = time.Duration(0)
)

const (
	refresherOn  = true
	refresherOff = false
)

type refreshNotifier struct {
	c   *cache
	key interface{}
}

type Cache interface {
	Get(key interface{}) (value interface{}, ok bool)
	Set(key interface{}, item interface{}, expire time.Duration)
	Del(key interface{})
	SetExpire(key interface{}, expire time.Duration)
	SetRefresher(fn func(c Cache, key interface{}, currentValue interface{}) (newValue interface{}, err error))
	StartRefresher(key interface{}, refreshInterval time.Duration)
	StopRefresher(key interface{})
}

type cache struct {
	namespace string
	items     map[interface{}]cacheItem
	mutex     sync.RWMutex
	refresher func(c Cache, key interface{}, currentValue interface{}) (newValue interface{}, err error)
	// for refresher
	refresherStatus map[interface{}]bool
	refresherMutex  sync.RWMutex
}

var _ Cache = (*cache)(nil)

type cacheItem struct {
	expiration time.Time
	value      interface{}
	deleter    *time.Timer
}

func New(name string) Cache {
	c := &cache{
		namespace:       name,
		items:           make(map[interface{}]cacheItem),
		mutex:           sync.RWMutex{},
		refresherStatus: make(map[interface{}]bool),
		refresherMutex:  sync.RWMutex{},
	}
	return c
}

// TODO: implement MSet (multi set) for performance
func (c *cache) Set(key, value interface{}, expire time.Duration) {
	var (
		expiration time.Time
		deleter    *time.Timer
	)
	if expire != NoExpire {
		expiration = time.Now().Add(expire)
		deleter = time.AfterFunc(expire, func() {
			c.mutexDelete(key)
		})
	}
	item := cacheItem{
		value:      value,
		expiration: expiration,
		deleter:    deleter,
	}

	c.mutexSet(key, item)
}

func (c *cache) Del(key interface{}) {
	c.mutexDelete(key)
}

func (c *cache) SetExpire(key interface{}, expire time.Duration) {
	if item, ok := c.mutexGet(key); ok {
		if expire == NoExpire {
			item.expiration = time.Time{}
			if item.deleter != nil {
				item.deleter.Stop()
				item.deleter = nil
			}
		} else {
			item.expiration = time.Now().Add(expire)
			if item.deleter != nil {
				item.deleter.Reset(expire)
			} else {
				item.deleter = time.AfterFunc(expire, func() {
					c.mutexDelete(key)
				})
			}
		}
		c.mutexSet(key, item)
	}
}

func (c *cache) SetRefresher(fn func(c Cache, key interface{}, currentValue interface{}) (newValue interface{}, err error)) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.refresher = fn
}

func (c *cache) StartRefresher(key interface{}, refreshInterval time.Duration) {
	c.refresherMutex.RLock()
	running := c.refresherStatus[key]
	c.refresherMutex.RUnlock()
	if !running {
		c.refresherMutex.Lock()
		c.refresherStatus[key] = refresherOn
		c.refresherMutex.Unlock()
		time.AfterFunc(refreshInterval, c.startBackGroundRefresh(key, refreshInterval))
	}
}

func (c *cache) StopRefresher(key interface{}) {
	c.refresherMutex.Lock()
	c.refresherStatus[key] = refresherOff
	c.refresherMutex.Unlock()
}

func (c *cache) startBackGroundRefresh(key interface{}, interval time.Duration) func() {
	return func() {
		c.refresherMutex.RLock()
		running := c.refresherStatus[key]
		c.refresherMutex.RUnlock()
		if !running {
			return
		}
		c.refresh(key)
		time.AfterFunc(interval, c.startBackGroundRefresh(key, interval))
	}
}

/*
func (c *cache) OnExpire(key string, fn func(item interface{}) error) error {
	return fn(c.item)
}

func (c *cache) OnRefresh(key string, fn func(item interface{}) error) error {
	return fn(c.item)
}
*/

// TODO: implement MGet (multi get) for performance
func (c *cache) Get(key interface{}) (interface{}, bool) {
	item, ok := c.mutexGet(key)

	// no expire on expiration is zero value
	if !ok {
		return nil, false
	}

	return item.value, ok
}

func (c *cache) mutexGet(key interface{}) (cacheItem, bool) {
	c.mutex.RLock()
	item, ok := c.items[key]
	c.mutex.RUnlock()

	return item, ok
}

func (c *cache) mutexDelete(key interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// for avoid goroutine leak
	if c.items[key].deleter != nil {
		c.items[key].deleter.Stop()
	}

	delete(c.items, key)
}

func (c *cache) mutexSet(key interface{}, item cacheItem) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = item
}

func (c *cache) mutexSetValue(key, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item := c.items[key]
	item.value = value
	c.items[key] = item

}

func (c *cache) refresh(key interface{}) {
	// avoid panic
	// see: https://github.com/azihsoyn/lwcache/issues/5
	if c.refresher == nil {
		return
	}

	item, ok := c.mutexGet(key)
	if !ok {
		// item deleted. not refresh
		return
	}
	val, err := c.refresher(c, key, item.value)
	if err != nil {
		// TODO?: notifyError
		return
	}
	c.mutexSetValue(key, val)
}
