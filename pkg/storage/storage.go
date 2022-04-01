package storage

import (
    "errors"
    "sync"

    "github.com/mojo-lang/core/go/pkg/mojo/core"
    "github.com/ncraft-io/ncraft-go/pkg/config"
    "github.com/ncraft-io/ncraft-go/pkg/logs"
)

type initializer func(cfg *Config) Storage

var initializers map[string]initializer

func Register(name string, init initializer) {
    (&sync.Once{}).Do(func() {
        initializers = make(map[string]initializer)
    })

    initializers[name] = init
}

type Storage interface {
    BucketName() string
    SetBucket(name string) error

    Read(key string, options core.Options) (error, *Object)
    Write(object *Object, options core.Options) error

    Download(key string, path string, options core.Options) error
    Upload(localFile string, key string, options core.Options) error
}

func NewStorage(cfg *Config) Storage {
    if init, ok := initializers[cfg.Vendor]; ok {
        return init(cfg)
    }
    return nil
}

var storage Storage

func GetStorage() Storage {
    (&sync.Once{}).Do(func() {
        conf := &Config{}
        if err := config.ScanKey("storage", conf); err != nil {
            logs.Warnw("failed to get the server config", "error", err.Error())
            storage = NewDummyStorage()
        } else {
            if storage = NewStorage(conf); storage == nil {
                storage = NewDummyStorage()
            }
        }
    })
    return storage
}

func Download(key string, path string, options core.Options) error {
    return GetStorage().Download(key, path, options)
}
func Upload(localFile string, key string, options core.Options) error {
    return GetStorage().Upload(localFile, key, options)
}

type DummyStorage struct {
    err error
}

func NewDummyStorage() Storage                                                          { return &DummyStorage{err: errors.New("DummyStorage: not implement")} }
func (s *DummyStorage) BucketName() string                                              { return "dummy" }
func (s *DummyStorage) SetBucket(name string) error                                     { return s.err }
func (s *DummyStorage) Read(key string, options core.Options) (error, *Object)          { return s.err, nil }
func (s *DummyStorage) Write(object *Object, options core.Options) error                { return s.err }
func (s *DummyStorage) Download(key string, path string, options core.Options) error    { return s.err }
func (s *DummyStorage) Upload(localFile string, key string, options core.Options) error { return s.err }
