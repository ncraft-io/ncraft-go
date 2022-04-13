package messaging

import (
    "context"
    "github.com/mojo-lang/core/go/pkg/mojo/core"
)

type PushEndpoint struct {
    Service string
    Method  string
    Url     core.Url
}

type Subscription struct {
    Name  string
    Topic string

    PushEndpoint *PushEndpoint

    AutoAck bool
}

type Handler func(ctx context.Context, m *Message) error

type Subscriber interface {
    Subscribe(subscription *Subscription, handler Handler)
}
