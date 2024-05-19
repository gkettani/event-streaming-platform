package raft

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

type Storage struct {
	db *leveldb.DB
}

func NewStorage(path string) *Storage {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Fatalf("Failed to open LevelDB: %v", err)
	}
	return &Storage{db: db}
}

func (s *Storage) SaveLogEntry(entry LogEntry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return s.db.Put([]byte(fmt.Sprintf("log-%d", entry.Index)), data, nil)
}

func (s *Storage) GetLogEntry(index int) (LogEntry, error) {
	data, err := s.db.Get([]byte(fmt.Sprintf("log-%d", index)), nil)
	if err != nil {
		return LogEntry{}, err
	}

	var entry LogEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return LogEntry{}, err
	}
	return entry, nil
}
