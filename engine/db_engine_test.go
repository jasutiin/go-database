package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
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

	return &Options{DbName: dbName}, dbPath
}
