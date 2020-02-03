# PubSub
Toy publisher/subscriber service, in Go. Uses websockets to for persistent connections with subscribers.

## Getting started
### Building
```bash
go build -o pubsub main.go
```

### Running
Provide a port as argument, or leave it blank to start the server on an OS-assigned port.
```bash
pubsub 8080
```

The two endpoints are /subscribe, which uses the websocket protocol, and /publish, which uses HTTP and JSON to publish messages to all subscribers.

## Notes
### SubscriberGroup
This type was designed with an Actor-based model in mind, as would be found in Erlang. The struct is nothing more than a container for the Go channels that form the message queue of the Actor, itself a simple goroutine.

### Scalability
Currently, the server uses a single collection of subscribers, leading to potentially degraded performance when publishing to a large number of subscribers. However, this can be alleviated by using multiple SubscriberGroups and keeping their number of subscribers balanced. This would allow to parallelize publishing by using a different goroutine for each group.
