package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"

	"github.com/gkettani/event-streaming-platform/pkg/common"
)

func main() {
	server := &Server{
		messages: make(map[int]common.Message),
		offset:   1,
	}

	rpc.Register(server)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Error starting RPC server:", err)
	}
	log.Println("RPC server started on port 1234")
	http.Serve(l, nil)
}

func (s *Server) Send(msg common.Message, reply *common.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Assign offset to the message
	msg.Offset = s.offset
	s.offset++

	// Store message in the log
	s.messages[msg.Offset] = msg

	*reply = common.Message{
		Type:   "send_ack",
		Offset: msg.Offset,
	}

	// Log the sent message
	log.Printf("Sent message: Offset=%d, Type=%s, Key=%s, Msg=%d\n", msg.Offset, msg.Type, msg.Key, msg.Msg)

	return nil
}

func (s *Server) Poll(offset int, reply *common.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the message exists in the log
	msg, ok := s.messages[offset]
	if !ok {
		return rpc.ServerError("Message not found for offset")
	}

	*reply = msg
	return nil
}

type Server struct {
	messages map[int]common.Message // Log of messages
	offset   int                    // Current offset in the log
	mu       sync.Mutex             // Mutex for thread-safe access to messages and offset
}
