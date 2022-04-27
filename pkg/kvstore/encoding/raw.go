package encoding

import (
	"errors"
)

type RawCodec struct{}

// Marshal encodes a Go value to JSON.
func (c RawCodec) Marshal(v interface{}) ([]byte, error) {
	switch v.(type) {
	case []byte:
		return v.([]byte), nil
	case string:
		return []byte(v.(string)), nil
	default:
		return nil, errors.New("not support")
	}
}

// Unmarshal decodes a JSON value into a Go value.
func (c RawCodec) Unmarshal(data []byte, v interface{}) error {
	switch v.(type) {
	case *[]byte:
		v = data
		return nil
	default:
		return errors.New("not support")
	}
}
