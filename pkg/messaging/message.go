package messaging

import "time"

type Message struct {
    Id         string
    Attributes map[string]string
    Data       []byte
}

type ReceivedMessage struct {
    Message

    AckId        string
    Subscription string
    PublishTime  time.Time
}
