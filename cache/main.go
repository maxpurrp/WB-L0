package cache

import (
	"client/order"
)

type Cache struct {
	store map[string]*order.Order
}

func NewCache() *Cache {
	return &Cache{
		store: make(map[string]*order.Order),
	}
}

func (c *Cache) Set(key string, value *order.Order) {
	c.store[key] = value
}

func (c *Cache) Get(key string) (*order.Order, bool) {
	val, found := c.store[key]
	return val, found
}

func SaveInMemory(cache *Cache, data *order.Order) {
	cache.Set(data.OrderUID, data)
}

func GetValue(cache *Cache, key string) *order.Order {
	value, found := cache.Get(key)
	if found {
		return value
	} else {
		return nil
	}
}
