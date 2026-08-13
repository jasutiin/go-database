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

func LoadMemTable() (*memTable, error) {
	return &memTable{
		entries: make(map[string][]memTableEntry),
		maxSize: 1024,
	}, nil
}
