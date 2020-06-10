package exporter

import (
	"fmt"

	"github.com/patrickmn/go-cache"
)

// PNCache Self implemented interface
type PNCache struct {
	C *cache.Cache
}

// Init Initialize the in-memory cache instance
func Initcache() PNCache {
	c := cache.New(0, 0)
	myCache := PNCache{c}
	return myCache
}

// AddKeyValueCache Add or update a cache value with the given key
func (c *PNCache) AddKeyValueCache(key string, value interface{}) {
	// Create a cache without any default expiration time and cleanup time
	ssl := value.(SSLInfoArray)
	c.C.Set(key, ssl, cache.NoExpiration)

	value, found := c.C.Get(key)

	if found {
		fmt.Printf("Insert new key-value pair %s:%v\n", key, value)
	} else {
		fmt.Println("value is not found.")
	}
}

// GetValue Retrieve the value from the given key
func (c *PNCache) GetValue(key string) interface{} {
	value, found := c.C.Get(key)
	if found {
		fmt.Println(value)
	} else {
		fmt.Println("value is not found.")
		return nil
	}
	return value
}
