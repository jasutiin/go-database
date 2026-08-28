package engine

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
)

type wal struct {
	path string
	file *os.File
	size int64
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

type walEntryKind byte

const (
	walEntryPut walEntryKind = iota + 1
	walEntryDelete
)

type walEntry struct {
	crc         uint32
	logNumber   uint32
	kind        walEntryKind
	keyLength   uint16
	valueLength uint16
	key         []byte
	value       []byte
}

const walEntryHeaderSize = 13

// [crc:4][log number:4][kind:1][key length:4][value length:4][key][value].
func parseWALEntry(file *os.File) (*walEntry, error) {
	header := make([]byte, walEntryHeaderSize)
	_, err := io.ReadFull(file, header)
	if errors.Is(err, io.EOF) {
		return nil, io.EOF
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
	crc := binary.LittleEndian.Uint32(header[1:5])
	if crc != crc32.ChecksumIEEE(header[5:13]) {
		return nil, fmt.Errorf("invalid CRC")
	}

	logNumber := binary.LittleEndian.Uint32(header[5:9])
	keyLength := binary.LittleEndian.Uint16(header[9:13])
	valueLength := binary.LittleEndian.Uint16(header[13:17])

	key := make([]byte, int(keyLength))
	if _, err := io.ReadFull(file, key); err != nil {
		return nil, fmt.Errorf("read WAL entry key: %w", err)
	}

	value := make([]byte, int(valueLength))
	if _, err := io.ReadFull(file, value); err != nil {
		return nil, fmt.Errorf("read WAL entry value: %w", err)
	}

	return &walEntry{
		logNumber:   logNumber,
		kind:        kind,
		keyLength:   keyLength,
		valueLength: valueLength,
		key:         key,
		value:       value,
	}, nil
}

func (log *wal) GetEntriesFromWAL() ([]*walEntry, error) {
	file, err := os.Open(log.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var walEntries []*walEntry

	for {
		entry, err := parseWALEntry(file)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parse WAL entry: %w", err)
		}

		walEntries = append(walEntries, entry)
	}

	return walEntries, nil
}
