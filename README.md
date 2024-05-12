## Description
This is a simple example of a event streaming system implemented in Go. The idea is to build a distributed version using multi-nodes and leverage the power of Go's concurrency model. RPC is used to communicate between nodes. 

## Usage
Run server
```bash
go run cmd/server/main.go
```

Run producer
```bash
go run cmd/producer/main.go
```

Run consumer
```bash
go run cmd/consumer/main.go
```

## Specifications:
1. **Communication Protocol:** The application uses RPC (Remote Procedure Call) over TCP for communication between the server, producer, and consumer.

2. **Message Format:** Messages exchanged between components are encoded in JSON format.

3. **Server:** The server maintains a log of messages and provides two RPC methods:
- *Send:* Accepts a message from the producer, assigns it an offset, stores it in the log, and returns an acknowledgment.
- *Poll:* Accepts an offset from the consumer, retrieves the corresponding message from the log, and returns it.

4. **Producer:** The producer periodically sends messages to the server using the Send RPC method.

5. **Consumer:** The consumer periodically polls messages from the server using the Poll RPC method.

6. **Log:** The log is a simple in-memory data structure that stores messages in the order they were received. Each message is assigned an offset that corresponds to its position in the log.

## Features
1. Message Sending: The producer can send messages to the server, which are then stored in the log.

2. Message Retrieval: The consumer can retrieve messages from the server using the message offset.

3. Offset Management: The server assigns a unique offset to each message it receives, allowing consumers to retrieve messages by their offset.

4. Concurrency: The server handles concurrent access to the message log using a mutex to ensure thread safety.

5. Error Handling: The application handles errors gracefully, providing appropriate error messages and responses when encountered.

## Checklist
- [x] Implement a simple event streaming system
- [x] Implement a simple producer
- [x] Implement a simple consumer
- [ ] Use goroutines for concurrency
- [ ] Distributed version
- [ ] Partitioning
- [ ] Fault tolerance
- [ ] High availability
