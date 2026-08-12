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
