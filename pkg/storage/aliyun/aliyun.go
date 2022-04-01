package aliyun

import (
    "fmt"
    "github.com/mojo-lang/core/go/pkg/mojo/core"
    "github.com/ncraft-io/ncraft-go/pkg/logs"
    "github.com/ncraft-io/ncraft-go/pkg/storage"

    "github.com/aliyun/aliyun-oss-go-sdk/oss"
    "github.com/pkg/errors"
)

func init() {
    storage.Register("aliyun", NewAliyun)
}

type Aliyun struct {
    client *oss.Client
    bucket *oss.Bucket
}

func NewAliyun(cfg *storage.Config) storage.Storage {
    m := &Aliyun{}

    if client, err := oss.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey); err != nil {
        logs.Errorw(fmt.Sprintf("failed to aliyun oss connect %s", cfg.Endpoint), "error", err)
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

func (a *Aliyun) BucketName() string {
    if a != nil && a.bucket != nil {
        return a.bucket.BucketName
    }
    return ""
}

func (a *Aliyun) SetBucket(name string) error {
    if len(name) > 0 && (a.bucket == nil || name != a.bucket.BucketName) {
        if ok, err := a.client.IsBucketExist(name); err != nil {
            return err
        } else if !ok {
            if err = a.client.CreateBucket(name); err != nil {
                return errors.Wrap(err, fmt.Sprintf("minio failed to create bucket %s", name))
            }
        }
        if bucket, err := a.client.Bucket(name); err != nil {
            return err
        } else {
            a.bucket = bucket
        }
    }
    return nil
}

func (a *Aliyun) Read(key string, options core.Options) (error, *storage.Object) {
    return nil, nil
}

func (a *Aliyun) Write(object *storage.Object, options core.Options) error {
    return nil
}

func (a *Aliyun) Download(key string, path string, options core.Options) error {
    if a.bucket != nil {
        if err := a.bucket.GetObjectToFile(key, path); err != nil {
            return logs.NewErrorw(fmt.Sprintf("aliyun oss failed to download %s to %s", key, path), "error", err)
        }
    } else {
        return logs.NewError("should set proper bucket name before download storage object")
    }

    return nil
}

func (a *Aliyun) Upload(localFile string, key string, options core.Options) error {
    if a.bucket != nil {
        if err := a.bucket.PutObjectFromFile(key, localFile); err != nil {
            return logs.NewErrorw(fmt.Sprintf("aliyun oss failed to upload %s to %s", localFile, key), "error", err)
        }
    } else {
        return logs.NewError("should set proper bucket name before upload storage object")
    }

    return nil
}
