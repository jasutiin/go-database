package lsm

import (
	"strings"

	"github.com/jasutiin/go-database/storage/skiplist"
)

type memTableEntry struct {
	value     []byte
	tombstone bool
}

type memTable struct {
	entries *skiplist.SkipList[string, memTableEntry]
}

func LoadMemTable(opts *Options, log *wal) (*memTable, error) {
	entries, err := log.GetEntriesFromWAL()

	if err != nil {
		return nil, err
	}

	skip := skiplist.New[string, memTableEntry](
		opts.SkipListMaxLevel,
		opts.SkipListMaxSize,
		strings.Compare,
	)

	table := &memTable{
		entries: skip,
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
		value:     append([]byte(nil), value...),
		tombstone: tombstone,
	}

	return table.entries.Insert(string(key), entry)
}

func (table *memTable) Get(key []byte) (value []byte, tombstone bool, found bool) {
	entry, found := table.entries.Find(string(key))
	if !found {
		return nil, false, false
	}

	return append([]byte(nil), entry.value...), entry.tombstone, true
}

func (table *memTable) Size() int {
	return table.entries.Size()
}

func (table *memTable) MaxSize() int {
	return table.entries.MaxSize()
}
