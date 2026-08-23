package engine

import (
	"bytes"
	"encoding/binary"
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

func TestGetEntriesFromWALReadsVariableLengthEntries(t *testing.T) {
	opts, _ := testOptions(t)

	log, err := LoadWAL(opts)
	if err != nil {
		t.Fatalf("LoadWAL() error = %v", err)
	}
	defer log.file.Close()

	want := []walEntry{
		{
			kind:  walEntryPut,
			key:   []byte("a"),
			value: []byte("short"),
		},
		{
			kind:  walEntryPut,
			key:   []byte("a-much-longer-key"),
			value: []byte("a value that is longer than eight bytes"),
		},
		{
			kind: walEntryDelete,
			key:  []byte("deleted-key"),
		},
	}

	for _, entry := range want {
		writeTestWALEntry(t, log.file, entry)
	}

	got, err := log.GetEntriesFromWAL()
	if err != nil {
		t.Fatalf("GetEntriesFromWAL() error = %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("entry count = %d, want %d", len(got), len(want))
	}

	for index := range want {
		if got[index].kind != want[index].kind {
			t.Errorf("entry %d kind = %d, want %d", index, got[index].kind, want[index].kind)
		}
		if !bytes.Equal(got[index].key, want[index].key) {
			t.Errorf("entry %d key = %q, want %q", index, got[index].key, want[index].key)
		}
		if !bytes.Equal(got[index].value, want[index].value) {
			t.Errorf("entry %d value = %q, want %q", index, got[index].value, want[index].value)
		}
	}
}

func writeTestWALEntry(t *testing.T, file *os.File, entry walEntry) {
	t.Helper()

	header := make([]byte, walEntryHeaderSize)
	header[0] = byte(entry.kind)
	binary.LittleEndian.PutUint32(header[1:5], uint32(len(entry.key)))
	binary.LittleEndian.PutUint32(header[5:9], uint32(len(entry.value)))

	for _, data := range [][]byte{header, entry.key, entry.value} {
		if _, err := file.Write(data); err != nil {
			t.Fatalf("write test WAL entry: %v", err)
		}
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
