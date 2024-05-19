package raft

type LogEntry struct {
	Index   int32
	Term    int32
	Command string
}

func (n *Node) appendEntry(entry LogEntry) {
	n.log = append(n.log, entry)
}
