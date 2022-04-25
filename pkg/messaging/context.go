package messaging

import "context"

const (
    TopicKey             = "messaging/topic"
    MessageIdKey         = "messaging/message.id"
    MessageAttributesKey = "messaging/message.attributes"
)

func GetContextTopic(ctx context.Context) string {
    if value, ok := ctx.Value(TopicKey).(string); ok {
        return value
    }
    return ""
}

func GetContextMessageId(ctx context.Context) string {
    if value, ok := ctx.Value(MessageIdKey).(string); ok {
        return value
    }
    return ""
}

func GetContextMessageAttributes(ctx context.Context) map[string]string {
    if value, ok := ctx.Value(MessageAttributesKey).(map[string]string); ok {
        return value
    }
    return nil
}
