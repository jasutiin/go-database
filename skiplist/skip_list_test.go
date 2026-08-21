package skiplist

import (
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
	list := NewSkipList[int](8, compareInts)

	err := list.Insert(&Node[int]{data: 10})
	if err != nil {
		t.Fatalf("Insert() error = %v", err)
	}

	got := levelZeroValues(list)
	want := []int{10}

	if !slices.Equal(got, want) {
		t.Fatalf("values = %v, want %v", got, want)
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
			list := NewSkipList[int](8, compareInts)

			for _, value := range test.input {
				if err := list.Insert(&Node[int]{data: value}); err != nil {
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
	list := NewSkipList[testEntry](8, compareEntries)

	if err := list.Insert(&Node[testEntry]{
		data: testEntry{key: 1, value: "old"},
	}); err != nil {
		t.Fatal(err)
	}

	if err := list.Insert(&Node[testEntry]{
		data: testEntry{key: 1, value: "new"},
	}); err != nil {
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
}
