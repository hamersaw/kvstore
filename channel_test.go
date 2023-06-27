package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChannelKVStoreDelete(t *testing.T) {
	ctx := context.TODO()
	value := "bar"

	t.Run("Success", func(t *testing.T) {
		kvstore := NewChannelKVStore(ctx, 1)
		channelKVStore, ok := kvstore.(*ChannelKVStore)
		assert.True(t, ok)

		channelKVStore.store["foo"] = &value
		assert.Equal(t, 1, len(channelKVStore.store))

		err := kvstore.Delete(ctx, "foo")
		assert.Nil(t, err)
		assert.Equal(t, 0, len(channelKVStore.store))
	})

	t.Run("FailureNotFound", func(t *testing.T) {
		kvstore := NewChannelKVStore(ctx, 1)
		err := kvstore.Delete(ctx, "foo")
		assert.ErrorIs(t, err, ErrNotFound)
	})
}

func TestChannelKVStoreSet(t *testing.T) {
	ctx := context.TODO()
	value := "bar"

	t.Run("Success", func(t *testing.T) {
		kvstore := NewChannelKVStore(ctx, 1)

		channelKVStore, ok := kvstore.(*ChannelKVStore)
		assert.True(t, ok)
		assert.Equal(t, 0, len(channelKVStore.store))

		err := kvstore.Set(ctx, "foo", &value)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(channelKVStore.store))
	})

	t.Run("FailureMaxCapacity", func(t *testing.T) {
		kvstore := NewChannelKVStore(ctx, 0)
		err := kvstore.Set(ctx, "foo", &value)
		assert.ErrorIs(t, err, ErrMaxCapacity)
	})
}

func TestChannelKVStoreGet(t *testing.T) {
	ctx := context.TODO()
	value := "bar"

	t.Run("Success", func(t *testing.T) {
		kvstore := NewChannelKVStore(ctx, 1)
		channelKVStore, ok := kvstore.(*ChannelKVStore)
		assert.True(t, ok)
		channelKVStore.store["foo"] = &value

		getValue, err := kvstore.Get(ctx, "foo")
		assert.Nil(t, err)
		assert.Equal(t, &value, getValue)
	})

	t.Run("FailureNotFound", func(t *testing.T) {
		kvstore := NewChannelKVStore(ctx, 1)

		getValue, err := kvstore.Get(ctx, "foo")
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Nil(t, getValue)
	})
}

func TestChannelKVStoreUpdate(t *testing.T) {
	ctx := context.TODO()
	value := "bar"
	updatedValue := "baz"

	t.Run("Success", func(t *testing.T) {
		kvstore := NewChannelKVStore(ctx, 1)
		channelKVStore, ok := kvstore.(*ChannelKVStore)
		assert.True(t, ok)
		channelKVStore.store["foo"] = &value

		err := kvstore.Update(ctx, "foo", &updatedValue)
		assert.Nil(t, err)
		assert.Equal(t, &updatedValue, channelKVStore.store["foo"])
	})

	t.Run("FailureNotFound", func(t *testing.T) {
		kvstore := NewChannelKVStore(ctx, 1)

		err := kvstore.Update(ctx, "foo", &updatedValue)
		assert.ErrorIs(t, err, ErrNotFound)
	})
}
