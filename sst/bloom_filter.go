package sst

import "hash/maphash"

type Filter struct {
	bits  []uint64
	m     uint64
	k     uint64
	seed1 maphash.Seed
	seed2 maphash.Seed
}
