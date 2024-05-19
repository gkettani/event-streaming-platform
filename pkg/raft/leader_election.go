package raft

import pb "github.com/gkettani/event-streaming-platform/pkg/api"

func (n *Node) handleRequestVote(request *pb.RequestVoteRequest) *pb.RequestVoteResponse {
	n.mu.Lock()
	defer n.mu.Unlock()

	response := &pb.RequestVoteResponse{
		Term:        n.currentTerm,
		VoteGranted: false,
	}

	if request.Term < n.currentTerm {
		return response
	}

	if request.Term > n.currentTerm {
		n.currentTerm = request.Term
		n.votedFor = -1
		n.state = Follower
	}

	if (n.votedFor == -1 || n.votedFor == request.CandidateId) &&
		(len(n.log) == 0 || (request.LastLogTerm > n.log[len(n.log)-1].Term ||
			(request.LastLogTerm == n.log[len(n.log)-1].Term && request.LastLogIndex >= n.log[len(n.log)-1].Index))) {
		response.VoteGranted = true
		n.votedFor = request.CandidateId
		n.electionTimeout.Reset(randomElectionTimeout())
	}

	return response
}

func (n *Node) handleAppendEntries(request *pb.AppendEntriesRequest) *pb.AppendEntriesResponse {
	n.mu.Lock()
	defer n.mu.Unlock()

	response := &pb.AppendEntriesResponse{
		Term:    n.currentTerm,
		Success: false,
	}

	if request.Term < n.currentTerm {
		return response
	}

	n.electionTimeout.Reset(randomElectionTimeout())

	if request.Term > n.currentTerm {
		n.currentTerm = request.Term
		n.state = Follower
	}

	// Ensure log consistency
	if len(n.log) < int(request.PrevLogIndex) || (request.PrevLogIndex > 0 && n.log[request.PrevLogIndex-1].Term != request.PrevLogTerm) {
		return response
	}

	n.log = append(n.log[:request.PrevLogIndex], fromProtoLogEntries(request.Entries)...)
	if request.LeaderCommit > n.commitIndex {
		n.commitIndex = min(request.LeaderCommit, int32(len(n.log)-1))
	}

	response.Success = true
	return response
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
