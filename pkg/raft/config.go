package raft

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	NodeID            int32    `yaml:"node_id"`
	Peers             []string `yaml:"peers"`
	LogDir            string   `yaml:"log_dir"`
	RPCPort           int      `yaml:"rpc_port"`
	HeartbeatInterval int      `yaml:"heartbeat_interval"`
}

// ClusterConfig represents the configuration for the entire cluster
type ClusterConfig struct {
	Nodes []Config `yaml:"nodes"`
}

// LoadClusterConfig loads the cluster configuration from a file
func LoadClusterConfig(filePath string) (*ClusterConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config ClusterConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
