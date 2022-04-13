package messaging

type Config struct {
    Provider    string `json:"provider"`
    ServiceName string `json:"serviceName"`

    Subscriptions []*Subscription `json:"subscriptions"`
}

func (c *Config) GetSubscriptions() []*Subscription {
    if c != nil {
        return c.Subscriptions
    }
    return nil
}
