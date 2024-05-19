package raft

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	pb "github.com/gkettani/event-streaming-platform/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type State int

const (
	Follower State = iota
	Candidate
	Leader
)

type Node struct {
	mu          sync.Mutex
	state       State
	currentTerm int32
	votedFor    int32
	log         []LogEntry
	commitIndex int32
	// lastApplied       int32
	nextIndex         map[string]int // map of peers to the index of the next log entry to send to that peer
	matchIndex        map[string]int // map of peers to the index of the highest log entry known to be replicated on that peer
	peers             []string
	id                int32
	electionTimeout   *time.Timer
	heartbeatInterval time.Duration
	applyCh           chan ApplyMsg // channel to send committed log entries to the client
	storage           *Storage
}

type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

func NewNode(config *Config) *Node {
	return &Node{
		state:             Follower,
		id:                config.NodeID,
		peers:             config.Peers,
		nextIndex:         make(map[string]int),
		matchIndex:        make(map[string]int),
		log:               make([]LogEntry, 0),
		applyCh:           make(chan ApplyMsg),
		storage:           NewStorage(config.LogDir),
		heartbeatInterval: time.Duration(config.HeartbeatInterval) * time.Millisecond,
		electionTimeout:   time.NewTimer(randomElectionTimeout()),
	}
}

func (n *Node) Start() error {
	go n.run()
	return nil
}

func (n *Node) run() {
	for {
		switch n.state {
		case Follower:
			n.runFollower()
		case Candidate:
			n.runCandidate()
		case Leader:
			n.runLeader()
		}
	}
}

func (n *Node) runFollower() {
	log.Printf("Node %d is a follower", n.id)
	n.electionTimeout.Reset(randomElectionTimeout())
	for n.state == Follower {
		<-n.electionTimeout.C
		n.state = Candidate
	}
}

func (n *Node) runCandidate() {
	log.Printf("Node %d is a candidate", n.id)
	n.startElection()
	for n.state == Candidate {
		<-n.electionTimeout.C
		n.startElection()
	}
}

func (n *Node) runLeader() {
	log.Printf("Node %d is a leader", n.id)
	n.electionTimeout.Stop()
	ticker := time.NewTicker(n.heartbeatInterval)
	defer ticker.Stop()

	// We send periodic heartbeats to all peers as long as we are the leader
	for n.state == Leader {
		<-ticker.C
		n.sendHeartbeats()
	}
}

// Inititated by candidates to gather votes whenever a leader goes down
func (n *Node) startElection() {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.state = Candidate
	n.currentTerm++
	n.votedFor = n.id

	votes := 1 // vote for self

	for _, peer := range n.peers {
		go func(peer string) {
			var lastLogIndex int32 = 0
			var lastLogTerm int32 = 0
			if len(n.log) > 0 {
				lastLogIndex = n.log[len(n.log)-1].Index
				lastLogTerm = n.log[len(n.log)-1].Term
			}
			request := &pb.RequestVoteRequest{
				Term:         n.currentTerm,
				CandidateId:  n.id,
				LastLogIndex: lastLogIndex,
				LastLogTerm:  lastLogTerm,
			}
			// Send RequestVote RPC to peer to see if it votes for us
			response := n.sendRequestVote(peer, request)
			if response != nil && response.VoteGranted {
				n.mu.Lock()
				votes++
				if votes > len(n.peers)/2 && n.state == Candidate {
					n.state = Leader
					n.mu.Unlock()
					return
				}
				n.mu.Unlock()
			}
		}(peer)
	}

	n.electionTimeout.Reset(randomElectionTimeout())
}

func (n *Node) sendHeartbeats() {
	for _, peer := range n.peers {
		go n.sendAppendEntries(peer)
	}
}

// So far the idea is to send a empty entry to all peers
// to notify them that we are still the leader
// This is only called by the leader
func (n *Node) sendAppendEntries(peer string) {
	log.Printf("Sending AppendEntries to peer %s", peer)
	n.mu.Lock()
	// define what logs should be sent to that peer
	prevLogIndex := n.nextIndex[peer] - 1
	prevLogTerm := int32(0)
	if prevLogIndex >= 0 {
		prevLogTerm = n.log[prevLogIndex].Term
	}
	entries := n.log[n.nextIndex[peer]:]
	n.mu.Unlock()

	request := &pb.AppendEntriesRequest{
		Term:         n.currentTerm,
		LeaderId:     n.id,
		PrevLogIndex: int32(prevLogIndex),
		PrevLogTerm:  prevLogTerm,
		Entries:      toProtoLogEntries(entries),
		LeaderCommit: int32(n.commitIndex),
	}

	conn, err := grpc.NewClient(peer, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		log.Printf("Failed to connect to peer %s: %v", peer, err)
		return
	}
	defer conn.Close()

	client := pb.NewRaftClient(conn)
	response, err := client.AppendEntries(context.Background(), &pb.AppendEntriesRequest{
		Term:         request.Term,
		LeaderId:     request.LeaderId,
		PrevLogIndex: request.PrevLogIndex,
		PrevLogTerm:  request.PrevLogTerm,
		Entries:      request.Entries,
		LeaderCommit: request.LeaderCommit,
	})

	if err != nil {
		log.Printf("Failed to send AppendEntries to peer %s: %v", peer, err)
		return
	}

	if response.Term > n.currentTerm {
		n.mu.Lock()
		n.currentTerm = response.Term
		n.state = Follower
		n.votedFor = -1
		n.mu.Unlock()
		return
	}

	if response.Success {
		// Keep track that this peer received the log entries
		// update nextIndex and matchIndex
		log.Printf("Peer %s received log entries", peer)
	} else {
		// Something is wrong
		// maybe we'll need to send previous log entries??
		log.Printf("Peer %s did not receive log entries", peer)
	}
}

func randomElectionTimeout() time.Duration {
	return time.Duration(150+rand.Intn(150)) * time.Millisecond
}

// Convert from protobuf LogEntry to internal LogEntry
func fromProtoLogEntry(entry *pb.LogEntry) LogEntry {
	return LogEntry{
		Term:    entry.Term,
		Index:   entry.Index,
		Command: entry.Command,
	}
}

// Convert from internal LogEntry to protobuf LogEntry
func toProtoLogEntry(entry LogEntry) *pb.LogEntry {
	return &pb.LogEntry{
		Term:    entry.Term,
		Index:   entry.Index,
		Command: entry.Command,
	}
}

// Convert slice of protobuf LogEntry to slice of internal LogEntry
func fromProtoLogEntries(entries []*pb.LogEntry) []LogEntry {
	logEntries := make([]LogEntry, len(entries))
	for i, entry := range entries {
		logEntries[i] = fromProtoLogEntry(entry)
	}
	return logEntries
}

// Convert slice of internal LogEntry to slice of protobuf LogEntry
func toProtoLogEntries(entries []LogEntry) []*pb.LogEntry {
	logEntries := make([]*pb.LogEntry, len(entries))
	for i, entry := range entries {
		logEntries[i] = toProtoLogEntry(entry)
	}
	return logEntries
}

// sendRequestVote sends a RequestVote RPC to a peer
func (n *Node) sendRequestVote(peer string, request *pb.RequestVoteRequest) *pb.RequestVoteResponse {
	log.Printf("Node %d sending RequestVote to peer %s", n.id, peer)
	// Implementation of sending RequestVote RPC to a peer
	// This would involve setting up a gRPC client and making the request
	return &pb.RequestVoteResponse{} // Placeholder return
}
