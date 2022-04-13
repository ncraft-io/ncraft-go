package messaging

import "time"

type Message struct {
    Id         string            `json:"id"`
    Attributes map[string]string `json:"attributes"`
    Data       []byte            `json:"data"`

    didAck bool
    ack    func()
}

func NewMessage(data interface{}, attributes ...string) *Message {
    return nil
}

func NewStringMessage(data string, attributes ...string) *Message {
    return &Message{
        Id:   "",
        Data: nil,
    }
}

func (m *Message) SetAck(ack func()) *Message {
    if m != nil {
        m.ack = ack
    }
    return m
}

func (m *Message) Ack() {
    if m != nil && !m.didAck {
        m.ack()
        m.didAck = true
    }
}

func (m *Message) SetAttribute(key string, value string) *Message {
    if m != nil {
        if m.Attributes == nil {
            m.Attributes = make(map[string]string)
        }
        m.Attributes[key] = value
    }
    return m
}

type ReceivedMessage struct {
    Message

    AckId        string
    Subscription string
    PublishTime  time.Time
}
