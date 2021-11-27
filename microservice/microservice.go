package microservice

import (
    "sync"
    "math/rand"
    "encoding/hex"
    "fmt"

    "pubsubservice/pubsub"
)

// Microservice schema
type Ms struct {
    id               string
    name             string
    messages         []pubsub.Message
    msgChannel       chan pubsub.Message
    stopChannel      chan bool
    subscriptions    map[string]pubsub.Subscription      // Maps topics to subscriptions
    lock             sync.RWMutex
}


// CreateMS creates a new Ms
func CreateMs(name string) (*Ms, error) {
    id := make([]byte, 16)
    if _, err := rand.Read(id); err != nil {
        return nil, err
    }

    return &Ms{
        id:             hex.EncodeToString(id),
        name:           name,
        msgChannel:     make(chan pubsub.Message, 100),
        messages:       make([]pubsub.Message, 0),
        stopChannel:    make(chan bool, 1),
        lock:           sync.RWMutex{},
        subscriptions:  make(map[string]pubsub.Subscription),
    }, nil
}

// Start starts the microservice to listen to messages from pubsub
func (ms *Ms) Start() {
    for {
        select {
        case msg := <-ms.msgChannel:
            ms.messages = append(ms.messages, msg)
            pubsub.Ack(msg.GetSubscriptionID(), msg.GetID())
        case <-ms.stopChannel:
            return
        }
    }
}

// Stop stops the microservice
func (ms *Ms) Stop() {
    ms.stopChannel <- true
}

// GetMessages returns a channel of Messages to listem on
func (ms *Ms) PrintMessages() {
    fmt.Println("Microservice: ", ms.name, )
    for _, msg := range ms.messages {
        fmt.Println(msg)
    }
}

// AddSubscription adds the subscription received from pubsub
func (ms *Ms) AddSubscription(s pubsub.Subscription) {
    ms.lock.Lock()
    defer ms.lock.Unlock()
    ms.subscriptions[s.GetTopic()] = s
}

//RemoveSubscription removes the subscription
func (ms *Ms) RemoveSubscription(topic string) {
    ms.lock.Lock()
    delete(ms.subscriptions, topic)
    ms.lock.Unlock()
}

// GetName returns the microservice name
func (ms *Ms) GetName() string {
	return ms.name
}

// GetID returns the microservice ID
func (ms *Ms) GetID() string {
	return ms.id
}

// GetMsgChannel returns a channel of Messages
func (ms *Ms) GetMsgChannel() chan pubsub.Message {
    return ms.msgChannel
}
