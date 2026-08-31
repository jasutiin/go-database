package engine

import "testing"

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

	entry, found := table.entries.Find(memTableEntry{key: "name"})
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
