package storage

type Config struct {
    Vendor string `json:"vendor" default:"minio"`

    Endpoint  string `json:"endpoint"`
    AccessKey string `json:"AccessKey"`
    SecretKey string `json:"SecretKey"`

    BucketName string `json:"bucketName"`
}
