package engine

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadWALCreatesFile(t *testing.T) {
	opts, dbPath := testOptions(t)

	log, err := LoadWAL(opts)
	if err != nil {
		t.Fatalf("LoadWAL() error = %v", err)
	}
	defer log.file.Close()

	wantPath := filepath.Join(dbPath, "wal.log")
	if log.path != wantPath {
		t.Fatalf("WAL path = %q, want %q", log.path, wantPath)
	}

	info, err := os.Stat(wantPath)
	if err != nil {
		t.Fatalf("stat WAL: %v", err)
	}

	if info.IsDir() {
		t.Fatalf("WAL path %q is a directory", wantPath)
	}
}

func TestLoadWALReopensExistingFileWithoutTruncatingIt(t *testing.T) {
	opts, _ := testOptions(t)
	want := []byte("existing WAL data")

	first, err := LoadWAL(opts)
	if err != nil {
		t.Fatalf("first LoadWAL() error = %v", err)
	}

	if _, err := first.file.Write(want); err != nil {
		first.file.Close()
		t.Fatalf("write WAL: %v", err)
	}

	if err := first.file.Close(); err != nil {
		t.Fatalf("close first WAL: %v", err)
	}

	second, err := LoadWAL(opts)
	if err != nil {
		t.Fatalf("second LoadWAL() error = %v", err)
	}
	defer second.file.Close()

	got := make([]byte, len(want))
	if _, err := second.file.ReadAt(got, 0); err != nil {
		t.Fatalf("read reopened WAL: %v", err)
	}

	if !bytes.Equal(got, want) {
		t.Fatalf("reopened WAL data = %q, want %q", got, want)
	}
}
