package sst

import "os"

type SST struct {
	sparseIndex map[int]int // key-byte offset mapping for blocks
	bloomFilter Filter      // bloom filter for key existence
	minKey      int         // minimum key in the SST
	maxKey      int         // maximum key in the SST
	file        *os.File    // underlying file for the SST
}

func NewSST() {

}

func AppendEntry() {

}
