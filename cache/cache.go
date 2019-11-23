package cache

import (
	"context"
	"sync"
	"time"
)

const defaultFlushDuration = 10 * time.Second

type Cache struct {
	sync.RWMutex
	values        map[string]*cacheItem
	flushDuration time.Duration
	cancelFunc    context.CancelFunc
}

type cacheItem struct {
	value   interface{}
	setTime time.Time
	dataTtl time.Duration
}

func NewCache() *Cache {
	ctx, cancelFunc := context.WithCancel(context.Background())

	c := &Cache{
		values:        make(map[string]*cacheItem),
		flushDuration: defaultFlushDuration,
		cancelFunc:    cancelFunc,
	}

	go c.flushData(ctx)

	return c
}

func (c *Cache)GetTTL(key string) time.Duration {
	val := c.values[key]
	ttl := val.dataTtl - (time.Now().Sub(val.setTime))
	if ttl < 0 {
		return 0
	}

	return ttl
}

func (c *Cache)StopCache() {
	c.cancelFunc()
}

func (c *Cache) Set(key string, data interface{}) {
	c.values[key] = &cacheItem{
		value:   data,
		setTime: time.Now(),
		dataTtl: 0,
	}
}

func (c *Cache) SetTTL(key string, data interface{}, ttl time.Duration) {
	c.values[key] = &cacheItem{
		value:   data,
		setTime: time.Now(),
		dataTtl: ttl,
	}
}

func (c *Cache) Get(key string) interface{} {
	i, ok := c.values[key]
	if !ok {
		return nil
	}

	if i.dataTtl == 0 {
		return i
	}

	if time.Now().After(i.setTime.Add(i.dataTtl)) {
		return nil
	}

	return i.value
}

func (c *Cache) GetInt(key string) int {
	i := c.Get(key)
	if i == nil {
		return 0
	}

	num, ok := i.(int)
	if !ok {
		return 0
	}

	return num
}

func (c *Cache) GetString(key string) string {
	i := c.Get(key)
	if i == nil {
		return ""
	}

	str, ok := i.(string)
	if !ok {
		return ""
	}

	return str
}

func (c *Cache) GetBool(key string) bool {
	i := c.Get(key)
	if i == nil {
		return false
	}

	b, ok := i.(bool)
	if !ok {
		return false
	}

	return b
}

func (c *Cache) SetFlushDuration(d time.Duration) {
	c.flushDuration = d
}

func (c *Cache) flushData(ctx context.Context) {
	for {
		select {
		case <- ctx.Done():
			return
		default:
			for key, val := range c.values {
				if time.Now().After(val.setTime.Add(val.dataTtl)) {
					delete(c.values, key)
				}
			}
		}

		time.Sleep(c.flushDuration)
	}
}