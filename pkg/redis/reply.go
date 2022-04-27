package redis

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

func cannotConvert(d reflect.Kind, s interface{}) error {
	return fmt.Errorf("redigo: Scan cannot convert from %s to %s",
		reflect.TypeOf(s), d)
}

// Int is a helper that converts a command reply to an integer. If err is not
// equal to nil, then Int returns 0, err. Otherwise, Int converts the
// reply to an int as follows:
//
//  Reply type    Result
//  integer       int(reply), nil
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
func Int(reply interface{}, err error) (int, error) {
	//return cluster.Int(reply, err)
	n, err := Int64(reply, err)
	return int(n), err
}

// Int64 is a helper that converts a command reply to 64 bit integer. If err is
// not equal to nil, then Int returns 0, err. Otherwise, Int64 converts the
// reply to an int64 as follows:
//
//  Reply type    Result
//  integer       reply, nil
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
func Int64(reply interface{}, err error) (int64, error) {
	//return cluster.Int64(reply, err)
	if err != nil {
		return 0, err
	}

	if n, ok := reply.(int64); ok {
		return n, nil
	} else {
		return 0, cannotConvert(reflect.Int64, reply)
	}
}

// Float64 is a helper that converts a command reply to 64 bit float. If err is
// not equal to nil, then Float64 returns 0, err. Otherwise, Float64 converts
// the reply to an int as follows:
//
//  Reply type    Result
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
func Float64(reply interface{}, err error) (float64, error) {
	//return cluster.Float64(reply, err)
	if err != nil {
		return 0, err
	}

	if n, ok := reply.(float64); ok {
		return n, nil
	} else {
		return 0, cannotConvert(reflect.Float64, reply)
	}
}

// String is a helper that converts a command reply to a string. If err is not
// equal to nil, then String returns "", err. Otherwise String converts the
// reply to a string as follows:
//
//  Reply type      Result
//  bulk string     string(reply), nil
//  simple string   reply, nil
//  nil             "",  ErrNil
//  other           "",  error
func String(reply interface{}, err error) (string, error) {
	//return cluster.String(reply, err)
	if err != nil {
		return "", err
	}

	switch reply.(type) {
	case string:
		return reply.(string), nil
	case []byte:
		return string(reply.([]byte)), nil
	default:
		return "", cannotConvert(reflect.String, reply)
	}
}

// Bytes is a helper that converts a command reply to a slice of bytes. If err
// is not equal to nil, then Bytes returns nil, err. Otherwise Bytes converts
// the reply to a slice of bytes as follows:
//
//  Reply type      Result
//  bulk string     reply, nil
//  simple string   []byte(reply), nil
//  nil             nil, ErrNil
//  other           nil, error
func Bytes(reply interface{}, err error) ([]byte, error) {
	//return cluster.Bytes(reply, err)
	if err != nil {
		return nil, err
	}

	switch reply.(type) {
	case string:
		return []byte(reply.(string)), nil
	case []byte:
		return reply.([]byte), nil
	default:
		return nil, cannotConvert(reflect.Slice, reply)
	}
}

// Bool is a helper that converts a command reply to a boolean. If err is not
// equal to nil, then Bool returns false, err. Otherwise Bool converts the
// reply to boolean as follows:
//
//  Reply type      Result
//  integer         value != 0, nil
//  bulk string     strconv.ParseBool(reply)
//  nil             false, ErrNil
//  other           false, error
func Bool(reply interface{}, err error) (bool, error) {
	//return cluster.Bool(reply, err)
	if err != nil {
		return false, err
	}

	if b, ok := reply.(bool); ok {
		return b, nil
	} else {
		return false, cannotConvert(reflect.Bool, reply)
	}
}

// Values is a helper that converts an array command reply to a []interface{}.
// If err is not equal to nil, then Values returns nil, err. Otherwise, Values
// converts the reply as follows:
//
//  Reply type      Result
//  array           reply, nil
//  nil             nil, ErrNil
//  other           nil, error
func Values(reply interface{}, err error) ([]interface{}, error) {
	//return cluster.Values(reply, err)
	if err != nil {
		return nil, err
	}

	if slice, ok := reply.([]interface{}); ok {
		return slice, nil
	} else {
		return nil, cannotConvert(reflect.Slice, reply)
	}
}

// Ints is a helper that converts an array command reply to a []int.
// If err is not equal to nil, then Ints returns nil, err.
func Ints(reply interface{}, err error) ([]int, error) {
	//return cluster.Ints(reply, err)
	if err != nil {
		return nil, err
	}

	values, err := Values(reply, err)
	if err != nil {
		return nil, err
	}

	nums := make([]int, 0, len(values))
	for _, value := range values {
		if n, ok := value.(int64); ok {
			nums = append(nums, int(n))
		} else {
			nums = append(nums, 0)
		}
	}

	return nums, nil
}

// Strings is a helper that converts an array command reply to a []string. If
// err is not equal to nil, then Strings returns nil, err. Nil array items are
// converted to "" in the output slice. Strings returns an error if an array
// item is not a bulk string or nil.
func Strings(reply interface{}, err error) ([]string, error) {
	//return cluster.Strings(reply, err)
	if err != nil {
		return nil, err
	}

	values, err := Values(reply, err)
	if err != nil {
		return nil, err
	}

	strs := make([]string, 0, len(values))
	for _, value := range values {
		switch value.(type) {
		case string:
			strs = append(strs, value.(string))
		case []byte:
			strs = append(strs, string(value.([]byte)))
		default:
			strs = append(strs, "")
		}
	}

	return strs, nil
}

// StringMap is a helper that converts an array of strings (alternating key, value)
// into a map[string]string. The HGETALL and CONFIG GET commands return replies in this format.
// Requires an even number of values in result.
func StringMap(result interface{}, err error) (map[string]string, error) {
	//return cluster.StringMap(result, err)
	if err != nil {
		return nil, err
	}

	values, err := Values(result, err)
	if err != nil {
		return nil, err
	}

	if len(values)%2 != 0 {
		return nil, errors.New("bulk string number not even")
	}

	m := make(map[string]string, len(values)/2)
	for idx := 0; idx < len(values); {
		var key, value string

		switch values[idx].(type) {
		case string:
			key = values[idx].(string)
		case []byte:
			key = string(values[idx].([]byte))
		default:
			return nil, cannotConvert(reflect.String, values[idx])
		}

		switch values[idx+1].(type) {
		case string:
			value = values[idx+1].(string)
		case []byte:
			value = string(values[idx+1].([]byte))
		default:
			return nil, cannotConvert(reflect.String, values[idx+1])
		}

		m[key] = value
		idx += 2
	}

	return m, nil
}

// Scan copies from src to the values pointed at by dest.
//
// The values pointed at by dest must be an integer, float, boolean, string,
// []byte, interface{} or slices of these types. Scan uses the standard strconv
// package to convert bulk strings to numeric and boolean types.
//
// If a dest value is nil, then the corresponding src value is skipped.
//
// If a src element is nil, then the corresponding dest value is not modified.
//
// To enable easy use of Scan in a loop, Scan returns the slice of src
// following the copied values.
//func Scan(src []interface{}, dst ...interface{}) ([]interface{}, error) {
//	return cluster.Scan(src, dst...)
//}
