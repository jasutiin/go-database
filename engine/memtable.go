package engine

import "github.com/jasutiin/go-database/skiplist"

type memTableEntry struct {
	key       string
	value     []byte
	tombstone bool
}

type memTable struct {
	entries *skiplist.SkipList[memTableEntry]
	size    int
	maxSize int
}

func LoadMemTable() (*memTable, error) {
	return &memTable{
		entries: new(skiplist.SkipList[memTableEntry]),
		maxSize: 1024,
	}, nil
}
