package redis

import (
    "time"
)

type ClusterSlot struct {
    Start int           `json:"start"`
    End   int           `json:"end"`
    Nodes []ClusterNode `json:"nodes"`
}

type ClusterNode struct {
    Id   string `json:"id"`
    Addr string `json:"addr"`
}

type Config struct {
    Connections []string `json:"connections"`
    //ClusterSlots []redis.ClusterSlot `json:"clusterSlots"`
    Password string `json:"password"`

    DbNumber int `json:"dbNumber"`

    // Maximum number of idle connections in the pool.
    // Deprecate in future version, using MinIdleConnections instead
    MaxIdleConnections int `json:"maxIdleConnections"`
    // Minimum number of idle connections in the pool.
    // Default is 100 connections
    MinIdleConnections int `json:"minIdleConnections"`
    // Maximum number of socket connections.
    // Default is 300 connections
    MaxActiveConnections int `json:"maxActiveConnections"`
    // Close connections after remaining idle for this duration. If the value
    // is zero, then idle connections are not closed. Applications should set
    // the timeout to a value less than the server's timeout.
    IdleTimeout time.Duration `json:"idleTimeout"`
    // Frequency of idle checks made by idle connections reaper.
    // Default is 1 minute. -1 disables idle connections reaper,
    // but idle connections are still discarded by the client
    // when getting connection if IdleTimeout is set.
    IdleCheckFrequency time.Duration

    // Connection age at which client retires (closes) the connection.
    // Default is to not close aged connections.
    ConnectionTimeout time.Duration `json:"connectionTimeout"`
    // Timeout for socket reads. If reached, commands will fail
    // with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
    // Default is 1 seconds.
    ReadTimeout time.Duration `json:"readTimeout"`
    // Timeout for socket writes. If reached, commands will fail
    // with a timeout instead of blocking.
    // Default is ReadTimeout.
    WriteTimeout time.Duration `json:"writeTimeout"`

    // Deprecated
    KeepAlive int `json:"keepAlive"`
    // Deprecated
    AliveTime time.Duration `json:"aliveTime"`

    // Maximum number of retries before giving up.
    // Default is to not retry failed commands.
    MaxRetries int `json:"maxRetries"`
    // Minimum backoff between each retry.
    // Default is 8 milliseconds; -1 disables backoff.
    MinRetryBackoff time.Duration `json:"minRetryBackoff"`
    // Maximum backoff between each retry.
    // Default is 512 milliseconds; -1 disables backoff.
    MaxRetryBackoff time.Duration `json:"maxRetryBackoff"`
}
