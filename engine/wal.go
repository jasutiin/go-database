package engine

import (
	"errors"
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

// TODO: check for if wal has unfinished entries, deal with those too
func (log *wal) GetEntriesFromWAL() ([]*walEntry, error) {
	file, err := os.Open(log.path)

	if err != nil {
		return nil, err
	}

	walEntries := make([]*walEntry, 0)

	for {
		kindBytes := make([]byte, 1)
		_, err := file.Read(kindBytes)

		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}

		keyBytes := make([]byte, 8)
		_, err = file.Read(keyBytes)

		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}

		valueBytes := make([]byte, 8)
		_, err = file.Read(valueBytes)

		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}

		kind := walEntryKind(kindBytes[0])
		switch kind {
		case walEntryPut, walEntryDelete:
			entries = append(entries, &walEntry{
				kind:  kind,
				key:   keyBytes,
				value: valueBytes,
			})
	}

	return walEntries, nil
}
