package minio

import (
    "fmt"
    "github.com/minio/minio-go"
    "github.com/mojo-lang/core/go/pkg/mojo/core"
    "github.com/ncraft-io/ncraft-go/pkg/logs"
    "github.com/ncraft-io/ncraft-go/pkg/storage"
    "github.com/pkg/errors"
)

func init() {
    storage.Register("minio", NewMinio)
}

type Minio struct {
    client     *minio.Client
    bucketName string
}

func NewMinio(cfg *storage.Config) storage.Storage {
    m := &Minio{}

    if client, err := minio.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey, false); err != nil {
        logs.Errorw(fmt.Sprintf("failed to minio connect %s", cfg.Endpoint), "error", err)
        return nil
    } else {
        m.client = client
    }

    if err := m.SetBucket(cfg.BucketName); err != nil {
        logs.Errorw(err.Error())
        return nil
    }

    return m
}

func (m *Minio) BucketName() string {
    return m.bucketName
}

func (m *Minio) SetBucket(name string) error {
    if len(name) > 0 && name != m.bucketName {
        if ok, err := m.client.BucketExists(name); err != nil {
            return err
        } else if !ok {
            if err = m.client.MakeBucket(name, ""); err != nil {
                return errors.Wrap(err, fmt.Sprintf("minio failed to create bucket %s", name))
            }
        }
        m.bucketName = name
    }
    return nil
}

func (m *Minio) Read(key string, options core.Options) (error, *storage.Object) {
    return nil, nil
}

func (m *Minio) Write(object *storage.Object, options core.Options) error {
    return nil
}

func (m *Minio) Download(key string, path string, options core.Options) error {
    if err := m.client.FGetObject(m.bucketName, key, path, minio.GetObjectOptions{}); err != nil {
        return logs.NewErrorw("minio failed to download object", "key", key, "path", path, "error", err.Error())
    }
    return nil
}

func (m *Minio) Upload(localFile string, key string, options core.Options) error {
    if _, err := m.client.FPutObject(m.bucketName, key, localFile, minio.PutObjectOptions{}); err != nil {
        return errors.Wrap(err, fmt.Sprintf("minio failed to upload %s to %s", localFile, key))
    }
    return nil
}
