package pubsub

import (
    "testing"
    "errors"
    //"log"
)

var (
    ch = make(chan Message)
    topic = "TestTopic"
    msID = "TestMsID"
)

func TestCreateTopic(t *testing.T) {
    CreateTopic(topic)

    var exist = TopicExists(topic)
    if !exist {
        t.Errorf("Topic was not created.")
    }
}

func TestDeleteTopic(t *testing.T) {
    DeleteTopic(topic)

    var tExist = TopicExists(topic)
    if tExist {
        t.Errorf("Topic was not deleted.")
    }
}

func TestAddSubscription(t *testing.T)  {
    // Must receive "topic doesn't exist" error
    var targetErr = errors.New("Topic doesn't exist.")
    s, err := AddSubscription(topic, msID, ch)
    t.Run("TestNoTopicExist", func(t *testing.T) {
        if err.Error() != targetErr.Error() {
            t.Errorf("Expected: %q; Got: %q", targetErr, err)
        }
    })

    // Create Topic
    CreateTopic(topic)

    s, err = AddSubscription(topic, msID, ch)
    if err != nil {
        t.Errorf("Error: %q", err)
    }

    t.Run("TestSubscriptionAdded", func(t *testing.T) {
        if _, found := subscriptions[s.id]; !found {
            t.Errorf("Subscription with given ID not added.")
        }
    })

    t.Run("TestSubscriptionMappingAdded", func(t *testing.T) {
        if _, found := subMap[topic][s.id]; !found {
            t.Errorf("Subscription not mapped to topic.")
        }
    })

    t.Run("TestSubscriptionExists", func(t *testing.T) {
        sCopy, _ := AddSubscription(topic, msID, ch)
        if sCopy.id != s.id {
            t.Errorf("Multiple subscriptions for a topic, msID pair.")
        }
    })
}

func TestDeleteSubscription(t *testing.T) {
    //Add subscription
    s, err := AddSubscription(topic, msID, ch)
    if err != nil {
        t.Errorf("Error: %q", err)
    }

    //Delete subscription
    DeleteSubscription(s.GetID())

    t.Run("TestSubscriptionMappingDelete", func(t *testing.T) {
        if _, found := subMap[topic][s.id]; found {
            t.Errorf("Subscription still mapped to topic.")
        }
    })
    t.Run("TestSubscriptionDelete", func(t *testing.T) {
        if _, found := subscriptions[s.id]; found {
            t.Errorf("Subscription not deleted.")
        }
    })
}

func TestUnsubscribe(t *testing.T) {
    //Add subscription
    s, err := AddSubscription(topic, msID, ch)
    if err != nil {
        t.Errorf("Error: %q", err)
    }

    Unsubscribe(s.id)
    if subscriptions[s.id].active {
        t.Errorf("Subscription still active")
    }
}
