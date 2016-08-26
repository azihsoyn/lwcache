package cache

import (
	"sync"
	"time"
)

type Getter interface {
	Get(key interface{}) (value interface{}, ok bool)
}

type expireNotifier struct {
	key    interface{}
	expire time.Time
}

type refreshNotifier struct {
	key interface{}
}

type Cache interface {
	Getter
	Set(key interface{}, item interface{}, expire time.Duration)
	SetExpire(key interface{}, expire time.Duration)
	SetRefresher(key interface{}, fn func(key interface{}, currentValue interface{}) (newValue interface{}, err error), refreshInterval time.Duration)
}

type cache struct {
	items   map[interface{}]cacheItem
	mutex   sync.RWMutex
	expChan chan expireNotifier
	refChan chan refreshNotifier
}

var _ Cache = (*cache)(nil)

type cacheItem struct {
	expire          *time.Timer
	value           interface{}
	refreshInterval time.Duration
	refresher       func(key interface{}, currentValue interface{}) (newValue interface{}, err error)
}

func New(name string) Cache {
	c := &cache{
		items:   make(map[interface{}]cacheItem),
		mutex:   sync.RWMutex{},
		expChan: make(chan expireNotifier, 128),  // TODO: configurable
		refChan: make(chan refreshNotifier, 128), // TODO: configurable
	}
	go func(c *cache) {
		for {
			select {
			case exp := <-c.expChan:
				c.mutex.Lock()
				delete(c.items, exp.key)
				c.mutex.Unlock()
			case ref := <-c.refChan:
				// TODO: define func
				c.mutex.Lock()
				item, ok := c.items[ref.key]
				if !ok {
					// item deleted. stop refresher
					c.mutex.Unlock()
					continue
				}
				val, err := item.refresher(ref.key, item.value)
				if err != nil {
					// TODO?: notifyError
					c.mutex.Unlock()
					continue
				}
				item.value = val
				c.items[ref.key] = item
				c.mutex.Unlock()
			default:
			}
		}
	}(c)
	return c
}

// TODO: implement MSet (multi set) for performance
func (c *cache) Set(key interface{}, item interface{}, expire time.Duration) {
	timer := time.NewTimer(expire)
	c.mutex.Lock()
	c.items[key] = cacheItem{
		value:  item,
		expire: timer,
	}
	c.mutex.Unlock()

	go c.notifyExpire(key, timer)
}

func (c *cache) SetExpire(key interface{}, expire time.Duration) {
	c.mutex.Lock()
	item := c.items[key]
	timer := item.expire
	c.mutex.Unlock()

	timer.Reset(expire)
}

func (c *cache) SetRefresher(key interface{}, fn func(key interface{}, currentValue interface{}) (newValue interface{}, err error), refreshInterval time.Duration) {
	c.mutex.Lock()
	item := c.items[key]
	item.refresher = fn
	item.refreshInterval = refreshInterval
	c.items[key] = item
	c.mutex.Unlock()

	go c.startBackGroundRefresh(key, refreshInterval)
}

func (c *cache) notifyExpire(key interface{}, timer *time.Timer) {
	expiredAt := <-timer.C
	c.expChan <- expireNotifier{
		key:    key,
		expire: expiredAt,
	}
}

func (c *cache) startBackGroundRefresh(key interface{}, interval time.Duration) {
	timer := time.NewTimer(interval)
	<-timer.C
	c.refChan <- refreshNotifier{
		key: key,
	}
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	// item exists
	if _, ok := c.items[key]; ok {
		go c.startBackGroundRefresh(key, interval)
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

// TODO: implement MGet (multi set) for performance
// TODO?: check expire if need more accuracy
func (c *cache) Get(key interface{}) (interface{}, bool) {
	c.mutex.RLock()
	v, ok := c.items[key]
	c.mutex.RUnlock()

	return v.value, ok
}
