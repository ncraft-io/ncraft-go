package redis

import (
    "context"
    "testing"
)

func TestString(t *testing.T) {
    client := NewGoRedis(Config{Connections: []string{"127.0.0.1:6379"}})

    testKey := `test_key\ntest_key`
    testValue := `test_value\ntest_value`
    if _, err := client.Do(context.Background(), "set", testKey, testValue); err != nil {
        t.Fatal(err)
    }

    if str, err := String(client.Do(context.Background(), "get", testKey)); err != nil || str != testValue {
        t.Fatal(err)
    }
}

func TestStrings(t *testing.T) {
    client := NewGoRedis(Config{Connections: []string{"127.0.0.1:6379"}})

    testKey := `hm_test_key`
    testField1 := `test_field1`
    testField2 := `test_field2`
    testValue1 := `test_value1`
    testValue2 := `test_value2`

    invalidField := `invalid_field`

    _, _ = client.Do(context.Background(), "del", testKey)
    if _, err := client.Do(context.Background(), "hmset", testKey, testField1, testValue1, testField2, testValue2); err != nil {
        t.Fatal(err)
    }

    if strs, err := Strings(client.Do(context.Background(), "hmget", testKey, invalidField, testField1, testField2)); err != nil {
        t.Fatal(err)
    } else {
        if len(strs) != 3 {
            t.Fatal("length unexpected")
        }
        expectStrs := []string{"", testValue1, testValue2}
        for idx := range strs {
            if strs[idx] != expectStrs[idx] {
                t.Fatal("value not equal")
            }
        }
    }
}

func TestInt(t *testing.T) {
    client := NewGoRedis(Config{Connections: []string{"127.0.0.1:6379"}})

    _, _ = client.Do(context.Background(), "del", "int_key")
    resp, err := client.Do(context.Background(), "incr", "int_key")
    if err != nil {
        t.Fatal(err)
    }

    if n, err := Int(resp, err); err != nil || n != 1 {
        t.Fatal("expect 1, got ", n)
    }
}

func TestInt64(t *testing.T) {
    client := NewGoRedis(Config{Connections: []string{"127.0.0.1:6379"}})

    _, _ = client.Do(context.Background(), "del", "int_key")
    resp, err := client.Do(context.Background(), "incr", "int_key")
    if err != nil {
        t.Fatal(err)
    }

    if n, err := Int64(resp, err); err != nil || n != 1 {
        t.Fatal("expect 1, got ", n)
    }
}

func TestFloat64(t *testing.T) {

}

func TestBytes(t *testing.T) {
    client := NewGoRedis(Config{Connections: []string{"127.0.0.1:6379"}})
    //client := NewPool(Config{Connections: []string{"127.0.0.1:6379"}})
    testKey := `byte_key`
    testValue := []byte{12, 36, 134, 112, 12, 4, 6, 77, 22}
    //testValue := []byte("sfnowb[qmfn")
    if resp, err := client.Do(context.Background(), "set", testKey, testValue); err != nil {
        t.Fatal(err)
    } else {
        _ = resp
    }

    if bs, err := Bytes(client.Do(context.Background(), "get", testKey)); err != nil {
        t.Fatal(err)
    } else if !bytesEqual(bs, testValue) {
        t.Fatal("bytes not equal")
    }
}

func TestBool(t *testing.T) {
}

func TestInts(t *testing.T) {
}

func TestStringMap(t *testing.T) {
    client := NewGoRedis(Config{Connections: []string{"127.0.0.1:6379"}})

    testKey := `hm_test_key`
    testField1 := `test_field1`
    testField2 := `test_field2`
    testField3 := []byte(`test_field3`)
    testValue1 := `test_value1`
    testValue2 := []byte{0, 1, 2, 3, 4, '\n', 5}
    testValue3 := []byte{1, 2, 3, 4, 5, '\n', 5}

    _, _ = client.Do(context.Background(), "del", testKey)
    if _, err := client.Do(context.Background(), "hmset", testKey, testField1, testValue1, testField2, testValue2, testField3, testValue3); err != nil {
        t.Fatal(err)
    }

    if m, err := StringMap(client.Do(context.Background(), "hgetall", testKey)); err != nil {
        t.Fatal(err)
    } else {
        if len(m) != 3 {
            t.Fatal("length unexpected")
        }
        if m[testField1] != testValue1 || m[testField2] != string(testValue2) || m[string(testField3)] != string(testValue3) {
            t.Fatal("value unexpected")
        }
    }
}

func TestStringMap2(t *testing.T) {
    client := NewGoRedis(Config{Connections: []string{"127.0.0.1:6379"}})

    testKey := `hm_test_key`
    testField1 := []byte{110, 111, 21, 113, 14, 115, 16, 17}
    testField2 := []byte{10, 11, 112, 13, 141, 15, 16, 17}
    testField3 := []byte{20, 11, 112, 13, 14, 115, 16, 27}
    testValue1 := []byte{10, 11, 112, 3, 41, '\n', 5}
    testValue2 := []byte{0, 1, 211, 31, 4, '\n', 5}
    testValue3 := []byte{1, 2, 3, 4, 51, '\n', 5}

    _, _ = client.Do(context.Background(), "del", testKey)
    if _, err := client.Do(context.Background(), "hmset", testKey, testField1, testValue1, testField2, testValue2, testField3, testValue3); err != nil {
        t.Fatal(err)
    }

    if m, err := StringMap(client.Do(context.Background(), "hgetall", testKey)); err != nil {
        t.Fatal(err)
    } else {
        if len(m) != 3 {
            t.Fatal("length unexpected")
        }
        if m[string(testField1)] != string(testValue1) || m[string(testField2)] != string(testValue2) || m[string(testField3)] != string(testValue3) {
            t.Fatal("value unexpected")
        }
    }
}

func bytesEqual(bs1, bs2 []byte) bool {
    if len(bs1) != len(bs2) {
        return false
    }
    for idx := range bs1 {
        if bs1[idx] != bs2[idx] {
            return false
        }
    }
    return true
}
