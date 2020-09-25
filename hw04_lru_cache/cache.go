package hw04_lru_cache //nolint:golint,stylecheck
import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*listItem
	mux      sync.Mutex
}

// Set - sets value for provided key in cache.
func (lru *lruCache) Set(key Key, value interface{}) bool {
	// Defense from concurrent writing
	lru.mux.Lock()
	defer lru.mux.Unlock()

	item, keyExists := lru.items[key]
	// If key is already presented in cache - just update it's value and move to the front of queue
	// (make it most "hot" item)
	if keyExists {
		item.value = cacheItem{Key: key, Value: value}
		lru.queue.MoveToFront(item)
	} else {
		// Otherwise new queue item is created and pushed to the front of the queue,
		// and target key is stored in helper map
		item = lru.queue.PushFront(cacheItem{Key: key, Value: value})
		lru.items[key] = item
		// If queue becomes bigger than cache capacity,
		// last item of queue (or the least "hot" one) must be removed and it's key either.
		if lru.queue.Len() > lru.capacity {
			lastItem := lru.queue.Back()
			lru.queue.Remove(lastItem)
			delete(lru.items, lastItem.value.(cacheItem).Key)
		}
	}
	return keyExists
}

// Get - get value for provided key from cache.
func (lru *lruCache) Get(key Key) (interface{}, bool) {
	// Defense from concurrent reading
	lru.mux.Lock()
	defer lru.mux.Unlock()

	item, keyExists := lru.items[key]
	// If key is presented in cache, than it must be set as the most "hot" and returned
	if keyExists {
		lru.queue.MoveToFront(item)
		return item.value.(cacheItem).Value, true
	}
	return nil, false
}

func (lru *lruCache) Clear() {
	lru.items = make(map[Key]*listItem)
	lru.queue = &list{}
}

type cacheItem struct {
	Key   Key
	Value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		items:    make(map[Key]*listItem),
		queue:    &list{},
	}
}
