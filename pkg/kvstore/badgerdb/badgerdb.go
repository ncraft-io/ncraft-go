package badgerdb

import (
    "context"
    "errors"
    "github.com/dgraph-io/badger"
    "github.com/ncraft-io/ncraft-go/pkg/kvstore"
)

// Store is a gokv.Store implementation for BadgerDB.
type Store struct {
    db         *badger.DB
    bucketName string
}

func init() {
    kvstore.RegisterStore("badgerdb",
        func(options kvstore.Options) kvstore.KvStore {
            o := Options{}
            err := kvstore.ResetOptions(options, &o)
            if err != nil {
                return nil
            }

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
func (s *Store) Put(ctx context.Context, k, v []byte) error {
    if err := kvstore.CheckKey(k); err != nil {
        return err
    }

    // First turn the passed object into something that BadgerDB can handle
    err := s.db.Update(func(txn *badger.Txn) error {
        return txn.Set(kvstore.BucketKey(s.bucketName, k), v)
    })
    if err != nil {
        return err
    }
    return nil
}

func (s *Store) BatchPut(ctx context.Context, keys, values [][]byte) error {
    return errors.New("not implement")
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
    err = s.db.View(func(txn *badger.Txn) error {
        item, err := txn.Get(kvstore.BucketKey(s.bucketName, k))
        if err != nil {
            return err
        }
        // item.Value() is only valid within the transaction.
        // We can either copy it ourselves or use the ValueCopy() method.
        // TODO: Benchmark if it's faster to copy + close tx,
        // or to keep the tx open until unmarshalling is done.
        data, err = item.ValueCopy(nil)
        if err != nil {
            return err
        }
        return nil
    })
    // If no value was found return false
    if err == badger.ErrKeyNotFound {
        return nil, nil
    } else if err != nil {
        return nil, err
    }

    return data, nil
}

func (s Store) BatchGet(ctx context.Context, keys [][]byte) ([][]byte, error) {
    return nil, errors.New("not implement")
}

// Delete deletes the stored value for the given key.
// Deleting a non-existing key-value pair does NOT lead to an error.
// The key must not be "".
func (s *Store) Delete(ctx context.Context, key []byte) error {
    if err := kvstore.CheckKey(key); err != nil {
        return err
    }

    return s.db.Update(func(txn *badger.Txn) error {
        return txn.Delete(kvstore.BucketKey(s.bucketName, key))
    })
}

func (s *Store) BatchDelete(ctx context.Context, keys [][]byte) error {
    return errors.New("not implement")
}

func (s Store) Scan(ctx context.Context, startKey, endKey []byte, limit int) (keys [][]byte, values [][]byte, err error) {
    return nil, nil, errors.New("not implement")
}

func (s *Store) DeleteRange(ctx context.Context, startKey []byte, endKey []byte) error {
    return errors.New("not implement")
}

// Close closes the store.
// It must be called to make sure that all pending updates make their way to disk.
func (s Store) Close() error {
    return s.db.Close()
}

// Options are the options for the BadgerDB store.
type Options struct {
    // Directory for storing the DB files.
    // Optional ("BadgerDB" by default).
    Dir string `json:"dir"`

    BucketName string `json:"bucketName"`
}

// DefaultOptions is an Options object with default values.
// Dir: "BadgerDB", Codec: encoding.JSON
var DefaultOptions = Options{
    Dir: "BadgerDB",
}

// NewStore creates a new BadgerDB store.
// Note: BadgerDB uses an exclusive write lock on the database directory so it cannot be shared by multiple processes.
// So when creating multiple clients you should always use a new database directory (by setting a different Path in the options).
//
// You must call the Close() method on the store when you're done working with it.
func NewStore(options Options) (kvstore.KvStore, error) {
    result := &Store{}

    // Set default values
    if options.Dir == "" {
        options.Dir = DefaultOptions.Dir
    }

    result.bucketName = options.BucketName

    // Open the Badger database located in the options.Dir directory.
    // It will be created if it doesn't exist.
    opts := badger.DefaultOptions(options.Dir)
    db, err := badger.Open(opts)
    if err != nil {
        return result, err
    }

    result.db = db
    return result, nil
}
