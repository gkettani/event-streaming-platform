package raft

import (
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type Cluster struct {
	// nodes   []*Node
	servers []*grpc.Server
	wg      sync.WaitGroup
}

func NewCluster() *Cluster {
	return &Cluster{}
}

func (nm *Cluster) AddNode(config *Config) {
	nm.wg.Add(1)
	go func() {
		defer nm.wg.Done()
		node := NewNode(config)
		// nm.nodes = append(nm.nodes, node)
		if err := node.Start(); err != nil {
			log.Fatalf("Failed to start node %d: %v", config.NodeID, err)
		}

		server := NewRPCServer(node)
		nm.servers = append(nm.servers, server)
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.RPCPort))
		if err != nil {
			log.Fatalf("Failed to listen on port %d: %v", config.RPCPort, err)
		}
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server on port %d: %v", config.RPCPort, err)
		}
	}()
}

func (c *Cluster) Wait() {
	c.wg.Wait()
}

func (c *Cluster) Stop() {
	for _, server := range c.servers {
		server.GracefulStop()
	}
}
