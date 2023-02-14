package lru

import "errors"

type Cacher interface {
	Get(interface{}) (interface{}, error)
	Put(interface{}, interface{}) error
}

type lruCache struct {
	size      int
	remaining int
	cache     map[string]string
	queue     []string
}

func NewCache(size int) Cacher {
	return &lruCache{size: size, remaining: size, cache: make(map[string]string), queue: make([]string, 0)}
}

func (lru *lruCache) Get(key interface{}) (interface{}, error) {
	// Your code here....

	_, keyInMap := lru.cache[key.(string)] // Checks whether key is in map or not

	if keyInMap == false {

		return key, errors.New("Key is not in map") // If key is not in map, output read error

	} else {

		lru.queue = append(lru.queue, key.(string)) // Else append to queue

	}

	return lru.cache[key.(string)], nil // Returns mapped value and error value
}

func (lru *lruCache) Put(key, val interface{}) error {
	// Your code here....

	_, keyInMap := lru.cache[key.(string)] // Checks whether key is in map or not

	if lru.remaining == 0 {

		tempDel := lru.queue[0]    // Temp variable to hold queue value
		delete(lru.cache, tempDel) // Delete value out of cache
		lru.qDel(tempDel)          // Delete value out of queue

		lru.remaining++ // Update remaining

	} else {

		if keyInMap == true {
			lru.qDel(key.(string)) // Delete value out of queue
		}

		lru.queue = append(lru.queue, key.(string)) // Add value to queue
		lru.cache[key.(string)] = val.(string)      // Add mapping association

		lru.remaining-- // Update remaining

	}

	return nil
}

// Delete element from queue
func (lru *lruCache) qDel(ele string) {
	for i := 0; i < len(lru.queue); i++ {
		if lru.queue[i] == ele {
			oldlen := len(lru.queue)
			copy(lru.queue[i:], lru.queue[i+1:])
			lru.queue = lru.queue[:oldlen-1]
			break
		}
	}
}
