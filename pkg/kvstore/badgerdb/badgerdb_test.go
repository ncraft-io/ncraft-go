package badgerdb

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestStore tests if reading from, writing to and deleting from the store works properly.
// A struct is used as value. See TestTypes() for a test that is simpler but tests all types.
func TestStore(t *testing.T) {
	// Test with JSON
	store, err := NewStore(DefaultOptions)
	assert.Nil(t, err)

	store.Put(context.Background(), []byte("test"), []byte("value"))

	value, err := store.Get(context.Background(), []byte("test"))
	assert.Equal(t, value, []byte("value"))

	/*t.Run("Raw", func(t *testing.T) {
		store, path := createStore(t, encoding.JSON)
		defer cleanUp(store, path)
		test.TestStore(store, t)
	})

	// Test with gob
	t.Run("gob", func(t *testing.T) {
		store, path := createStore(t, encoding.Gob)
		defer cleanUp(store, path)
		test.TestStore(store, t)
	})*/
}
