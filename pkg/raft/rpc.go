package raft

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/gkettani/event-streaming-platform/pkg/api"
)

type rpcServer struct {
	pb.UnimplementedRaftServer
	node *Node
}

func (s *rpcServer) AppendEntries(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	return s.node.handleAppendEntries(req), nil
}

func (s *rpcServer) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	return s.node.handleRequestVote(req), nil
}

func NewRPCServer(node *Node) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterRaftServer(s, &rpcServer{node: node})
	return s
}

// func StartRPCServer(node *Node, port int) {
// 	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
// 	if err != nil {
// 		log.Fatalf("Failed to listen: %v", err)
// 	}

// 	s := grpc.NewServer()
// 	pb.RegisterRaftServer(s, &rpcServer{node: node})

// 	if err := s.Serve(lis); err != nil {
// 		log.Fatalf("Failed to serve: %v", err)
// 	}
// }
