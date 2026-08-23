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

func LoadMemTable(opts *Options, log *wal) (*memTable, error) {
	entries, err := log.GetEntriesFromWAL()

	if err != nil {
		return nil, err
	}

	skip := skiplist.New[memTableEntry](opts.SkipListMaxLevel, compareEntries)

	for _, value := range entries {
		entry := memTableEntry{
			key:       string(value.key),
			value:     value.value,
			tombstone: false,
		}

		skip.Insert(skiplist.NewNode(entry))
	}

	return &memTable{
		entries: skip,
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
