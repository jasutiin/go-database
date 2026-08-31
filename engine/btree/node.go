package btree

type nodeType uint8

type node struct {
	pageID     pageID
	kind       nodeType
	keys       [][]byte
	values     [][]byte
	children   []pageID
	nextLeafID pageID
}

type pathElement struct {
	pageID     pageID
	childIndex int
}
