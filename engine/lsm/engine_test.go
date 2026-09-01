package lsm

import (
	"errors"
	"fmt"

	"os"
	"path/filepath"
	"testing"
	"time"

	errs "github.com/jasutiin/go-database/engine/errors"
)

func TestEngineStartupInitializesStorage(t *testing.T) {
	opts, dbPath := testOptions(t)

	db, err := Startup(opts)
	if err != nil {
		t.Fatalf("Startup() error = %v", err)
	}
	defer db.writeAheadLog.file.Close()

	if db.table == nil {
		t.Fatal("Startup() did not initialize the memtable")
	}

	if db.writeAheadLog == nil {
		t.Fatal("Startup() did not initialize the WAL")
	}

	wantWALPath := filepath.Join(dbPath, "wal.log")
	if db.writeAheadLog.path != wantWALPath {
		t.Fatalf("WAL path = %q, want %q", db.writeAheadLog.path, wantWALPath)
	}

	if _, err := os.Stat(wantWALPath); err != nil {
		t.Fatalf("stat WAL: %v", err)
	}
}

func TestEnginePutAppendsWALAndMemTable(t *testing.T) {
	opts, _ := testOptions(t)
	opts.SkipListMaxLevel = 16
	opts.SkipListMaxSize = 1024

	db, err := Startup(opts)
	if err != nil {
		t.Fatalf("Startup() error = %v", err)
	}
	defer db.writeAheadLog.file.Close()

	if err := db.Put([]byte("name"), []byte("Alice")); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	walEntries, err := db.writeAheadLog.GetEntriesFromWAL()
	if err != nil {
		t.Fatalf("GetEntriesFromWAL() error = %v", err)
	}
	if len(walEntries) != 1 {
		t.Fatalf("WAL entry count = %d, want 1", len(walEntries))
	}
	if string(walEntries[0].key) != "name" || string(walEntries[0].value) != "Alice" {
		t.Fatalf("WAL entry key/value = %q/%q", walEntries[0].key, walEntries[0].value)
	}

	entry, found := db.table.entries.Find("name")
	if !found {
		t.Fatal("memtable entry was not found")
	}
	if string(entry.value) != "Alice" {
		t.Fatalf("memtable value = %q, want %q", entry.value, "Alice")
	}
}

func TestEngineGetReturnsValuesAndHidesTombstones(t *testing.T) {
	opts, _ := testOptions(t)
	db, err := Startup(opts)
	if err != nil {
		t.Fatalf("Startup() error = %v", err)
	}
	defer db.writeAheadLog.file.Close()

	if err := db.Put([]byte("name"), []byte("Alice")); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	value, err := db.Get([]byte("name"))
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if string(value) != "Alice" {
		t.Fatalf("Get() value = %q, want %q", value, "Alice")
	}

	if err := db.Delete([]byte("name")); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = db.Get([]byte("name"))
	if !errors.Is(err, errs.ErrKeyNotFound) {
		t.Fatalf("Get() after Delete() error = %v, want %v", err, errs.ErrKeyNotFound)
	}

	_, err = db.Get([]byte("missing"))
	if !errors.Is(err, errs.ErrKeyNotFound) {
		t.Fatalf("Get() missing key error = %v, want %v", err, errs.ErrKeyNotFound)
	}
}

func testOptions(t *testing.T) (*Options, string) {
	t.Helper()

	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable() error = %v", err)
	}

	dbName := fmt.Sprintf("go-database-test-%d", time.Now().UnixNano())
	dbPath := filepath.Join(filepath.Dir(executable), dbName)
	t.Cleanup(func() {
		if err := os.RemoveAll(dbPath); err != nil {
			t.Errorf("remove test database: %v", err)
		}
	})

	return &Options{
		DbName:           dbName,
		SkipListMaxLevel: 16,
		SkipListMaxSize:  1024,
	}, dbPath
}
