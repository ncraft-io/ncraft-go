package kvstore

import (
    "errors"
    jsoniter "github.com/json-iterator/go"
)

func ResetOptions(options map[string]interface{}, newOptions interface{}) error {
    bytes, err := jsoniter.ConfigFastest.Marshal(options)
    if err != nil {
        return err
    }
    return jsoniter.ConfigFastest.Unmarshal(bytes, newOptions)
}

// CheckKey returns an error if k == ""
func CheckKey(k []byte) error {
    if len(k) == 0 {
        return errors.New("the passed key is an empty bytes, which is invalid")
    }
    return nil
}

// CheckValue returns an error if v == nil
func CheckValue(v interface{}) error {
    if v == nil {
        return errors.New("the passed value is nil, which is not allowed")
    }
    return nil
}
