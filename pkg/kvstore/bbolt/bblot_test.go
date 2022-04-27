package bbolt

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

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

func TestBatchPut(t *testing.T) {
	// Test with JSON
	store, err := NewStore(DefaultOptions)
	assert.Nil(t, err)

	keys := make([][]byte, 0, 1000)
	values := make([][]byte, 0, 1000)
	for i := 0; i < 1000; i++ {
		keys = append(keys, []byte(fmt.Sprintf("%d", i)))
		values = append(values, []byte(fmt.Sprintf("%d", 2*i)))
	}

	if err := store.BatchPut(context.Background(), keys, values); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		value, err := store.Get(context.Background(), []byte(fmt.Sprintf("%d", i)))
		if err != nil {
			t.Fatal(err)
		}
		v, err := strconv.ParseInt(string(value), 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, v, int64(2*i))
	}
}

func TestBatchScan(t *testing.T) {
	store, err := NewStore(DefaultOptions)
	assert.Nil(t, err)

	keys := make([][]byte, 0, 1000)
	values := make([][]byte, 0, 1000)
	for i := 0; i < 1000; i++ {
		keys = append(keys, []byte(fmt.Sprintf("%d", i)))
		values = append(values, []byte(fmt.Sprintf("%d", 2*i)))
	}

	if err := store.BatchPut(context.Background(), keys, values); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		value, err := store.Get(context.Background(), []byte(fmt.Sprintf("%d", i)))
		if err != nil {
			t.Fatal(err)
		}
		v, err := strconv.ParseInt(string(value), 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, v, int64(2*i))
	}

	keys, values, err = store.Scan(context.Background(), []byte(fmt.Sprintf("%d", 10)),
		[]byte(fmt.Sprintf("%d", 101)), 10000)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(keys), 2)
	for i := range keys {
		v, err := strconv.ParseInt(string(values[i]), 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		k, err := strconv.ParseInt(string(keys[i]), 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, v, int64(2*k))
	}
}