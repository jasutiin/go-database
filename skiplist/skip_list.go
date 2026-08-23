package skiplist

import "math/rand/v2"

type node[T any] struct {
	data T

	// left-most is the "widest" level
	// as you go to the right it narrows
	next []*node[T]
}

type SkipList[T any] struct {
	head     *node[T]
	level    int
	maxLevel int
	compare  func(a, b T) int
}

func New[T any](maxLevel int, compare func(a, b T) int) *SkipList[T] {
	return &SkipList[T]{
		head: &node[T]{
			next: make([]*node[T], maxLevel),
		},
		level:    0,
		maxLevel: maxLevel,
		compare:  compare,
	}
}

// TODO: this insert has logic that depends on the node having a data field, change it later
func (s *SkipList[T]) Insert(value T) error {
	randomLevel := rand.IntN(s.maxLevel) + 1
	node := &node[T]{
		next: make([]*node[T], randomLevel),
	}
	node.data = value

	for level := range randomLevel {
		current := s.head

		if current.next[level] == nil {
			current.next[level] = node
		} else {
			// while the node is more than the next node
			for current.next[level] != nil && s.compare(node.data, current.next[level].data) > 0 {
				current = current.next[level]
			}

			// if key already exists, replace its data
			if current.next[level] != nil && s.compare(node.data, current.next[level].data) == 0 {
				current.next[level].data = node.data
				return nil
			}

			node.next[level] = current.next[level]
			current.next[level] = node
		}
	}

	if randomLevel-1 > s.level {
		s.level = randomLevel - 1
	}

	return nil
}

func (s *SkipList[T]) Delete(value T) error {
	return nil
}

func (s *SkipList[T]) Find(value T) (*node[T], error) {
	return nil, nil
}
