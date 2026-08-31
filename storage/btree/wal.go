package engine

import "os"

type walRecordType uint8

type walRecord struct {
	txID   uint64
	kind   walRecordType
	pageID pageID
	data   []byte
}

type wal struct {
	file *os.File
}
