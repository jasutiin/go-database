package engine

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"os"
	"path/filepath"
	"strings"
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

func TestWALInsertAppendsEntries(t *testing.T) {
	opts, _ := testOptions(t)

	log, err := LoadWAL(opts)
	if err != nil {
		t.Fatalf("LoadWAL() error = %v", err)
	}
	defer log.file.Close()

	if err := log.Insert([]byte("name"), []byte("Alice"), false); err != nil {
		t.Fatalf("Insert() put error = %v", err)
	}
	if err := log.Insert([]byte("deleted"), nil, true); err != nil {
		t.Fatalf("Insert() delete error = %v", err)
	}

	entries, err := log.GetEntriesFromWAL()
	if err != nil {
		t.Fatalf("GetEntriesFromWAL() error = %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("entry count = %d, want 2", len(entries))
	}

	if entries[0].kind != walEntryPut ||
		!bytes.Equal(entries[0].key, []byte("name")) ||
		!bytes.Equal(entries[0].value, []byte("Alice")) {
		t.Fatalf("put entry = %+v", entries[0])
	}

	if entries[1].kind != walEntryDelete ||
		!bytes.Equal(entries[1].key, []byte("deleted")) ||
		len(entries[1].value) != 0 {
		t.Fatalf("delete entry = %+v", entries[1])
	}

	if log.size <= 0 {
		t.Fatalf("WAL size = %d, want positive size", log.size)
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
			logNumber: 1,
			kind:      walEntryPut,
			key:       []byte("a"),
			value:     []byte("short"),
		},
		{
			logNumber: 2,
			kind:      walEntryPut,
			key:       []byte("a-much-longer-key"),
			value:     []byte("a value that is longer than eight bytes"),
		},
		{
			logNumber: 3,
			kind:      walEntryDelete,
			key:       []byte("deleted-key"),
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
		if got[index].logNumber != want[index].logNumber {
			t.Errorf("entry %d log number = %d, want %d", index, got[index].logNumber, want[index].logNumber)
		}
		if got[index].kind != want[index].kind {
			t.Errorf("entry %d kind = %d, want %d", index, got[index].kind, want[index].kind)
		}
		if got[index].keyLength != uint16(len(want[index].key)) {
			t.Errorf("entry %d key length = %d, want %d", index, got[index].keyLength, len(want[index].key))
		}
		if got[index].valueLength != uint16(len(want[index].value)) {
			t.Errorf("entry %d value length = %d, want %d", index, got[index].valueLength, len(want[index].value))
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

	data := encodeTestWALEntry(entry)
	if _, err := file.Write(data); err != nil {
		t.Fatalf("write test WAL entry: %v", err)
	}
}

func encodeTestWALEntry(entry walEntry) []byte {
	header := make([]byte, walEntryHeaderSize)
	binary.LittleEndian.PutUint32(header[4:8], entry.logNumber)
	header[8] = byte(entry.kind)
	binary.LittleEndian.PutUint16(header[9:11], uint16(len(entry.key)))
	binary.LittleEndian.PutUint16(header[11:13], uint16(len(entry.value)))

	checksum := crc32.NewIEEE()
	_, _ = checksum.Write(header[4:])
	_, _ = checksum.Write(entry.key)
	_, _ = checksum.Write(entry.value)
	binary.LittleEndian.PutUint32(header[0:4], checksum.Sum32())

	data := make([]byte, 0, len(header)+len(entry.key)+len(entry.value))
	data = append(data, header...)
	data = append(data, entry.key...)
	data = append(data, entry.value...)
	return data
}

func TestGetEntriesFromWALRejectsInvalidCRC(t *testing.T) {
	opts, _ := testOptions(t)

	log, err := LoadWAL(opts)
	if err != nil {
		t.Fatalf("LoadWAL() error = %v", err)
	}
	defer log.file.Close()

	data := encodeTestWALEntry(walEntry{
		logNumber: 1,
		kind:      walEntryPut,
		key:       []byte("key"),
		value:     []byte("value"),
	})
	data[len(data)-1] ^= 0xff

	if _, err := log.file.Write(data); err != nil {
		t.Fatalf("write corrupt WAL entry: %v", err)
	}

	_, err = log.GetEntriesFromWAL()
	if err == nil || !strings.Contains(err.Error(), "invalid CRC") {
		t.Fatalf("GetEntriesFromWAL() error = %v, want invalid CRC", err)
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
