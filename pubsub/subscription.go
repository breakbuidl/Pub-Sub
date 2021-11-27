package pubsub

import (
	"encoding/hex"
	"math/rand"
    "time"
	"sync"
)

// Subscription metadata
type Subscription struct {
    id                 string
    topic              string
    msID               string
    msgChannel         chan Message
	ackChannel         chan string
	ackDeadlineSeconds int
	msgAttempts        int
    active             bool
    createdAt          int64
	sLock              sync.RWMutex
}

func CreateSubscription(topic, msID string, messages chan Message) (Subscription, error) {
    id := make([]byte, 16)
	if _, err := rand.Read(id); err != nil {
		return Subscription{}, err
	}

    return Subscription{
        id:                 hex.EncodeToString(id),
        topic:              topic,
        msID:               msID,
        active:             true,
        createdAt:          time.Now().UnixNano(),
        msgChannel:         messages,
		ackChannel:         make(chan string),
		ackDeadlineSeconds: 3,
		msgAttempts:        3,
    }, nil
}

func (s *Subscription) GetTopic() string {
    return s.topic
}

func (s *Subscription) GetID() string {
    return s.id
}

func (s *Subscription) GetStatus() bool {
	sLock.RLock()
	defer sLock.RUnlock()
	return s.active
}

func (s *Subscription) SetStatus(status bool) {
	sLock.Lock()
	defer sLock.Unlock()
	s.active = status
}
