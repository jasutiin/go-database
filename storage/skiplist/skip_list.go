package skiplist

import (
	"errors"
	"math/rand/v2"
)

type node[K, V any] struct {
	key   K
	value V

	// smallest index is the level with no gaps
	next []*node[K, V]
}

var ErrMaxSizeReached = errors.New("skip list maximum size reached")

type SkipList[K, V any] struct {
	head     *node[K, V]
	level    int
	maxLevel int
	size     int
	maxSize  int
	compare  func(a, b K) int
}

func New[K, V any](maxLevel, maxSize int, compare func(a, b K) int) *SkipList[K, V] {
	return &SkipList[K, V]{
		head: &node[K, V]{
			next: make([]*node[K, V], maxLevel),
		},
		level:    0,
		maxLevel: maxLevel,
		maxSize:  maxSize,
		compare:  compare,
	}
}

const promotionProbability = 0.25

func randomLevel(maxLevel int) int {
	level := 1
	for level < maxLevel && rand.Float64() < promotionProbability {
		level++
	}

	return level
}

func (s *SkipList[K, V]) Insert(key K, value V) error {
	return s.insertAtLevel(key, value, randomLevel(s.maxLevel))
}

func (s *SkipList[K, V]) insertAtLevel(key K, value V, nodeMaxLevel int) error {
	update := make([]*node[K, V], s.maxLevel)
	current := s.head

	for level := s.level; level >= 0; level-- {
		for current.next[level] != nil && s.compare(key, current.next[level].key) > 0 {
			current = current.next[level]
		}

		update[level] = current
	}

	existing := update[0].next[0]
	if existing != nil && s.compare(key, existing.key) == 0 {
		existing.value = value
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

	node := &node[K, V]{
		key:   key,
		value: value,
		next:  make([]*node[K, V], nodeMaxLevel),
	}

	for level := range nodeMaxLevel {
		node.next[level] = update[level].next[level]
		update[level].next[level] = node
	}

	s.size++
	return nil
}

func (s *SkipList[K, V]) Find(key K) (V, bool) {
	current := s.head

	for level := s.level; level >= 0; level-- {
		for current.next[level] != nil {
			comparison := s.compare(key, current.next[level].key)
			if comparison > 0 {
				current = current.next[level]
			} else if comparison == 0 {
				return current.next[level].value, true
			} else {
				break
			}
		}
	}

	var empty V
	return empty, false
}

func (s *SkipList[K, V]) RemoveFront() (K, V, bool) {
	front := s.head.next[0]
	if front == nil {
		var emptyKey K
		var emptyValue V
		return emptyKey, emptyValue, false
	}

	for level := range front.next {
		if s.head.next[level] == front {
			s.head.next[level] = front.next[level]
		}
	}

	s.size--
	s.reduceLevel()
	return front.key, front.value, true
}

func (s *SkipList[K, V]) RemoveBack() (K, V, bool) {
	if s.head.next[0] == nil {
		var emptyKey K
		var emptyValue V
		return emptyKey, emptyValue, false
	}

	back := s.head
	for level := s.level; level >= 0; level-- {
		for back.next[level] != nil {
			back = back.next[level]
		}
	}

	for level := 0; level <= s.level; level++ {
		previous := s.head
		for previous.next[level] != nil && previous.next[level] != back {
			previous = previous.next[level]
		}
		if previous.next[level] == back {
			previous.next[level] = back.next[level]
		}
	}

	s.size--
	s.reduceLevel()
	return back.key, back.value, true
}

func (s *SkipList[K, V]) reduceLevel() {
	for s.level > 0 && s.head.next[s.level] == nil {
		s.level--
	}
}

func (s *SkipList[K, V]) Size() int {
	return s.size
}

func (s *SkipList[K, V]) MaxSize() int {
	return s.maxSize
}
