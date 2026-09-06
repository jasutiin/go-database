package lsm

import (
	"encoding/binary"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestEngineFlushWritesMemTableEntriesToSST(t *testing.T) {
	opts, _ := testOptions(t)
	opts.SkipListMaxSize = 2

	engine, err := Startup(opts)
	if err != nil {
		t.Fatalf("Startup() error = %v", err)
	}
	defer engine.writeAheadLog.file.Close()

	if err := engine.Put([]byte("second"), []byte("two")); err != nil {
		t.Fatalf("Put(second) error = %v", err)
	}
	if err := engine.Put([]byte("first"), []byte("one")); err != nil {
		t.Fatalf("Put(first) error = %v", err)
	}

	var table *sst
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		engine.mutex.RLock()
		if len(engine.sstables) == 1 {
			table = engine.sstables[0]
		}
		engine.mutex.RUnlock()

		if table != nil {
			break
		}
		time.Sleep(time.Millisecond)
	}

	if table == nil {
		t.Fatal("flush did not publish an SST")
	}
	defer table.file.Close()

	data, err := os.ReadFile(table.path)
	if err != nil {
		t.Fatalf("read SST: %v", err)
	}

	entries, err := decodeSSTEntriesForTest(data)
	if err != nil {
		t.Fatalf("decode SST: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("SST entry count = %d, want 2", len(entries))
	}

	if entries[0].key != "first" || string(entries[0].value) != "one" {
		t.Fatalf("first SST entry = %q/%q, want first/one", entries[0].key, entries[0].value)
	}
	if entries[1].key != "second" || string(entries[1].value) != "two" {
		t.Fatalf("second SST entry = %q/%q, want second/two", entries[1].key, entries[1].value)
	}
}

func decodeSSTEntriesForTest(data []byte) ([]sstEntry, error) {
	var entries []sstEntry

	for offset := 0; offset < len(data); {
		if len(data)-offset < sstEntryHeaderSize {
			return nil, fmt.Errorf("incomplete SST entry header")
		}

		entryOffset := offset
		keyLength := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
		valueLength := int(binary.LittleEndian.Uint32(data[offset+4 : offset+8]))
		tombstone := data[offset+8] == 1
		offset += sstEntryHeaderSize

		entryLength := keyLength + valueLength
		if len(data)-offset < entryLength {
			return nil, fmt.Errorf("incomplete SST entry")
		}

		entries = append(entries, sstEntry{
			key:       string(data[offset : offset+keyLength]),
			value:     append([]byte(nil), data[offset+keyLength:offset+entryLength]...),
			tombstone: tombstone,
			offset:    int64(entryOffset),
		})
		offset += entryLength
	}

	return entries, nil
}
