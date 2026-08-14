package skiplist

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
}

func (s *SkipList[T]) Insert(node *Node[T]) error {
	return nil
}

func (s *SkipList[T]) Delete(node *Node[T]) error {
	return nil
}

func (s *SkipList[T]) Find(node *Node[T]) (*Node[T], error) {

}
