package engine

import (
	"os"
	"sync"
)

type pager struct {
	file       *os.File
	pageSize   uint32
	pageCount  uint64
	syncWrites bool
	mu         sync.Mutex
}
