package lsm

import (
	"bytes"
	"testing"
)

func TestMemTableInsertStoresValue(t *testing.T) {
	opts, _ := testOptions(t)
	opts.SkipListMaxLevel = 16
	opts.SkipListMaxSize = 1024

	log, err := LoadWAL(opts)
	if err != nil {
		t.Fatalf("LoadWAL() error = %v", err)
	}
	defer log.file.Close()

	table, err := LoadMemTable(opts, log)
	if err != nil {
		t.Fatalf("LoadMemTable() error = %v", err)
	}

	if err := table.Insert([]byte("name"), []byte("Alice"), false); err != nil {
		t.Fatalf("Insert() error = %v", err)
	}

	entry, found := table.entries.Find("name")
	if !found {
		t.Fatal("inserted entry was not found")
	}
	if string(entry.value) != "Alice" {
		t.Fatalf("value = %q, want %q", entry.value, "Alice")
	}
	if entry.tombstone {
		t.Fatal("inserted value is marked as a tombstone")
	}
	if table.Size() != 1 {
		t.Fatalf("size = %d, want 1", table.Size())
	}
}

func TestMemTableGet(t *testing.T) {
	opts, _ := testOptions(t)
	log, err := LoadWAL(opts)
	if err != nil {
		t.Fatalf("LoadWAL() error = %v", err)
	}
	defer log.file.Close()

	table, err := LoadMemTable(opts, log)
	if err != nil {
		t.Fatalf("LoadMemTable() error = %v", err)
	}

	if err := table.Insert([]byte("name"), []byte("Alice"), false); err != nil {
		t.Fatalf("Insert() error = %v", err)
	}
	if err := table.Insert([]byte("deleted"), nil, true); err != nil {
		t.Fatalf("Insert() error = %v", err)
	}

	value, tombstone, found := table.Get([]byte("name"))
	if !found || tombstone || !bytes.Equal(value, []byte("Alice")) {
		t.Fatalf("Get(name) = %q, tombstone=%t, found=%t", value, tombstone, found)
	}

	value[0] = 'a'
	value, tombstone, found = table.Get([]byte("name"))
	if !found || tombstone || !bytes.Equal(value, []byte("Alice")) {
		t.Fatalf("Get(name) after result mutation = %q, tombstone=%t, found=%t", value, tombstone, found)
	}

	value, tombstone, found = table.Get([]byte("deleted"))
	if !found || !tombstone || value != nil {
		t.Fatalf("Get(deleted) = %q, tombstone=%t, found=%t", value, tombstone, found)
	}

	value, tombstone, found = table.Get([]byte("missing"))
	if found || tombstone || value != nil {
		t.Fatalf("Get(missing) = %q, tombstone=%t, found=%t", value, tombstone, found)
	}
}

func TestLoadMemTableInitializesEmptyTable(t *testing.T) {
	opts, _ := testOptions(t)
	opts.SkipListMaxLevel = 16
	opts.SkipListMaxSize = 1024

	log, err := LoadWAL(opts)
	if err != nil {
		t.Fatalf("LoadWAL() error = %v", err)
	}
	defer log.file.Close()

	table, err := LoadMemTable(opts, log)
	if err != nil {
		t.Fatalf("LoadMemTable() error = %v", err)
	}

	if table.entries == nil {
		t.Fatal("LoadMemTable() entries = nil")
	}

	if table.Size() != 0 {
		t.Fatalf("LoadMemTable() size = %d, want 0", table.Size())
	}

	if table.MaxSize() != 1024 {
		t.Fatalf("LoadMemTable() maxSize = %d, want 1024", table.MaxSize())
	}
}
