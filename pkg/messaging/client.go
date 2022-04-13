package messaging

import (
    "context"
    "errors"
    "github.com/ncraft-io/ncraft-go/pkg/config"
)

type Client struct {
    Config   *Config
    Provider Queue
}

func NewClient() (*Client, error) {
    cfg := &Config{}
    if err := config.ScanKey("messaging", cfg); err != nil {
        return nil, err
    }
    return NewClientWith(cfg)
}

func NewClientWith(config *Config) (*Client, error) {
    if config == nil {
    }

    if queues != nil {
        if c, ok := queues[config.Provider]; ok {
            if q, err := c(config); err != nil {
                return nil, err
            } else {
                return &Client{Provider: q}, nil
            }
        }
    }

    return nil, errors.New("")
}

func (c *Client) GetConfig() *Config {
    if c != nil {
        return c.Config
    }
    return nil
}

func (c *Client) Shutdown() {
    if c != nil {
        c.Provider.Shutdown()
    }
}

func (c *Client) Subscribe(subscription *Subscription, handler Handler) {
    if c != nil {
        c.Provider.Subscribe(subscription, handler)
    }
}

func (c *Client) Publish(ctx context.Context, topic string, messages ...*Message) error {
    if c != nil {
        c.Provider.Publish(ctx, topic, messages...)
    }

    return nil
}
