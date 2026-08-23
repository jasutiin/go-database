the motivation for this project is to understand the concepts taught in the following books:
- designing data-intensive applications
- database internals

<br>

also looking at the following for inspo:
- https://github.com/google/leveldb
- https://www.youtube.com/@benjdicken

<br>

to run tests, use `go test -v ./...`.

<br>

to run the skip list benchmarks, use `go test ./skiplist -run '^$' -bench '.' -benchmem`.

```text
BenchmarkInsertBuild/size=1000-12    8    13836212 ns/op    1000 inserts/op    13836 ns/insert    103401 B/op    2001 allocs/op
```
how to read benchmark results:
- `BenchmarkInsertBuild/size=1000-12` describes the benchmark that ran
  - `BenchmarkInsertBuild` is the benchmark name
  - `size=1000` means one operation creates a skip list containing 1,000 inserted values
  - `-12` means `GOMAXPROCS` was 12
    - go was allowed to execute go code on up to 12 operating system threads at the same time
    - this describes the runtime's available parallel capacity, not how many benchmark operations were running
    - the benchmark uses one sequential goroutine, so it finishes each operation before beginning the next and generally only needs one thread
    - the other available threads matter when code starts multiple goroutines or when a benchmark explicitly uses `b.RunParallel`
- `8` is the number of operations go ran
  - each operation created one skip list and performed 1,000 inserts
  - 8 operations resulted in 8 separate skip lists and 8,000 total inserts
- timing measurements
  - `13836212 ns/op` means one complete 1,000-value build took about 13.8 milliseconds on average
  - `1000 inserts/op` confirms that each operation contained 1,000 inserts
  - `13836 ns/insert` means one insert took about 13.8 microseconds on average
- memory measurements
  - `103401 B/op` means one complete build allocated about 103 kilobytes of heap memory
  - `2001 allocs/op` means one complete build performed about 2,001 heap allocations, or roughly 2 allocations per insert
