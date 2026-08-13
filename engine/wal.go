package engine

import (
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

	exPath := filepath.Dir(ex)
	walPath := exPath + opts.DbName + "/wal.log"
	wal := new(wal)

	// create if missing, otherwise allow read/write. writes go to the end of the file
	// 0o644 means owner = r/w, group = r, others = r
	file, err := os.OpenFile(walPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}

	// TODO: set size for wal?
	wal.path = walPath
	wal.file = file

	if err != nil {
		return nil, err
	}

	return wal, nil
}
