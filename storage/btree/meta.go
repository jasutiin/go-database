package engine

type meta struct {
	pageSize       uint32
	rootPageID     pageID
	freelistPageID pageID
	pageCount      uint64
}
