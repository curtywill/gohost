package requests

import (
	"time"
)

type c struct {
	data    []byte
	timeout time.Time
}

// cache is indexed by cookie, and then indexed by endpoint
type cacheEntry map[string]c

var cache = make(map[string]cacheEntry)

func addToCache(cookie, endpoint string, data []byte) {
	if cache[cookie] == nil {
		cache[cookie] = make(cacheEntry)
	}
	cache[cookie][endpoint] = c{data, time.Now()}
}

func getFromCache(cookie, endpoint string) []byte {
	cachedData, ok := cache[cookie][endpoint]

	if !ok {
		return nil
	}
	t := time.Now()
	if t.Sub(cachedData.timeout).Seconds() > 30 {
		return nil
	}
	return cachedData.data
}
