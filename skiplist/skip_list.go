package skiplist

import "math/rand/v2"

type Node[T any] struct {
	data T

	// left-most is the "widest" level
	// as you go to the right it narrows
	next []*Node[T]
}

func NewNode[T any](data T) *Node[T] {
	return &Node[T]{
		data: data,
	}
}

type SkipList[T any] struct {
	head     *Node[T]
	level    int
	maxLevel int
	compare  func(a, b T) int
}

func New[T any](maxLevel int, compare func(a, b T) int) *SkipList[T] {
	return &SkipList[T]{
		head: &Node[T]{
			next: make([]*Node[T], maxLevel),
		},
		level:    0,
		maxLevel: maxLevel,
		compare:  compare,
	}
}

// TODO: this insert has logic that depends on the node having a data field, change it later
func (s *SkipList[T]) Insert(value T) error {
	randomLevel := rand.IntN(s.maxLevel) + 1
	node := &Node[T]{
		next: make([]*Node[T], randomLevel),
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

func (s *SkipList[T]) Delete(node *Node[T]) error {
	return nil
}

func (s *SkipList[T]) Find(node *Node[T]) (*Node[T], error) {
	return nil, nil
}
