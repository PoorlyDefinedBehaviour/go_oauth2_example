package inmemorycache

import (
	"runtime"
	"sync"
	"time"
)

type entry struct {
	value     interface{}
	expiresAt *time.Time
}

type InMemoryCache struct {
	mutex sync.RWMutex
	store map[interface{}]entry
}

func New() *InMemoryCache {
	cache := &InMemoryCache{
		mutex: sync.RWMutex{},
		store: make(map[interface{}]entry),
	}

	cleaner := expiredEntryCleaner{Stop: make(chan struct{})}

	runtime.SetFinalizer(cache, func(cache *InMemoryCache) {
		cleaner.Stop <- struct{}{}
	})

	go cleaner.deleteEntriesWhenTheyExpire(cache)

	return cache
}

type expiredEntryCleaner struct {
	Stop chan struct{}
}

func (cleaner *expiredEntryCleaner) deleteEntriesWhenTheyExpire(cache *InMemoryCache) {
	itsTimeToDeleteExpiredEntries := time.Tick(1 * time.Second)

	for {
		select {
		case <-itsTimeToDeleteExpiredEntries:
			now := time.Now()

			cache.mutex.Lock()

			for key, entry := range cache.store {
				if entry.expiresAt != nil && entry.expiresAt.Before(now) {
					delete(cache.store, key)
				}
			}

			cache.mutex.Unlock()

		case <-cleaner.Stop:
			return
		}
	}
}

func (cache *InMemoryCache) SetWithExpiration(key interface{}, value interface{}, expiration time.Duration) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	expiresAt := time.Now().Add(expiration)

	cache.store[key] = entry{
		value:     value,
		expiresAt: &expiresAt,
	}
}

func (cache *InMemoryCache) Get(key interface{}) (interface{}, bool) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	entry, ok := cache.store[key]
	if !ok {
		return nil, false
	}

	return entry.value, true
}

func (cache *InMemoryCache) Has(key interface{}) bool {
	_, keyFound := cache.Get(key)
	return keyFound
}

func (cache *InMemoryCache) Delete(key interface{}) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	delete(cache.store, key)
}
