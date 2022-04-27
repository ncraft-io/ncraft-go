package encoding

import (
	jsoniter "github.com/json-iterator/go"
)

// JSONcodec encodes/decodes Go values to/from JSON.
// You can use encoding.JSON instead of creating an instance of this struct.
type JsonCodec struct{}

// Marshal encodes a Go value to JSON.
func (c JsonCodec) Marshal(v interface{}) ([]byte, error) {
	return jsoniter.ConfigFastest.Marshal(v)
}

// Unmarshal decodes a JSON value into a Go value.
func (c JsonCodec) Unmarshal(data []byte, v interface{}) error {
	return jsoniter.ConfigFastest.Unmarshal(data, v)
}
