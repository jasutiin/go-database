package engine

import "os"

type wal struct {
	path string
	file *os.File
	size int64
}

func LoadWAL() (*wal, error) {
	return &wal{
		path: "wal.log",
	}, nil
}
