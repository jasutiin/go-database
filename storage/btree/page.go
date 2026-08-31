package engine

type pageID uint64

type pageType uint8

type page struct {
	id   pageID
	kind pageType
	data []byte
}
