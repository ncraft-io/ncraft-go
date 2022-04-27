package redis

import (
    "context"
    "errors"
    redigo "github.com/gomodule/redigo/redis"
)

type Pool struct {
    pool *redigo.Pool
}

type redisCommand struct {
    name string
    args []interface{}
}

type redisBatch struct {
    pool *redigo.Pool
    cmds []redisCommand
}

func NewPool(cfg Config) Redis {
    if cfg.MaxIdleConnections == 0 {
        cfg.MaxIdleConnections = 100
    }

    if cfg.MaxActiveConnections == 0 {
        cfg.MaxActiveConnections = 300
    }

    //if cfg.IdleTimeout == 0 {
    //	cfg.IdleTimeout = 180 * time.Second
    //}

    dialFunc := func() (c redigo.Conn, err error) {
        c, err = redigo.Dial("tcp", cfg.Connections[0])
        if err != nil {
            panic(err)
            return nil, err
        }

        if len(cfg.Password) > 0 {
            if _, err := c.Do("AUTH", cfg.Password); err != nil {
                c.Close()
                return nil, err
            }
        }

        _, err = c.Do("SELECT", cfg.DbNumber)
        if err != nil {
            c.Close()
            return nil, err
        }
        return
    }

    // initialize a new pool
    p := &redigo.Pool{
        MaxIdle:     cfg.MaxIdleConnections,
        IdleTimeout: cfg.IdleTimeout,
        MaxActive:   cfg.MaxActiveConnections,
        Dial:        dialFunc,
    }

    return &Pool{pool: p}
}

func (p *Pool) Do(ctx context.Context, cmd string, arguments ...interface{}) (interface{}, error) {
    connection := p.pool.Get()
    defer connection.Close()

    return connection.Do(cmd, arguments...)
}

func (p *Pool) Close() error {
    return p.pool.Close()
}

func (p *Pool) Stats() *Stats {
    return &Stats{}
}

// NewBatch implement the Cluster NewBatch method.
func (p *Pool) NewBatch() Batch {
    return &redisBatch{
        pool: p.pool,
        cmds: make([]redisCommand, 0),
    }
}

// RunBatch implement the Cluster RunBatch method.
func (p *Pool) RunBatch(ctx context.Context, batch Batch) ([]interface{}, error) {
    bat := batch.(*redisBatch)

    connection := p.pool.Get()
    defer connection.Close()

    for _, cmd := range bat.cmds {
        err := connection.Send(cmd.name, cmd.args...)
        if err != nil {
            return nil, err
        }
    }

    err := connection.Flush()
    if err != nil {
        return nil, err
    }

    var replies []interface{}
    for i := 0; i < len(bat.cmds); i++ {
        reply, err := connection.Receive()
        if err != nil {
            return nil, err
        }

        replies = append(replies, reply)
    }

    return replies, nil
}

func (b *redisBatch) Put(cmd string, args ...interface{}) error {
    if len(args) < 1 {
        return errors.New("no key found in args")
    }

    b.cmds = append(b.cmds, redisCommand{name: cmd, args: args})
    return nil
}
