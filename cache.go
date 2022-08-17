package main

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value     interface{}
	Timestamp time.Time
}

type Cache struct {
	sync.RWMutex
	Items    map[string]*CacheItem
	TickerGC *time.Ticker
	StopGC   chan bool
	TTL      time.Duration
}

func (c *Cache) Set(key string, value interface{}) {
	c.Lock()
	defer c.Unlock()

	c.Items[key] = &CacheItem{
		Value:     value,
		Timestamp: time.Now(),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	item, found := c.Items[key]
	if !found {
		return nil, false
	}

	return item.Value, true
}

func (c *Cache) Clean() {
	c.Lock()
	defer c.Unlock()

	now := time.Now()
	for key, item := range c.Items {
		if item.Timestamp.Add(c.TTL).Before(now) {
			delete(c.Items, key)
		}
	}
}

func (c *Cache) StartGC() {
	go func() {
		for {
			select {
			case <-c.TickerGC.C:
				c.Clean()
			case <-c.StopGC:
				return
			}
		}
	}()
}

func NewCache(ttl time.Duration, gcInterval time.Duration) *Cache {
	c := &Cache{
		Items:    make(map[string]*CacheItem),
		TickerGC: time.NewTicker(gcInterval),
		TTL:      ttl,
		StopGC:   make(chan bool),
	}
	c.StartGC()

	return c
}
