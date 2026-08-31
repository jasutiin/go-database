package skiplist

import (
	"errors"
	"math/rand/v2"
)

type node[T any] struct {
	data T

	// smallest index is the level with no gaps
	next []*node[T]
}

var ErrMaxSizeReached = errors.New("skip list maximum size reached")

type SkipList[T any] struct {
	head     *node[T]
	level    int
	maxLevel int
	size     int
	maxSize  int
	compare  func(a, b T) int
}

func New[T any](maxLevel, maxSize int, compare func(a, b T) int) *SkipList[T] {
	return &SkipList[T]{
		head: &node[T]{
			next: make([]*node[T], maxLevel),
		},
		level:    0,
		maxLevel: maxLevel,
		maxSize:  maxSize,
		compare:  compare,
	}
}

func randomLevel(maxLevel int) int {
	return rand.IntN(maxLevel) + 1
}

func (s *SkipList[T]) Insert(value T) error {
	return s.insertAtLevel(value, randomLevel(s.maxLevel))
}

func (s *SkipList[T]) insertAtLevel(value T, nodeMaxLevel int) error {
	update := make([]*node[T], s.maxLevel)
	current := s.head

	for level := s.level; level >= 0; level-- {
		for current.next[level] != nil && s.compare(value, current.next[level].data) > 0 {
			current = current.next[level]
		}

		update[level] = current
	}

	existing := update[0].next[0]
	if existing != nil && s.compare(value, existing.data) == 0 {
		existing.data = value
		return nil
	}

	if s.maxSize > 0 && s.size >= s.maxSize {
		return ErrMaxSizeReached
	}

	if nodeMaxLevel-1 > s.level {
		for level := s.level + 1; level < nodeMaxLevel; level++ {
			update[level] = s.head
		}

		s.level = nodeMaxLevel - 1
	}

	node := &node[T]{
		data: value,
		next: make([]*node[T], nodeMaxLevel),
	}

	for level := range nodeMaxLevel {
		node.next[level] = update[level].next[level]
		update[level].next[level] = node
	}

	s.size++
	return nil
}

func (s *SkipList[T]) Find(value T) (T, bool) {
	current := s.head

	for level := s.level; level >= 0; level-- {
		for current.next[level] != nil {
			if s.compare(value, current.next[level].data) > 0 {
				current = current.next[level]
			} else if s.compare(value, current.next[level].data) == 0 {
				return current.next[level].data, true
			} else {
				break // go to next level
			}
		}
	}

	var empty T
	return empty, false
}

func (s *SkipList[T]) Size() int {
	return s.size
}

func (s *SkipList[T]) MaxSize() int {
	return s.maxSize
}
