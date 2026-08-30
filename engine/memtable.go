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

	skip := skiplist.New[memTableEntry](
		opts.SkipListMaxLevel,
		opts.SkipListMaxSize,
		memTableEntry.compare,
	)

	table := &memTable{
		entries: skip,
		maxSize: opts.SkipListMaxSize,
	}

	for _, value := range entries {
		if err := table.Insert(
			value.key,
			value.value,
			value.kind == walEntryDelete,
		); err != nil {
			return nil, err
		}
	}

	return table, nil
}

func (table *memTable) Insert(key, value []byte, tombstone bool) error {
	entry := memTableEntry{
		key:       string(key),
		value:     append([]byte(nil), value...),
		tombstone: tombstone,
	}

	if err := table.entries.Insert(entry); err != nil {
		return err
	}

	table.size = table.entries.Size()
	return nil
}
