package lsm

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const sstEntryHeaderSize = 9

type sstEntry struct {
	key       string
	value     []byte
	tombstone bool
	offset    int64
}

type sst struct {
	path string
	file *os.File
	size int64
}

func newSST(opts *Options) (*sst, error) {
	executable, err := os.Executable()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(filepath.Dir(executable), opts.DbName)
	if err := os.MkdirAll(dbPath, 0o755); err != nil {
		return nil, err
	}

	file, err := os.CreateTemp(dbPath, "sst-*.db")
	if err != nil {
		return nil, fmt.Errorf("create SST: %w", err)
	}

	return &sst{
		path: file.Name(),
		file: file,
	}, nil
}

// SST entries are encoded as
// [key length:4][value length:4][tombstone:1][key][value].
func (table *sst) appendEntryToSSTFile(entry *sstEntry) error {
	data, err := encodeSSTEntry(entry)
	if err != nil {
		return err
	}

	entry.offset = table.size

	written, err := table.file.Write(data)
	if err != nil {
		return fmt.Errorf("append SST entry: %w", err)
	}
	if written != len(data) {
		return io.ErrShortWrite
	}

	table.size += int64(written)
	return nil
}

func encodeSSTEntry(entry *sstEntry) ([]byte, error) {
	const maxUint32 = uint64(^uint32(0))

	if uint64(len(entry.key)) > maxUint32 {
		return nil, fmt.Errorf("SST key is too large: %d bytes", len(entry.key))
	}
	if uint64(len(entry.value)) > maxUint32 {
		return nil, fmt.Errorf("SST value is too large: %d bytes", len(entry.value))
	}

	header := make([]byte, sstEntryHeaderSize)
	binary.LittleEndian.PutUint32(header[0:4], uint32(len(entry.key)))
	binary.LittleEndian.PutUint32(header[4:8], uint32(len(entry.value)))
	if entry.tombstone {
		header[8] = 1
	}

	data := make([]byte, 0, len(header)+len(entry.key)+len(entry.value))
	data = append(data, header...)
	data = append(data, entry.key...)
	data = append(data, entry.value...)
	return data, nil
}
