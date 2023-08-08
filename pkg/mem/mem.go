package mem

import (
	"fmt"
	"time"
)

const (
	NoExpiration      time.Duration = -1
	DefaultExpiration time.Duration = 0
)

type Item struct {
	Object     interface{}
	Expiration int64
}

type Cache struct {
	*cache
}

type cache struct {
	items             map[string]Item
	defaultExpiration time.Duration
}

func (c *cache) Get(key string) (interface{}, bool) {
	item, found := c.items[key]

	if !found {
		return nil, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}

	return item.Object, true
}

func (c *cache) Set(key string, val interface{}, timeout time.Duration) {
	var expiration int64

	if timeout > 0 {
		expiration = time.Now().Add(timeout).UnixNano()
	}

	c.items[key] = Item{
		Object:     val,
		Expiration: expiration,
	}
}

func (c *cache) Add(key string, val interface{}, timeout time.Duration) error {
	_, found := c.Get(key)

	if found {
		return fmt.Errorf("Item %s already exists", key)
	}

	c.Set(key, val, timeout)
	return nil
}

func (c *cache) Delete(key string) {
	delete(c.items, key)
}

func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	items := make(map[string]Item)
	c := &cache{
		defaultExpiration: defaultExpiration,
		items:             items,
	}
	return &Cache{c}
}
