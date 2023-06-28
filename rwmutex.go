package main

import (
	"context"
	"sync"
)

// RWMutexKVStore is a KVStore implementation using a RWMutex to enable concurrent requests to set,
// get, update, and delete KV objects.
type RWMutexKVStore struct {
	maxCapacity int
	rwLock      sync.RWMutex
	store       map[string]*string
}

// Delete removes the specified key from the KVStore. This function returns ErrNotFound if the key
// does not exist in the KVStore.
func (c *RWMutexKVStore) Delete(ctx context.Context, key string) error {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()

	if _, ok := c.store[key]; !ok {
		return ErrNotFound
	}

	delete(c.store, key)
	return nil
}

// Get returns the associated value for the key from the KVStore. This function returns ErrNotFound
// if the key does not exist in the KVStore.
func (c *RWMutexKVStore) Get(ctx context.Context, key string) (*string, error) {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()

	value, ok := c.store[key]
	if !ok {
		return nil, ErrNotFound
	}

	return value, nil
}

// Set stores the <key, value> pair in the KVStore. This function returns ErrMaxCapacity if the
// KVStore does not have capacity to insert the key.
func (c *RWMutexKVStore) Set(ctx context.Context, key string, value *string) error {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()

	if len(c.store) > c.maxCapacity {
		return ErrMaxCapacity
	}

	c.store[key] = value
	return nil
}

// Update updates the value associated with the specified key in the KVStore. This function return
// ErrNotFound if the key does not exist in the KVStore.
func (c *RWMutexKVStore) Update(ctx context.Context, key string, value *string) error {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()

	_, ok := c.store[key]
	if !ok {
		return ErrNotFound
	}

	c.store[key] = value
	return nil
}

// NewRWMutexKVStore creates a new RWMutexKVStore with the sepcified maximum capacity.
func NewRWMutexKVStore(ctx context.Context, maxCapacity int) *RWMutexKVStore {
	return &RWMutexKVStore{
		maxCapacity: maxCapacity,
		rwLock:      sync.RWMutex{},
		store:       map[string]*string{},
	}
}
