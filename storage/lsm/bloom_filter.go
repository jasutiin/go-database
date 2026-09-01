package lsm

import "hash/maphash"

type filter struct {
	bits  []uint64
	m     uint64
	k     uint64
	seed1 maphash.Seed
	seed2 maphash.Seed
}
