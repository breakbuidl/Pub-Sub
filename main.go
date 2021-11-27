package main

import (
    "log"
    "time"
    "sync"

    "pubsubservice/microservice"
    "pubsubservice/pubsub"
)

// Main Method
func main() {
    // Create microservices
    var ms = make([]*microservice.Ms, 0)
    var msNames = []string{"DB", "UI", "Payments", "Analytics"}

    for _, name := range msNames {
        m, err := microservice.CreateMs(name)
        if err != nil {
            log.Println("Error: ", err.Error())
        }
        ms = append(ms, m)
    }

    //Create topics
    var topics = []string{"Topic1", "Topic2", "Topic3", "Topic4", "Topic5", "Topic6"}
    for _, topic := range topics{
        pubsub.CreateTopic(topic)
    }

    //Create subscription map
    var subMap = make(map[*microservice.Ms][]string)
    subMap[ms[0]] = []string{"Topic1", "Topic4", "Topic5"}
    subMap[ms[1]] = []string{"Topic1", "Topic2", "Topic3", "Topic5"}
    subMap[ms[2]] = []string{"Topic4", "Topic5", "Topic6"}
    subMap[ms[3]] = []string{"Topic3", "Topic6"}

    // Add subscriptions
    for m, subT := range subMap {
        for _, t := range subT {
            s, err := pubsub.AddSubscription(t, m.GetID(), m.GetMsgChannel())
            if err != nil {
                log.Println("Error: ", err)
            }
            m.AddSubscription(s)
        }
    }

    var wg1 sync.WaitGroup
    var wg2 sync.WaitGroup

    // Start all microservices
    for _, m := range ms {
        go func(m *microservice.Ms) {
            wg1.Add(1)
            defer wg1.Done()
            m.Start()
        }(m)
    }

    // Wait for all services to start
    time.Sleep(2 * time.Second)

    // Publish messages to topics
    wg2.Add(5)
    go func() {
        defer wg2.Done()
        pubsub.Publish("Topic4", "PubSub")
    }()
    go func() {
        defer wg2.Done()
        pubsub.Publish("Topic2", "is")
    }()
    go func() {
        defer wg2.Done()
        pubsub.Publish("Topic3", "working")
    }()
    go func() {
        defer wg2.Done()
        pubsub.Publish("Topic4", "perfectly")
    }()
    go func() {
        defer wg2.Done()
        pubsub.Publish("Topic6", "fine")
    }()

    wg2.Wait()

    // Stop all microservices
    for _, m := range ms {
        go func(m *microservice.Ms) {
            wg1.Add(1)
            defer wg1.Done()
            m.Stop()
        }(m)
    }

    // Wait for all services to finish
    wg1.Wait()

    //
    for _, m := range ms {
        m.PrintMessages()
    }
}
