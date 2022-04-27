package bbolt

import (
    "bytes"
    "errors"
    bolt "go.etcd.io/bbolt"

    "context"
    "github.com/ncraft-io/ncraft-go/pkg/kvstore"
)

// Store is a gokv.Store implementation for bbolt (formerly known as Bolt / Bolt DB).
type Store struct {
    db         *bolt.DB
    bucketName string
}

func init() {
    kvstore.RegisterStore("bblot",
        func(options kvstore.Options) kvstore.KvStore {
            o := Options{}
            kvstore.ResetOptions(options, &o)
            store, err := NewStore(o)
            if err != nil {
                return nil
            }
            return store
        })
}

// Put stores the given value for the given key.
// Values are automatically marshalled to JSON or gob (depending on the configuration).
// The key must not be "" and the value must not be nil.
func (s *Store) Put(ctx context.Context, k []byte, v []byte) error {
    if err := kvstore.CheckKey(k); err != nil {
        return err
    }

    err := s.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(s.bucketName))
        return b.Put(k, v)
    })
    if err != nil {
        return err
    }
    return nil
}

func (s *Store) BatchPut(ctx context.Context, keys, values [][]byte) error {
    if len(keys) != len(values) {
        return errors.New("args error")
    }
    for i := range keys {
        if err := kvstore.CheckKey(keys[i]); err != nil {
            return err
        }
    }

    err := s.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(s.bucketName))
        for i := range keys {
            if err := b.Put(keys[i], values[i]); err != nil {
                return err
            }
        }
        return nil
    })
    if err != nil {
        return err
    }
    return nil
}

// Get retrieves the stored value for the given key.
// You need to pass a pointer to the value, so in case of a struct
// the automatic unmarshalling can populate the fields of the object
// that v points to with the values of the retrieved object's values.
// If no value is found it returns (false, nil).
// The key must not be "" and the pointer must not be nil.
func (s Store) Get(ctx context.Context, k []byte) (v []byte, err error) {
    if err := kvstore.CheckKey(k); err != nil {
        return nil, err
    }

    var data []byte
    err = s.db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(s.bucketName))
        txData := b.Get(k)
        // txData is only valid during the transaction.
        // Its value must be copied to make it valid outside of the tx.
        // TODO: Benchmark if it's faster to copy + close tx,
        // or to keep the tx open until unmarshalling is done.
        if txData != nil {
            // `data = append([]byte{}, txData...)` would also work, but the following is more explicit
            data = make([]byte, len(txData))
            copy(data, txData)
        }
        return nil
    })
    if err != nil {
        return nil, nil
    }

    // If no value was found return false
    if data == nil {
        return nil, nil
    }

    return data, nil
}

func (s Store) BatchGet(ctx context.Context, keys [][]byte) ([][]byte, error) {
    var datas [][]byte
    err := s.db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(s.bucketName))
        for _, key := range keys {
            if err := kvstore.CheckKey(key); err != nil {
                datas = append(datas, []byte(""))
            }

            txData := b.Get(key)
            // txData is only valid during the transaction.
            // Its value must be copied to make it valid outside of the tx.
            // TODO: Benchmark if it's faster to copy + close tx,
            // or to keep the tx open until unmarshalling is done.
            if txData != nil {
                // `data = append([]byte{}, txData...)` would also work, but the following is more explicit
                data := make([]byte, len(txData))
                copy(data, txData)
                datas = append(datas, data)
            } else {
                datas = append(datas, []byte(""))
            }
        }

        return nil
    })
    if err != nil {
        return nil, nil
    }

    // If no value was found return false
    if len(datas) == 0 {
        return nil, nil
    }

    return datas, nil
}

// Delete deletes the stored value for the given key.
// Deleting a non-existing key-value pair does NOT lead to an error.
// The key must not be "".
func (s Store) Delete(ctx context.Context, k []byte) error {
    if err := kvstore.CheckKey(k); err != nil {
        return err
    }

    return s.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(s.bucketName))
        return b.Delete([]byte(k))
    })
}

func (s *Store) BatchDelete(ctx context.Context, keys [][]byte) error {
    return errors.New("not implement")
}

func (s Store) Scan(ctx context.Context, startKey, endKey []byte, limit int) (keys [][]byte, values [][]byte, err error) {
    lessThanEndKey := func(k []byte) bool {
        if endKey == nil {
            return true
        } else {
            return bytes.Compare(k, endKey) < 0
        }
    }
    err = s.db.View(func(tx *bolt.Tx) error {
        c := tx.Bucket([]byte(s.bucketName)).Cursor()
        cnt := 0

        // Iterate over the keys
        for k, v := c.Seek(startKey); k != nil && lessThanEndKey(k); k, v = c.Next() {
            if cnt > limit {
                break
            }
            if v != nil {
                keys = append(keys, k)
                values = append(values, v)
                cnt++
            }
        }
        return nil
    })
    return keys, values, err
}

func (s *Store) DeleteRange(ctx context.Context, startKey []byte, endKey []byte) error {
    return errors.New("not implement")
}

// Close closes the store.
// It must be called to make sure that all open transactions finish and to release all DB resources.
func (s Store) Close() error {
    return s.db.Close()
}

// Options are the options for the bbolt store.
type Options struct {
    // Bucket name for storing the key-value pairs.
    // Optional ("default" by default).
    BucketName string `json:"bucketName"`
    // Path of the DB file.
    // Optional ("bbolt.db" by default).
    Path string `json:"path"`
}

// DefaultOptions is an Options object with default values.
// BucketName: "default", Path: "bbolt.db", Codec: encoding.JSON
var DefaultOptions = Options{
    BucketName: "default",
    Path:       "bbolt.db",
}

// NewStore creates a new bbolt store.
// Note: bbolt uses an exclusive write lock on the database file so it cannot be shared by multiple processes.
// So when creating multiple clients you should always use a new database file (by setting a different Path in the options).
//
// You must call the Close() method on the store when you're done working with it.
func NewStore(options Options) (kvstore.KvStore, error) {
    result := &Store{}

    // Set default values
    if options.BucketName == "" {
        options.BucketName = DefaultOptions.BucketName
    }
    if options.Path == "" {
        options.Path = DefaultOptions.Path
    }

    // Open DB
    db, err := bolt.Open(options.Path, 0600, nil)
    if err != nil {
        return result, err
    }

    // Create a bucket if it doesn't exist yet.
    // In bbolt key/value pairs are stored to and read from buckets.
    err = db.Update(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists([]byte(options.BucketName))
        if err != nil {
            return err
        }
        return nil
    })
    if err != nil {
        return result, err
    }

    result.db = db
    result.bucketName = options.BucketName
    return result, nil
}
