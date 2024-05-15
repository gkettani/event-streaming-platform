package common

import (
	"sync"
)

type Node struct {
	mu sync.Mutex
	wg sync.WaitGroup

	id      string
	peerIDs []string

	handlers  map[string]func(Message) error
	callbacks map[string]chan Message
}

func NewNode() *Node {
	return &Node{
		handlers:  make(map[string]func(Message) error),
		callbacks: make(map[string]chan Message),
	}
}

func (n *Node) Init(id string, peerIDs []string) {
	n.id = id
	n.peerIDs = peerIDs
}

func (n *Node) ID() string {
	return n.id
}

// TODO: Implement the SyncRPC method
func (n *Node) SyncRPC(any, any, any) (Message, error) {
	return Message{}, nil
}

// Run executes the node's main loop. it waits for incoming messages and
// dispatches them to the appropriate handler.
func (n *Node) Run() error {
	return nil
}
