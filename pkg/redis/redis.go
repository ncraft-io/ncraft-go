package redis

import (
    "context"
    "errors"
    "github.com/ncraft-io/ncraft-go/pkg/config"
)

type Redis interface {
    // Do excute a redis command with random number arguments. First argument will
    // be used as key to hash to a slot, so it only supports a subset of redis
    // commands.
    //
    // SUPPORTED: most commands of keys, strings, lists, sets, sorted sets, hashes.
    // NOT SUPPORTED: scripts, transactions, clusters.
    //
    // Particularly, MSET/MSETNX/MGET are supported using result aggregation.
    // To MSET/MSETNX, there's no atomicity gurantee that given keys are set at once.
    // It's possible that some keys are set, while others not.
    //
    // See README.md for more details.
    // See full redis command list: http://www.redis.io/commands
    Do(ctx context.Context, cmd string, arguments ...interface{}) (interface{}, error)

    Close() error

    // NewBatch create a new redisBatch to pack mutiple commands.
    NewBatch() Batch

    // RunBatch execute commands in redisBatch simutaneously. If multiple commands are
    // directed to the same node, they will be merged and sent at once using pipeling.
    RunBatch(ctx context.Context, batch Batch) ([]interface{}, error)

    Stats() *Stats
}

type Batch interface {
    // Put add a redis command to redisBatch.
    Put(cmd string, arguments ...interface{}) error
}

func NewRedis() Redis {
    var cfg Config

    err := config.Get("redis").Scan(&cfg)
    if err != nil {
        cfg.Connections = []string{":6379"}
    }

    if len(cfg.Connections) == 0 {
        cfg.Connections = []string{":6379"}
    }

    //if cfg.MinIdleConnections == 0 && cfg.MaxIdleConnections != 0 {
    //	cfg.MinIdleConnections = cfg.MaxIdleConnections
    //}

    //if len(cfg.Connections) == 1 { // signal redis mode
    //	return NewPool(cfg)
    //} else { // cluster redis mode
    //	return NewCluster(cfg)
    //}

    return NewGoRedis(cfg)
}

func MGet(redis Redis, keys ...string) (interface{}, error) {
    if redis != nil {
        var args []interface{}
        for _, key := range keys {
            args = append(args, key)
        }
        return redis.Do(context.Background(), "MGET", args...)
    }
    return nil, errors.New("the redis handler is nil")
}
