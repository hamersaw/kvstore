package main

import (
	"context"
)

// request defines a single KV operation.
type request struct {
	key		     string
	value	     *string
	operation    operation
	responseChan chan *string
}

// operation defines which action the KVStore should perform .
type operation int

const (
	DELETE operation = iota
	GET
	PUT
	UPDATE
)

// ChannelKVStore is a KVStore implementation using a golang channel to enable concurrent requests
// to set, get, update, and delete KV objects.
type ChannelKVStore struct {
	requestChannel chan request
	store          map[string]*string
}

// Delete removes the specified key from the KVStore. This function returns ErrNotFound if the key
// does not exist in the KVStore.
func (c *ChannelKVStore) Delete(ctx context.Context, key string) error {
	responseChan := make(chan *string)
	c.requestChannel <- request{key, nil, DELETE, responseChan}

	select {
	case <-ctx.Done():
		return nil
	case value := <-responseChan:
		// in DELETE nil means the key was not found
		if value == nil {
			return ErrNotFound
		}
		return nil
	}
}

// Get returns the associated value for the key from the KVStore. This function returns ErrNotFound
// if the key does not exist in the KVStore.
func (c *ChannelKVStore) Get(ctx context.Context, key string) (*string, error) {
	responseChan := make(chan *string)
	c.requestChannel <- request{key, nil, GET, responseChan}

	select {
	case <-ctx.Done():
		return nil, nil
	case value := <-responseChan:
		// in GET nil means the key was not found
		if value == nil {
			return nil, ErrNotFound
		}
		return value, nil
	}
}

// Set stores the <key, value> pair in the KVStore. This function returns ErrMaxCapacity if the
// KVStore does not have capacity to insert the key.
func (c *ChannelKVStore) Set(ctx context.Context, key string, value *string) error {
	responseChan := make(chan *string)
	c.requestChannel <- request{key, value, PUT, responseChan}

	select {
	case <-ctx.Done():
		return nil
	case value := <-responseChan:
		// in PUT nil means store is at maximum capacity
		if value == nil {
			return ErrMaxCapacity
		}
		return nil
	}
}

// Update updates the value associated with the specified key in the KVStore. This function returns
// ErrNotFound if the key does not exist in the KVStore.
func (c *ChannelKVStore) Update(ctx context.Context, key string, value *string) error {
	responseChan := make(chan *string)
	c.requestChannel <- request{key, value, UPDATE, responseChan}

	select {
	case <-ctx.Done():
		return nil
	case value := <-responseChan:
		// in UPDATE nil means key was not found
		if value == nil {
			return ErrNotFound
		}
		return nil
	}
}

// NewChannelKVStore creates a new ChannelKVStore and starts the request handling goroutine.
func NewChannelKVStore(ctx context.Context, maxSize int) *ChannelKVStore {
	requestChannel := make(chan request)
	store := make(map[string]*string)

	go func() {
		for {
			select {
			case request := <-requestChannel:
				switch request.operation {
				case DELETE:
					if value, ok := store[request.key]; ok {
						delete(store, request.key)
						request.responseChan <- value
					} else {
						request.responseChan <- nil
					}
				case GET:
					request.responseChan <- store[request.key]
				case PUT:
					if len(store) >= maxSize {
						request.responseChan <- nil
					} else {
						store[request.key] = request.value
						request.responseChan <- request.value
					}
				case UPDATE:
					if _, ok := store[request.key]; ok {
						store[request.key] = request.value
						request.responseChan <- request.value
					} else {
						request.responseChan <- nil
					}
				}
			case <-ctx.Done():
				break
			}
		}
	}()

	return &ChannelKVStore{
		requestChannel: requestChannel,
		store:          store,
	}
}
