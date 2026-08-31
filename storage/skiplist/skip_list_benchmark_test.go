package skiplist

import (
	"fmt"
	"math/rand"
	"testing"
)

var (
	benchmarkValue int
	benchmarkFound bool
	benchmarkSize  int
)

// benchmarkCompareInts orders integer values used by the skip list benchmarks.
func benchmarkCompareInts(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

// BenchmarkInsertBuild measures creating and filling a new skip list. One
// operation builds one list containing the number of values in the subtest name.
func BenchmarkInsertBuild(b *testing.B) {
	for _, size := range []int{100, 1_000, 5_000} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			values := rand.New(rand.NewSource(1)).Perm(size)
			b.ReportAllocs()
			b.ResetTimer()

			for iteration := 0; iteration < b.N; iteration++ {
				list := New[int](16, 0, benchmarkCompareInts)

				for _, value := range values {
					if err := list.Insert(value); err != nil {
						b.Fatal(err)
					}
				}

				benchmarkSize = list.Size()
			}

			b.ReportMetric(float64(size), "inserts/op")
			b.ReportMetric(
				float64(b.Elapsed().Nanoseconds())/float64(b.N*size),
				"ns/insert",
			)
		})
	}
}

// BenchmarkInsertExisting measures replacing a key in a prebuilt skip list.
// One operation is one Insert call for a key that already exists.
func BenchmarkInsertExisting(b *testing.B) {
	for _, size := range []int{100, 1_000, 5_000} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			list := newBenchmarkList(b, size)
			value := size / 2

			b.ReportAllocs()
			b.ResetTimer()

			for iteration := 0; iteration < b.N; iteration++ {
				if err := list.Insert(value); err != nil {
					b.Fatal(err)
				}
			}

			benchmarkSize = list.Size()
		})
	}
}

// BenchmarkFind measures lookups in a prebuilt skip list using a mixture of
// existing and missing values. One operation is one Find call.
func BenchmarkFind(b *testing.B) {
	for _, size := range []int{100, 1_000, 5_000} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			list := newBenchmarkList(b, size)
			queries := []int{0, size / 2, size - 1, size + 1}

			b.ReportAllocs()
			b.ResetTimer()

			for iteration := 0; iteration < b.N; iteration++ {
				benchmarkValue, benchmarkFound = list.Find(
					queries[iteration%len(queries)],
				)
			}
		})
	}
}

// newBenchmarkList creates a consistently shuffled list outside the timed
// portion of benchmarks that need prebuilt data.
func newBenchmarkList(b *testing.B, size int) *SkipList[int] {
	b.Helper()

	list := New[int](16, 0, benchmarkCompareInts)
	values := rand.New(rand.NewSource(1)).Perm(size)

	for _, value := range values {
		if err := list.Insert(value); err != nil {
			b.Fatal(err)
		}
	}

	return list
}
