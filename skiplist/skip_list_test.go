package skiplist

import (
	"errors"
	"slices"
	"testing"
)

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

func levelZeroValues(list *SkipList[int]) []int {
	var values []int

	for current := list.head.next[0]; current != nil; current = current.next[0] {
		values = append(values, current.data)
	}

	return values
}

func TestInsertIntoEmptyList(t *testing.T) {
	list := New[int](8, 100, compareInts)

	err := list.Insert(10)
	if err != nil {
		t.Fatalf("Insert() error = %v", err)
	}

	got := levelZeroValues(list)
	want := []int{10}

	if !slices.Equal(got, want) {
		t.Fatalf("values = %v, want %v", got, want)
	}
}

func TestInsertTracksUniqueValues(t *testing.T) {
	list := New[int](8, 100, compareInts)

	if got := list.Size(); got != 0 {
		t.Fatalf("initial size = %d, want 0", got)
	}

	for _, value := range []int{1, 2, 3} {
		if err := list.Insert(value); err != nil {
			t.Fatalf("Insert(%d) error = %v", value, err)
		}
	}

	if got := list.Size(); got != 3 {
		t.Fatalf("size after unique inserts = %d, want 3", got)
	}

	if err := list.Insert(2); err != nil {
		t.Fatalf("Insert(2) duplicate error = %v", err)
	}

	if got := list.Size(); got != 3 {
		t.Fatalf("size after duplicate insert = %d, want 3", got)
	}
}

func TestInsertRejectsNewValueAtMaximumSize(t *testing.T) {
	list := New[int](8, 2, compareInts)

	if got := list.MaxSize(); got != 2 {
		t.Fatalf("maximum size = %d, want 2", got)
	}

	if err := list.Insert(1); err != nil {
		t.Fatalf("Insert(1) error = %v", err)
	}
	if err := list.Insert(2); err != nil {
		t.Fatalf("Insert(2) error = %v", err)
	}

	if err := list.Insert(3); !errors.Is(err, ErrMaxSizeReached) {
		t.Fatalf("Insert(3) error = %v, want %v", err, ErrMaxSizeReached)
	}

	if got := list.Size(); got != 2 {
		t.Fatalf("size after rejected insert = %d, want 2", got)
	}
}

func TestInsertAllowsUpdateAtMaximumSize(t *testing.T) {
	list := New[testEntry](8, 1, compareEntries)

	if err := list.Insert(testEntry{key: 1, value: "old"}); err != nil {
		t.Fatal(err)
	}
	if err := list.Insert(testEntry{key: 1, value: "new"}); err != nil {
		t.Fatalf("duplicate update error = %v", err)
	}

	if got := list.Size(); got != 1 {
		t.Fatalf("size after duplicate update = %d, want 1", got)
	}

	if got := list.head.next[0].data.value; got != "new" {
		t.Fatalf("updated value = %q, want %q", got, "new")
	}
}

func TestInsertMaintainsSortedOrder(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{
			name:  "sorted input",
			input: []int{1, 2, 3},
			want:  []int{1, 2, 3},
		},
		{
			name:  "reverse input",
			input: []int{3, 2, 1},
			want:  []int{1, 2, 3},
		},
		{
			name:  "unordered input",
			input: []int{4, 1, 3, 2},
			want:  []int{1, 2, 3, 4},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			list := New[int](8, 100, compareInts)

			for _, value := range test.input {
				if err := list.Insert(value); err != nil {
					t.Fatalf("Insert(%d) error = %v", value, err)
				}
			}

			got := levelZeroValues(list)

			if !slices.Equal(got, test.want) {
				t.Fatalf("values = %v, want %v", got, test.want)
			}
		})
	}
}

type testEntry struct {
	key   int
	value string
}

func compareEntries(a, b testEntry) int {
	switch {
	case a.key < b.key:
		return -1
	case a.key > b.key:
		return 1
	default:
		return 0
	}
}

func TestInsertReplacesDuplicateKey(t *testing.T) {
	list := New[testEntry](8, 100, compareEntries)

	if err := list.Insert(testEntry{key: 1, value: "old"}); err != nil {
		t.Fatal(err)
	}

	if err := list.Insert(testEntry{key: 1, value: "new"}); err != nil {
		t.Fatal(err)
	}

	first := list.head.next[0]
	if first == nil {
		t.Fatal("list is empty")
	}

	if first.data.value != "new" {
		t.Fatalf("value = %q, want %q", first.data.value, "new")
	}

	if first.next[0] != nil {
		t.Fatal("duplicate key created a second node")
	}

	if got := list.Size(); got != 1 {
		t.Fatalf("size = %d, want 1", got)
	}
}
