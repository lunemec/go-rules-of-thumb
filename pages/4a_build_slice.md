## Building a slice from scratch

When you know the final length up front, should you `append` into a growing slice, `append` into a preallocated slice, or pre-size the slice and fill by index?

> [!TIP]
> if the final length is known, do not grow from `nil`
> use `make([]T, 0, n)` + `append` as the readable default
> on hot paths, benchmark `make([]T, n)` + index writes; it wins `23/27` cells here and beats preallocated `append` by up to about `54%`
> plain `append` growth is never best here and can be about `4x` slower with many more allocations

![build slice graph](assets/BenchmarkBuildSlice.png)

[Detailed line view](assets/BenchmarkBuildSlice-detail.png)

Across `8B`, `32B`, and `128B` elements from `10` to `100 000` items, plain growth never wins. `make([]T, 0, n)` + `append` consistently removes most of that cost, and pre-sized index writes usually take the remaining lead, especially once the output is non-trivial. At `32B x 100 000`, for example, growing with `append` takes about `967µs`, `16.4 MiB/op`, and `29 allocs/op`, while pre-sized indexing drops that to about `242µs`, `3.2 MiB/op`, and `2 allocs/op`.

[Benchmark results](assets/BenchmarkBuildSlice.txt)
