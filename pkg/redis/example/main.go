package main

import (
    "context"
    "flag"
    "fmt"
    "github.com/ncraft-io/ncraft-go/pkg/redis"
)

func main() {
    // Update addresses if they have been overwritten by flags
    flag.Parse()

    r := redis.NewRedis()

    batch := r.NewBatch()
    batch.Put("SET", "key.1", "1")
    batch.Put("SET", "key.2", "2")
    batch.Put("SET", "key.3", "3")
    batch.Put("MGET", "key.1", "key.2", "key.3")

    replies, err := r.RunBatch(context.Background(), batch)
    if err != nil {
    }

    r1, err := redis.String(replies[0], nil)
    fmt.Println("first is", r1)

    r2, err := redis.String(replies[1], nil)
    fmt.Println("second is", r2)

    r3, err := redis.String(replies[2], nil)
    fmt.Println("third is ", r3)

    r4, err := redis.Strings(replies[3], nil)
    fmt.Printf("fourth is %v", r4)
}
