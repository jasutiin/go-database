package engine

import "os"

type wal struct {
	path string
	file *os.File
	size int64
}
