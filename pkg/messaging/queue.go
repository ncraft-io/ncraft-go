package messaging

import "sync"

type Queue interface {
    Publisher
    Subscriber

    Shutdown()
}

type Constructor func(config *Config) (Queue, error)

var queues map[string]Constructor
var queuesOnce sync.Once

func Register(name string, constructor Constructor) {
    queuesOnce.Do(func() {
        queues = make(map[string]Constructor)
    })

    queues[name] = constructor
}
