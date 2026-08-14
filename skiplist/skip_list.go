package skiplist

import "math/rand/v2"

type Node[T any] struct {
	data T

	// left-most is the "widest" level
	// as you go to the right it narrows
	next []*Node[T]
}

type SkipList[T any] struct {
	head     *Node[T]
	level    int
	maxLevel int
	compare  func(a, b T) int
}

func NewSkipList[T any](maxLevel int, compare func(a, b T) int) *SkipList[T] {
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
func (s *SkipList[T]) Insert(node *Node[T]) error {
	randomLevel := rand.IntN(s.maxLevel)

	for level := 0; level < randomLevel; level++ {
		currentLevel := s.head.next[level]

		if currentLevel.next == nil {
			currentLevel.next[level] = node
		} else {
			current := s.head.next[level]
			// while the node is more than the next node
			for (s.compare(node.data, current.next[level].data)) > 0 {
				current = current.next[level]
			}

			// if key already exists, replace its data
			if s.compare(node.data, current.next[level].data) == 0 {
				current.next[level].data = node.data
			} else {
				node.next[level] = current.next[level]
				current.next[level] = node
			}
		}
	}

	return nil
}

func (s *SkipList[T]) Delete(node *Node[T]) error {
	return nil
}

func (s *SkipList[T]) Find(node *Node[T]) (*Node[T], error) {

}
