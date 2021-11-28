# Introduction

The module pubsubservice in this project replicates the Pub/Sub service that facilitates the communication among different microservices. It does so 
with the help of two packages viz. microservice and pubsub.

- microservice package serves as the client of the pubsub system in form of Subscriber and Publisher.
- pubsub package provides the core functionality of the system with Message, Subscription and Topic as key components.

# Setup

1. Download the repository.
2. `go build` at the root of the directory
3. `./pubsubservice` to run the program
4. `go test -v ./...` to run the tests

# Details

The design of the components is such that to keep them as loosely coupled as possible to simulate different services in different containers.

**microservice** package provides an Ms struct for microservice schema and its methods.
- **Ms.Start()** starts the service to listen on channel to recieve messages
- **Ms.Stop()** stops the service and no further messages can be sent

Here's how a typical functioning goes: (Subscriptions and Topics are created beforehand)
- `Ms` publishes messages to topic
- `pubsub` broadcasts the message on all the subscribed microservices' `msgChannel`
- If the subscribed microservice recieves the message, it sends an ack on `ackChannel`. If `pubsub` doesn't recieve ack within `ackDeadlineSeconds`,
it tries for `msgAttemps` number of times. (these configs are set in subscription)

