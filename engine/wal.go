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

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &wal{
		path: walPath,
		file: file,
		size: info.Size(),
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

// [crc:4][log number:4][kind:1][key length:2][value length:2][key][value].
func parseWALEntry(file *os.File) (*walEntry, error) {
	header := make([]byte, walEntryHeaderSize)
	_, err := io.ReadFull(file, header)
	if errors.Is(err, io.EOF) {
		return nil, io.EOF
	}
	if err != nil {
		return nil, fmt.Errorf("read WAL entry header: %w", err)
	}

	crc := binary.LittleEndian.Uint32(header[0:4])
	logNumber := binary.LittleEndian.Uint32(header[4:8])
	kind := walEntryKind(header[8])
	keyLength := binary.LittleEndian.Uint16(header[9:11])
	valueLength := binary.LittleEndian.Uint16(header[11:13])

	if kind != walEntryPut && kind != walEntryDelete {
		return nil, fmt.Errorf("unknown WAL entry kind %d", kind)
	}

	key := make([]byte, int(keyLength))
	if _, err := io.ReadFull(file, key); err != nil {
		return nil, fmt.Errorf("read WAL entry key: %w", err)
	}

	value := make([]byte, int(valueLength))
	if _, err := io.ReadFull(file, value); err != nil {
		return nil, fmt.Errorf("read WAL entry value: %w", err)
	}

	checksum := crc32.NewIEEE()
	_, _ = checksum.Write(header[4:])
	_, _ = checksum.Write(key)
	_, _ = checksum.Write(value)
	if crc != checksum.Sum32() {
		return nil, fmt.Errorf("invalid CRC")
	}

	return &walEntry{
		crc:         crc,
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

func (log *wal) Insert(key, value []byte, tombstone bool) error {
	kind := walEntryPut
	if tombstone {
		kind = walEntryDelete
	}

	entry := &walEntry{
		kind:  kind,
		key:   key,
		value: value,
	}

	data, err := encodeWALEntry(entry)
	if err != nil {
		return err
	}

	written, err := log.file.Write(data)
	if err != nil {
		return fmt.Errorf("append WAL entry: %w", err)
	}
	if written != len(data) {
		return io.ErrShortWrite
	}

	if err := log.file.Sync(); err != nil {
		return fmt.Errorf("sync WAL: %w", err)
	}

	log.size += int64(written)
	return nil
}

func encodeWALEntry(entry *walEntry) ([]byte, error) {
	const maxUint16 = int(^uint16(0))

	if len(entry.key) > maxUint16 {
		return nil, fmt.Errorf("WAL key is too large: %d bytes", len(entry.key))
	}
	if len(entry.value) > maxUint16 {
		return nil, fmt.Errorf("WAL value is too large: %d bytes", len(entry.value))
	}

	entry.keyLength = uint16(len(entry.key))
	entry.valueLength = uint16(len(entry.value))

	header := make([]byte, walEntryHeaderSize)
	binary.LittleEndian.PutUint32(header[4:8], entry.logNumber)
	header[8] = byte(entry.kind)
	binary.LittleEndian.PutUint16(header[9:11], entry.keyLength)
	binary.LittleEndian.PutUint16(header[11:13], entry.valueLength)

	checksum := crc32.NewIEEE()
	_, _ = checksum.Write(header[4:])
	_, _ = checksum.Write(entry.key)
	_, _ = checksum.Write(entry.value)
	entry.crc = checksum.Sum32()
	binary.LittleEndian.PutUint32(header[0:4], entry.crc)

	data := make([]byte, 0, len(header)+len(entry.key)+len(entry.value))
	data = append(data, header...)
	data = append(data, entry.key...)
	data = append(data, entry.value...)
	return data, nil
}
