package skiplist

import (
	"errors"
	"slices"
	"testing"
)

type testSkipList[K, V any] struct {
	*SkipList[K, V]
	levels []int
	index  int
}

func newTestSkipList[K, V any](
	maxLevel int,
	maxSize int,
	compare func(a, b K) int,
) *testSkipList[K, V] {
	return &testSkipList[K, V]{
		SkipList: New[K, V](maxLevel, maxSize, compare),
		levels:   []int{1, 3, 2, 1},
	}
}

func (s *testSkipList[K, V]) Insert(key K, value V) error {
	level := s.levels[s.index%len(s.levels)]
	s.index++

	if level > s.maxLevel {
		level = s.maxLevel
	}

	return s.insertAtLevel(key, value, level)
}

func compareInts(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

func levelZeroValues(list *testSkipList[int, int]) []int {
	return valuesAtLevel(list, 0)
}

func valuesAtLevel(list *testSkipList[int, int], level int) []int {
	var values []int

	for current := list.head.next[level]; current != nil; current = current.next[level] {
		values = append(values, current.value)
	}

	return values
}

func TestInsertIntoEmptyList(t *testing.T) {
	list := newTestSkipList[int, int](8, 100, compareInts)

	if err := list.Insert(10, 10); err != nil {
		t.Fatalf("Insert() error = %v", err)
	}

	if got, want := levelZeroValues(list), []int{10}; !slices.Equal(got, want) {
		t.Fatalf("values = %v, want %v", got, want)
	}
}

func TestInsertUsesGeneratedLevels(t *testing.T) {
	list := newTestSkipList[int, int](8, 100, compareInts)

	for _, value := range []int{10, 20, 15} {
		if err := list.Insert(value, value); err != nil {
			t.Fatalf("Insert(%d) error = %v", value, err)
		}
	}

	wantByLevel := [][]int{
		{10, 15, 20},
		{15, 20},
		{20},
	}

	for level, want := range wantByLevel {
		if got := valuesAtLevel(list, level); !slices.Equal(got, want) {
			t.Fatalf("level %d values = %v, want %v", level, got, want)
		}
	}

	if list.level != 2 {
		t.Fatalf("active level = %d, want 2", list.level)
	}
}

func TestInsertTracksUniqueKeys(t *testing.T) {
	list := newTestSkipList[int, int](8, 100, compareInts)

	if got := list.Size(); got != 0 {
		t.Fatalf("initial size = %d, want 0", got)
	}

	for _, value := range []int{1, 2, 3} {
		if err := list.Insert(value, value); err != nil {
			t.Fatalf("Insert(%d) error = %v", value, err)
		}
	}

	if got := list.Size(); got != 3 {
		t.Fatalf("size after unique inserts = %d, want 3", got)
	}

	if err := list.Insert(2, 20); err != nil {
		t.Fatalf("Insert(2) duplicate error = %v", err)
	}

	if got := list.Size(); got != 3 {
		t.Fatalf("size after duplicate insert = %d, want 3", got)
	}
}

func TestInsertRejectsNewKeyAtMaximumSize(t *testing.T) {
	list := newTestSkipList[int, int](8, 2, compareInts)

	if got := list.MaxSize(); got != 2 {
		t.Fatalf("maximum size = %d, want 2", got)
	}

	if err := list.Insert(1, 1); err != nil {
		t.Fatalf("Insert(1) error = %v", err)
	}
	if err := list.Insert(2, 2); err != nil {
		t.Fatalf("Insert(2) error = %v", err)
	}
	if err := list.Insert(3, 3); !errors.Is(err, ErrMaxSizeReached) {
		t.Fatalf("Insert(3) error = %v, want %v", err, ErrMaxSizeReached)
	}

	if got := list.Size(); got != 2 {
		t.Fatalf("size after rejected insert = %d, want 2", got)
	}

	if got := valuesAtLevel(list, 0); !slices.Equal(got, []int{1, 2}) {
		t.Fatalf("level 0 values = %v, want [1 2]", got)
	}
	if got := valuesAtLevel(list, 1); !slices.Equal(got, []int{2}) {
		t.Fatalf("level 1 values = %v, want [2]", got)
	}
}

type testEntry struct {
	name string
}

func TestInsertAllowsUpdateAtMaximumSize(t *testing.T) {
	list := newTestSkipList[int, testEntry](8, 1, compareInts)

	if err := list.Insert(1, testEntry{name: "old"}); err != nil {
		t.Fatal(err)
	}
	if err := list.Insert(1, testEntry{name: "new"}); err != nil {
		t.Fatalf("duplicate update error = %v", err)
	}

	if got := list.Size(); got != 1 {
		t.Fatalf("size after duplicate update = %d, want 1", got)
	}

	if got := list.head.next[0].value.name; got != "new" {
		t.Fatalf("updated value = %q, want %q", got, "new")
	}
}

func TestInsertMaintainsSortedOrder(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{name: "sorted input", input: []int{1, 2, 3}, want: []int{1, 2, 3}},
		{name: "reverse input", input: []int{3, 2, 1}, want: []int{1, 2, 3}},
		{name: "unordered input", input: []int{4, 1, 3, 2}, want: []int{1, 2, 3, 4}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			list := newTestSkipList[int, int](8, 100, compareInts)

			for _, value := range test.input {
				if err := list.Insert(value, value); err != nil {
					t.Fatalf("Insert(%d) error = %v", value, err)
				}
			}

			if got := levelZeroValues(list); !slices.Equal(got, test.want) {
				t.Fatalf("values = %v, want %v", got, test.want)
			}
		})
	}
}

func TestFindUsesKey(t *testing.T) {
	list := newTestSkipList[int, testEntry](8, 100, compareInts)
	if err := list.Insert(1, testEntry{name: "Ada"}); err != nil {
		t.Fatal(err)
	}
	if err := list.Insert(2, testEntry{name: "Grace"}); err != nil {
		t.Fatal(err)
	}

	entry, found := list.Find(2)
	if !found || entry.name != "Grace" {
		t.Fatalf("Find(2) = %+v, found=%t", entry, found)
	}

	_, found = list.Find(3)
	if found {
		t.Fatal("Find(3) unexpectedly found a value")
	}
}
