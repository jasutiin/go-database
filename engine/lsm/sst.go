package lsm

import "os"

type sstEntry struct {
	key    string
	offset int64
}

type sst struct {
	sparseIndex map[int]int // key-byte offset mapping for blocks
	bloomFilter filter      // bloom filter for key existence
	minKey      int         // minimum key in the SST
	maxKey      int         // maximum key in the SST
	file        *os.File    // underlying file for the SST
}

func newSST() {

}

func appendEntry() {

}
