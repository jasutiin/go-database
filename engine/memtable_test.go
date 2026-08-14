package engine

import "testing"

func TestLoadMemTableInitializesEmptyTable(t *testing.T) {
	table, err := LoadMemTable()
	if err != nil {
		t.Fatalf("LoadMemTable() error = %v", err)
	}

	if table.entries == nil {
		t.Fatal("LoadMemTable() entries = nil")
	}

	if table.size != 0 {
		t.Fatalf("LoadMemTable() size = %d, want 0", table.size)
	}

	if table.maxSize != 1024 {
		t.Fatalf("LoadMemTable() maxSize = %d, want 1024", table.maxSize)
	}
}
