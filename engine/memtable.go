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

func LoadMemTable(opts Options) (*memTable, error) {
	return &memTable{
		entries: skiplist.NewSkipList(opts.SkipListMaxLevel, compareEntries),
		maxSize: 1024,
	}, nil
}

func compareEntries(a, b memTableEntry) int {
	if a.key < b.key {
		return -1
	}
	if a.key > b.key {
		return 1
	}
	return 0
}
