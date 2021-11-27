package pubsub

import (
    "errors"
    "sync"
    "log"
    "time"
)

// Maps Subscription ID to Subscription
type subscriptionSet map[string]Subscription

// Maps Topic IDs to Subscription Ids
type subscriptionIDSet map[string]bool
type subscriptionMap map[string]subscriptionIDSet

var subMap = make(subscriptionMap)
var subscriptions = make(subscriptionSet)
var sLock sync.RWMutex
var tLock sync.RWMutex

// CreateTopic creates topic
func CreateTopic(topic string) {
    tLock.Lock()
    defer tLock.Unlock()

    // Check if Topic already exists, if not
    // Add key: TopicID, value: make subscriptionSet
    if _, found := subMap[topic]; !found {
        subMap[topic] = make(subscriptionIDSet)
    }
}

// DeleteTopic deletes topic
func DeleteTopic(topic string) {
    tLock.Lock()
    defer tLock.Unlock()

    //Check if Topic exists
    if _, found := subMap[topic]; found {
        delete(subMap, topic)
    }

    // TODO: delete all subscriptions for the given topics
    // TODO: delete topic from all microservices
}

// TopicExists checks if the given topic exists
func TopicExists(topic string) bool {
    tLock.RLock()
    tLock.RUnlock()

    if _, found := subMap[topic]; found {
        return true
    } else {
        return false
    }
}

// AddSubscription creates Subscription for msID and maps it to topic
func AddSubscription(topic, msID string, messages chan Message) (Subscription, error) {
    //Check if topic exists
    if !TopicExists(topic) {
        return Subscription{}, errors.New("Topic doesn't exist.")
    }

    // Check if such subscription already exists
    for sID := range subMap[topic] {
        if subscriptions[sID].msID == msID {
            return subscriptions[sID], nil
        }
    }

    // Create Subscription metadata
    s, err := CreateSubscription(topic, msID, messages)
    if err != nil {
        return Subscription{}, err
    }
    sLock.Lock()
    subscriptions[s.id] = s
    subMap[topic][s.id] = true
    sLock.Unlock()

    return s, nil
}

// DeleteSubscription deletes subscription and its mapping with topic
func DeleteSubscription(subscriptionID string) {
    sLock.Lock()
    defer sLock.Unlock()

    if s, found := subscriptions[subscriptionID]; found {
        delete(subMap[s.topic], subscriptionID)
        delete(subscriptions, subscriptionID)
    }
}

// Unsubscribe deactivates the subscription, doesn't delete it
func Unsubscribe(subscriptionID string) {
    // sLock.Lock()
    // defer sLock.Unlock()

    if s, found := subscriptions[subscriptionID]; found {
        s.SetStatus(false)
        subscriptions[subscriptionID] = s
    }
}

// Publish publishes the data on topic and broadcasts it to the
// msg channels of subscribed microservices
func Publish(topic string, data interface{})  {
    sLock.RLock()
    defer sLock.RUnlock()

    // Broadcast message to every subscription mapped to given topic
    for subscriptionID, _  := range subMap[topic] {
        // Create message
        msg, err := CreateMessage(topic, subscriptionID, data)
        if err != nil {
            log.Println("Error: ", err)
        }

        // Check if the subscription status is active
        if s := subscriptions[subscriptionID]; s.GetStatus() {
            msgAttemptsLoop:
                // Try msgAttempts number of times to send the message
                for i := 0; i < s.msgAttempts; i++ {
                    // Send msg on subscription msgChannel
                    s.msgChannel <- msg
                    // Wait for receive on subscription ackChannel for specified seconds
                    select {
                    // TODO: Check if ack is for the same message that is sent
                    case <-s.ackChannel:
                        break msgAttemptsLoop
                    case <-time.After(time.Duration(s.ackDeadlineSeconds) * time.Second):
                        log.Println("Timeout: Attempt ", i + 1, " failed for messageID: ", msg.GetID())
                    }
                }
        }
    }
}

// Ack is called by the microservice to acknowledge the received message
func Ack(subscriptionID, messageID string) {
    subscriptions[subscriptionID].ackChannel <- messageID
}
