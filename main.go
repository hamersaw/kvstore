package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

const (
	channel = "channel"
	rwmutex = "rwmutex"
)

var ErrNotFound = fmt.Errorf("not found")
var ErrMaxCapacity = fmt.Errorf("maximum capacity")

// KVStore represents an in-memory, concurrent <key, value> data store.
type KVStore interface {
	// Delete removes the specified key from the KVStore. This function returns ErrNotFound if the
	// key does not exist in the KVStore.
	Delete(ctx context.Context, key string) error

	// Get returns the associated value for the key from the KVStore. This function returns
	// ErrNotFound if the key does not exist in the KVStore.
	Get(ctx context.Context, key string) (*string, error)

	// Set stores the <key, value> pair in the KVStore. This function returns ErrMaxCapacity if the
	// KVStore does not have capacity to insert the key.
	Set(ctx context.Context, key string, value *string) error

	// Update updates the value associated with the specified key in the KVStore. This function
	// returns ErrNotFound if the key does not exist in the KVStore.
	Update(ctx context.Context, key string, value *string) error
}

// KeyValue represents a <key, value> pair.
type KeyValue struct {
	Key   string
	Value string
}

func main() {
	ctx := context.Background()

	// parse command line arguments
	concurrencyEngine := flag.String("concurrency", "channel", "concurrency engine to use (channel, rwmutex)")
	maxSize := flag.Int("max-size", 50000, "maximum size of the kvstore")
	flag.Parse()

	// initialize kvstore
	var kvstore KVStore
	switch *concurrencyEngine {
	case channel:
		kvstore = NewChannelKVStore(ctx, *maxSize)
	case rwmutex:
		kvstore = NewRWMutexKVStore(ctx, *maxSize)
	default:
		panic(fmt.Sprintf("invalid concurrency engine '%s'", *concurrencyEngine))
	}

	// initialize rest server
	router := chi.NewRouter()

	router.Delete("/{key}", func(w http.ResponseWriter, r *http.Request) {
		// retrieve input from request
		key := chi.URLParam(r, "key")

		err := kvstore.Delete(ctx, key)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Write([]byte(fmt.Sprintf("%s", err)))
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})

	router.Get("/{key}", func(w http.ResponseWriter, r *http.Request) {
		// retrieve input from request
		key := chi.URLParam(r, "key")

		value, err := kvstore.Get(ctx, key)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Write([]byte(fmt.Sprintf("%s", err)))
		} else if value == nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("value is nil"))
		} else {
			w.Write([]byte(*value))
		}
	})

	router.Post("/{key}", func(w http.ResponseWriter, r *http.Request) {
		// retrieve input from request
		key := chi.URLParam(r, "key")
		value := r.FormValue("value")
		if len(value) <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := kvstore.Set(ctx, key, &value)
		if err != nil {
			if errors.Is(err, ErrMaxCapacity) {
				w.WriteHeader(http.StatusInsufficientStorage)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Write([]byte(fmt.Sprintf("%s", err)))
		} else {
			w.WriteHeader(http.StatusCreated)
		}
	})

	router.Route("/bulk", func(router chi.Router) {
		router.Put("/", func(w http.ResponseWriter, r *http.Request) {
			// decode key values from request body
			var keyValues []KeyValue
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&keyValues)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("%s", err)))
				return
			}

			// update key values in kvstore
			notFoundKeys := make([]string, 0)
			for _, keyValue := range keyValues {
				err := kvstore.Update(ctx, keyValue.Key, &keyValue.Value)
				if err != nil {
					if errors.Is(err, ErrNotFound) {
						notFoundKeys = append(notFoundKeys, keyValue.Key)
					} else {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(fmt.Sprintf("%s", err)))
						return
					}
				}
			}

			// return not found keys if any
			if len(notFoundKeys) > 0 {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(fmt.Sprintf("keys not found: %s", notFoundKeys)))
			} else {
				w.WriteHeader(http.StatusOK)
			}
		})
	})

	http.ListenAndServe(":3000", router)

	// never reachable - implementing sig handler is out of the scope of this project
	ctx.Done()
}
