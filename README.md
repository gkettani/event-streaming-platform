## Description
This is a simple event streaming platform implemented using a replicated log service inspired from Kafka. The log service is pretty basic and uses RAFT consensus algorithm for replication. Nodes work in clusters. The communication between nodes is done using gRPC. The goal is to build a distributed event streaming system that is fault-tolerant and highly available. The logs are stored using an on-disk key-value store for persistence.

## Usage
Step1: Install the dependencies
```bash
go mod tidy
```

Step 2: gRPC code generation
```bash
./pkg/api/build.sh
```

Step 3: Define the configuration of the cluster in the `config/config.yaml` file.

Step 4: Start the nodes in the cluster by running the following command on each node:
```bash
go run cmd/main.go
```

Step 5:
Run clients (work in progress)

## Specifications
- Log Replication: The system should replicate the log across all nodes in the cluster to ensure consistency.
- Leader Election: The system should elect a leader node to coordinate log replication and client requests.
- Fault Tolerance: The system should be able to tolerate failures of individual nodes without losing data.
- High Availability: The system should be able to continue operating even if some nodes are unavailable.
- Scalability: The system should be able to scale horizontally by adding more nodes to the cluster.
- Persistence: The log should be stored on disk to ensure durability in case of node failures.
- Consistency: The system should provide strong consistency guarantees for reads and writes.
- Performance: The system should be able to handle a high volume of requests with low latency.
- Client API: The system should provide a simple API for clients to interact with the log service. It should provide a simple SDK.
- Monitoring: The system should provide monitoring and metrics to track the health and performance of the cluster.
- Testing: The system should be thoroughly tested to ensure correctness and reliability. It should include unit tests, integration tests, and stress tests. It should also include fault injection tests to validate the fault tolerance of the system. It should also include performance tests to measure the throughput and latency of the system under different workloads.

### references
1. [RAFT website](https://raft.github.io)
2. [RAFT paper](https://raft.github.io/raft.pdf)
3. [KAFKA paper](https://notes.stephenholiday.com/Kafka.pdf)
