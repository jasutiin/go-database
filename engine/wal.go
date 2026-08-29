package engine

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type wal struct {
	path string
	file *os.File
	size int64
}

type walEntryKind byte

const (
	walEntryPut walEntryKind = iota + 1
	walEntryDelete
)

type walEntry struct {
	kind  walEntryKind
	key   []byte
	value []byte
}

func LoadWAL(opts *Options) (*wal, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(filepath.Dir(ex), opts.DbName)

	if err := os.MkdirAll(dbPath, 0o755); err != nil {
		return nil, err
	}

	walPath := filepath.Join(dbPath, "wal.log")

	// create if missing, otherwise allow read/write. writes go to the end of the file
	// 0o644 means owner = r/w, group = r, others = r
	file, err := os.OpenFile(walPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}

	// TODO: set size for wal?
	return &wal{
		path: walPath,
		file: file,
	}, nil
}

const walEntryHeaderSize = 9

// GetEntriesFromWAL reads records encoded as:
// [kind:1][key length:4][value length:4][key][value].
func (log *wal) GetEntriesFromWAL() ([]*walEntry, error) {
	file, err := os.Open(log.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var walEntries []*walEntry

	for {
		header := make([]byte, walEntryHeaderSize)
		_, err := io.ReadFull(file, header)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read WAL entry header: %w", err)
		}

		kind := walEntryKind(header[0])
		if kind != walEntryPut && kind != walEntryDelete {
			return nil, fmt.Errorf("unknown WAL entry kind %d", kind)
		}

		// not using int because it depends on the platform. uint32 ensures
		// we are reading the same amount of bytes from the file in any platform
		keyLength := binary.LittleEndian.Uint32(header[1:5])
		valueLength := binary.LittleEndian.Uint32(header[5:9])

		key := make([]byte, int(keyLength))
		if _, err := io.ReadFull(file, key); err != nil {
			return nil, fmt.Errorf("read WAL entry key: %w", err)
		}

		value := make([]byte, int(valueLength))
		if _, err := io.ReadFull(file, value); err != nil {
			return nil, fmt.Errorf("read WAL entry value: %w", err)
		}

		walEntries = append(walEntries, &walEntry{
			kind:  kind,
			key:   key,
			value: value,
		})
	}

	return walEntries, nil
}

func (log *wal) Insert(key, val []byte, isTombstone bool) error {
	file, err := os.Open(log.path)
	return nil
}
