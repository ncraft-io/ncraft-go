package tikv

import (
	"context"
	"github.com/bmizerany/assert"
	"testing"
)

func TestTikv(t *testing.T) {
	s, err := NewStore(Options{Pds: []string{"172.3.0.122:2479"}})
	assert.Equal(t, nil, err)
	defer s.Close()

	// 1. test single key
	err = s.Put(context.TODO(), []byte("foo"), []byte("bar"))
	assert.Equal(t, nil, err)

	v, err := s.Get(context.TODO(), []byte("foo"))
	assert.Equal(t, nil, err)
	assert.Equal(t, "bar", string(v))

	err = s.Delete(context.TODO(), []byte("foo"))
	assert.Equal(t, nil, err)

	v, err = s.Get(context.TODO(), []byte("foo"))
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(nil), v)

	// 2. test batch keys
	err = s.BatchPut(context.TODO(), [][]byte{[]byte("foo1"), []byte("foo2")}, [][]byte{[]byte("bar1"), []byte("bar2")})
	assert.Equal(t, nil, err)

	values, err := s.BatchGet(context.TODO(), [][]byte{[]byte("foo1"), []byte("foo2")})
	assert.Equal(t, nil, err)
	assert.Equal(t, "bar1", string(values[0]))
	assert.Equal(t, "bar2", string(values[1]))

	err = s.BatchDelete(context.TODO(), [][]byte{[]byte("foo1"), []byte("foo2")})
	assert.Equal(t, nil, err)

	values, err = s.BatchGet(context.TODO(), [][]byte{[]byte("foo1"), []byte("foo2")})
	assert.Equal(t, nil, err)
	assert.Equal(t, [][]byte{nil, nil}, values)
}
