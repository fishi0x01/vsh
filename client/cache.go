package client

import (
	"sync"

	"github.com/hashicorp/vault/api"
)

// Cache is a thread-safe cache for vault queries
type Cache struct {
	vaultClient *api.Client
	listQueries map[string]*cacheElement
	mutex       sync.Mutex
}

type cacheElement struct {
	Secret *api.Secret
	Err    error
}

// NewCache create a new thread-safe cache object
func NewCache(vaultClient *api.Client) *Cache {
	return &Cache{
		vaultClient: vaultClient,
		listQueries: make(map[string]*cacheElement),
	}
}

// Clear removes all entries from the cache
func (cache *Cache) Clear() {
	cache.mutex.Lock()
	cache.listQueries = make(map[string]*cacheElement)
	cache.mutex.Unlock()
}

// List tries to get path from cache.
// If path is not available in cache, it uses the vault client to query and update cache.
func (cache *Cache) List(path string) (result *api.Secret, err error) {
	// try to get path from cache
	cache.mutex.Lock()
	value, hit := cache.listQueries[path]
	cache.mutex.Unlock()
	if hit {
		return value.Secret, value.Err
	}

	// not found in cache -> query vault
	// NOTE: this part is not mutexed
	result, err = cache.vaultClient.Logical().List(path)

	// update cache
	cache.mutex.Lock()
	cache.listQueries[path] = &cacheElement{
		Secret: result,
		Err:    err,
	}
	cache.mutex.Unlock()
	return result, err
}
