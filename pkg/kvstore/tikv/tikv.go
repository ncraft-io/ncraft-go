package tikv

import (
    "context"
    "github.com/ncraft-io/ncraft-go/pkg/kvstore"
    ticonfig "github.com/tikv/client-go/config"
    "github.com/tikv/client-go/rawkv"
)

type Store struct {
    cli        *rawkv.Client
    bucketName string
}

func init() {
    kvstore.RegisterStore("tikv",
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

type Options struct {
    Pds        []string `json:"pds"`
    BucketName string   `json:"bucketName"`
}

func NewStore(conf Options) (kvstore.KvStore, error) {
    cli, err := rawkv.NewClient(context.TODO(), conf.Pds, ticonfig.Default())
    if err != nil {
        return nil, err
    }

    return &Store{
        cli:        cli,
        bucketName: conf.BucketName,
    }, nil
}

func (s Store) Key(k []byte) []byte {
    return kvstore.BucketKey(s.bucketName, k)
}

func (s Store) Keys(ks [][]byte) [][]byte {
    keys := make([][]byte, 0, len(ks))
    for _, key := range ks {
        keys = append(keys, kvstore.BucketKey(s.bucketName, key))
    }
    return keys
}

func (s *Store) Put(ctx context.Context, k []byte, v []byte) error {
    return s.cli.Put(ctx, s.Key(k), v)
}

func (s *Store) BatchPut(ctx context.Context, keys, values [][]byte) error {
    return s.cli.BatchPut(ctx, s.Keys(keys), values)
}

func (s *Store) Get(ctx context.Context, k []byte) (v []byte, err error) {
    return s.cli.Get(ctx, s.Key(k))
}

func (s *Store) BatchGet(ctx context.Context, keys [][]byte) ([][]byte, error) {
    return s.cli.BatchGet(ctx, s.Keys(keys))
}

func (s *Store) Delete(ctx context.Context, k []byte) error {
    return s.cli.Delete(ctx, s.Key(k))
}

func (s *Store) BatchDelete(ctx context.Context, keys [][]byte) error {
    return s.cli.BatchDelete(ctx, s.Keys(keys))
}

func (s *Store) Scan(ctx context.Context, startKey, endKey []byte, limit int) (keys [][]byte, values [][]byte, err error) {
    if endKey == nil {
        return s.cli.Scan(ctx, s.Key(startKey), nil, limit)
    }
    return s.cli.Scan(ctx, s.Key(startKey), s.Key(endKey), limit)
}

func (s *Store) DeleteRange(ctx context.Context, startKey []byte, endKey []byte) error {
    return s.cli.DeleteRange(ctx, s.Key(startKey), s.Key(endKey))
}

func (s *Store) Close() error {
    return s.cli.Close()
}
