package pubsub

import (
	"time"
	"math/rand"
	"encoding/hex"
)
// Message The message metadata
type Message struct {
	id             string
	topic          string
	subscriptionID string
	data           interface{}
	createdAt      int64
}

func CreateMessage(topic, subscriptionID string, data interface {}) (Message, error) {
	id := make([]byte, 16)
    if _, err := rand.Read(id); err != nil {
        return Message{}, err
    }

	return Message{
		id:             hex.EncodeToString(id),
		topic:          topic,
		subscriptionID: subscriptionID,
		data:           data,
		createdAt:      time.Now().UnixNano(),
	}, nil
}

// GetTopic Return the topic of the current message
func (m *Message) GetID() string {
	return m.id
}

// GetSubscriptionID returns the subscription id
func (m *Message) GetSubscriptionID() string {
	return m.subscriptionID
}

// GetTopic Return the topic of the current message
func (m *Message) GetTopic() string {
	return m.topic
}

// GetPayload Get the data of the current message
func (m *Message) GetData() interface{} {
	return m.data
}

// GetCreatedAt Get the creation time of this message
func (m *Message) GetCreatedAt() int64 {
	return m.createdAt
}
