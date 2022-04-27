package redis

import (
    "context"
    "sync"
    "testing"
    "time"
)

var redisConnection1 = []string{"localhost:6379"}
var redisConnection2 = []string{"localhost:7000"}
var redisClusterConnections = []string{"localhost:7001", "localhost:7002", "localhost:7003"}
var invalidRedisConnection = []string{"localhost:7777"}

func TestNewGoRedis(t *testing.T) {
    var cfg Config
    var client *GoRedisClient
    var err error

    cfg.Connections = redisConnection1 // redis addr
    client = NewGoRedis(cfg)
    _, err = client.Do(context.Background(), "ping")
    if err != nil {
        t.Fatal("connect failed")
    }

    cfg.Connections = invalidRedisConnection // invalid redis addr
    client = NewGoRedis(cfg)
    _, err = client.Do(context.Background(), "ping")
    if err == nil {
        t.Fatal("expect connect failed but not")
    }

    cfg.Connections = redisClusterConnections // redis cluster addrs
    client = NewGoRedis(cfg)
    _, err = client.Do(context.Background(), "ping")
    if err != nil {
        t.Fatal("connect failed")
    }
}

type command struct {
    cmd  string
    args []interface{}
}

func TestRunBatch(t *testing.T) {
    var key = "test_key"
    var field1, field2, field3 = "field1", "field2", "field3"
    var value1, value2, value3, value22 = "value1", "value2", "value3", "value22"
    var err error
    var rets []interface{}
    var strs []string
    var strsExpect = []string{"", value22}
    var n int
    var nExpect = 2
    var m map[string]string
    var mExpect = map[string]string{field2: value22, field3: value3}

    cmds := []*command{
        {
            cmd:  "del",
            args: []interface{}{key},
        },
        {
            cmd:  "hmset",
            args: []interface{}{key, field1, value1, field2, value2},
        },
        {
            cmd:  "hdel",
            args: []interface{}{key, field1},
        },
        {
            cmd:  "hmset",
            args: []interface{}{key, field2, value22, field3, value3},
        },
        {
            cmd:  "hmget",
            args: []interface{}{key, field1, field2},
        },
        {
            cmd:  "hlen",
            args: []interface{}{key},
        },
        {
            cmd:  "hgetall",
            args: []interface{}{key},
        },
    }

    goRedisClient := NewGoRedis(Config{Connections: redisConnection1})
    batch := goRedisClient.NewBatch()
    for _, cmd := range cmds {
        if err = batch.Put(cmd.cmd, cmd.args...); err != nil {
            t.Fatal(err)
        }
    }

    if rets, err = goRedisClient.RunBatch(context.Background(), batch); err != nil {
        t.Fatal(err)
    }

    if len(rets) != len(cmds) {
        t.Fatal("length not equal")
    }

    if strs, err = Strings(rets[4], nil); err != nil {
        t.Fatal(err)
    }
    for idx := range strs {
        if strs[idx] != strsExpect[idx] {
            t.Fatal("value not equal")
        }
    }

    if n, err = Int(rets[5], nil); err != nil {
        t.Fatal(err)
    }
    if n != nExpect {
        t.Fatal("value not equal")
    }

    if m, err = StringMap(rets[6], nil); err != nil {
        t.Fatal(err)
    }
    if len(m) != len(mExpect) {
        t.Fatal("length not equal")
    }
    for key := range m {
        if m[key] != mExpect[key] {
            t.Fatal("value not equal")
        }
    }
}

func TestStats(t *testing.T) {
    minIdle := 100
    maxActive := 300
    cfg := Config{
        Connections:          redisConnection1,
        MinIdleConnections:   minIdle,
        MaxActiveConnections: maxActive,
    }
    client := NewGoRedis(cfg)

    time.Sleep(time.Second)
    stats := client.Stats()

    if stats.IdleCount != int64(minIdle) {
        t.Fatal("idle connection number wrong, expect ", minIdle, ", got ", stats.IdleCount)
    }

    if stats.TotalCount != int64(minIdle) {
        t.Fatal("total connection number wrong, expect ", minIdle, ", got ", stats.TotalCount)
    }
}

func TestConfig_Default(t *testing.T) {
    minIdle := 100
    cfg := Config{
        Connections: redisConnection1,
    }
    client := NewGoRedis(cfg)

    time.Sleep(time.Second)
    stats := client.Stats()

    if stats.IdleCount != int64(minIdle) {
        t.Fatal("idle connection number wrong, expect ", minIdle, ", got ", stats.IdleCount)
    }

    if stats.TotalCount != int64(minIdle) {
        t.Fatal("total connection number wrong, expect ", minIdle, ", got ", stats.TotalCount)
    }
}

func TestConfig_IdleTimeout(t *testing.T) {
    var stats *Stats
    minIdle := 10
    maxConn := 20
    concurrency := 30
    idleTimeout := 1 * time.Second
    cfg := Config{
        Connections:          redisConnection1,
        MinIdleConnections:   minIdle,
        IdleCheckFrequency:   time.Second,
        MaxActiveConnections: maxConn,
        IdleTimeout:          idleTimeout,
    }
    client := NewGoRedis(cfg)

    // stats at start
    time.Sleep(time.Second / 2)
    stats = client.Stats()
    if stats.IdleCount != int64(minIdle) {
        t.Fatal("idle connection number wrong, expect ", minIdle, ", got ", stats.IdleCount)
    }
    if stats.TotalCount != int64(minIdle) {
        t.Fatal("total connection number wrong, expect ", minIdle, ", got ", stats.TotalCount)
    }

    // stats after concurrent op
    wg := sync.WaitGroup{}
    wg.Add(concurrency)
    for i := 0; i < concurrency; i++ {
        go func() {
            batch := client.NewBatch()
            _ = batch.Put("dbsize")
            _, _ = client.RunBatch(context.Background(), batch)
            wg.Done()
        }()
    }
    wg.Wait()

    stats = client.Stats()
    if stats.IdleCount != int64(maxConn) {
        t.Fatal("idle connection number wrong, expect ", maxConn, ", got ", stats.IdleCount)
    }
    if stats.TotalCount != int64(maxConn) {
        t.Fatal("total connection number wrong, expect ", maxConn, ", got ", stats.TotalCount)
    }

    // stats after idleTime exceed
    time.Sleep(idleTimeout * 2)
    stats = client.Stats()
    if stats.IdleCount != int64(minIdle) {
        t.Fatal("idle connection number wrong, expect ", minIdle, ", got ", stats.IdleCount)
    }
    if stats.TotalCount != int64(minIdle) {
        t.Fatal("total connection number wrong, expect ", minIdle, ", got ", stats.TotalCount)
    }
}
