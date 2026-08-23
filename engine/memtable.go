package engine

import "github.com/jasutiin/go-database/skiplist"

type memTableEntry struct {
	key       string
	value     []byte
	tombstone bool
}

func (entry memTableEntry) compare(other memTableEntry) int {
	switch {
	case entry.key < other.key:
		return -1
	case entry.key > other.key:
		return 1
	default:
		return 0
	}
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

	skip := skiplist.New[memTableEntry](opts.SkipListMaxLevel, memTableEntry.compare)

	for _, value := range entries {
		entry := memTableEntry{
			key:       string(value.key),
			value:     value.value,
			tombstone: false,
		}

		skip.Insert(entry)
	}

	return &memTable{
		entries: skip,
		maxSize: 1024,
	}, nil
}
