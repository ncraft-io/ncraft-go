package kvstore

import "bytes"

func BucketKey(bucket string, key []byte) []byte {
    buff := bytes.NewBufferString(bucket)
    if len(bucket) > 0 {
        buff.WriteByte('.')
    }

    buff.Write(key)
    return buff.Bytes()
}
