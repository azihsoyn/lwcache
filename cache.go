package cache

import (
	"sync"
	"time"
)

var expChan = make(chan expireNotifier, 1000) // TODO: configurable

type Getter interface {
	Get(key interface{}) (interface{}, bool)
}

type expireNotifier struct {
	key    interface{}
	expire time.Time
}

type cache struct {
	items map[interface{}]cacheItem
	mutex sync.RWMutex
}

type cacheItem struct {
	expire *time.Timer
	value  interface{}
}

func New(name string) *cache {
	c := &cache{
		items: make(map[interface{}]cacheItem),
		mutex: sync.RWMutex{},
	}
	go func() {
		for {
			select {
			case exp := <-expChan:
				c.mutex.Lock()
				delete(c.items, exp.key)
				c.mutex.Unlock()
			default:
			}
		}
	}()
	return c
}

func (c *cache) Set(key interface{}, item interface{}, expire time.Duration) error {
	timer := time.NewTimer(expire)
	c.mutex.Lock()
	c.items[key] = cacheItem{
		value:  item,
		expire: timer,
	}
	c.mutex.Unlock()

	go notifyExpire(key, timer)
	return nil
}

func (c *cache) SetExpire(key interface{}, expire time.Duration) {
	c.mutex.RLock()
	item := c.items[key]
	timer := item.expire
	c.mutex.RUnlock()

	timer.Reset(expire)
}

func notifyExpire(key interface{}, timer *time.Timer) {
	expiredAt := <-timer.C
	expChan <- expireNotifier{
		key:    key,
		expire: expiredAt,
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

// TODO?: check expire if need more accuracy
func (c *cache) Get(key interface{}) (interface{}, bool) {
	c.mutex.RLock()
	v, ok := c.items[key]
	c.mutex.RUnlock()

	return v.value, ok
}
