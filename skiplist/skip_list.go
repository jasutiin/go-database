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
	node := &node[T]{
		data: value,
		next: make([]*node[T], nodeMaxLevel),
	}

	current := s.head
	// start at the level with the largest gaps
	for level := nodeMaxLevel - 1; level >= 0; level-- {
		for current.next[level] != nil && s.compare(node.data, current.next[level].data) > 0 {
			current = current.next[level]
		}

		// if key already exists, replace its data
		if current.next[level] != nil && s.compare(node.data, current.next[level].data) == 0 {
			current.next[level].data = node.data
			return nil
		}

		if level == 0 && s.maxSize > 0 && s.size >= s.maxSize {
			return ErrMaxSizeReached
		}

		node.next[level] = current.next[level]
		current.next[level] = node
	}

	if nodeMaxLevel-1 > s.level {
		s.level = nodeMaxLevel - 1
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
