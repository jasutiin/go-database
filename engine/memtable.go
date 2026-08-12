package engine

type memTableEntry struct {
	key       string
	value     []byte
	tombstone bool
}

type memTable struct {
	entries map[string][]memTableEntry
	size    int
	maxSize int
}

func NewMemTable() *memTable {
	return &memTable{
		entries: make(map[string][]memTableEntry),
		maxSize: 1024,
	}
}
