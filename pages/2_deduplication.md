## Deduplication

When is it more efficient to deduplicate a `slice` as opposed to using a `map[]struct{}` for the same purpose?

> **TL;DR**: use `map` when `len(haystack) > 100`. If you must reduce allocations, use sort + in-place slice
> deduplication. Suprisingly it is fast enough.

![deduplication graph](assets/BenchmarkDeduplication.png)

[Benchmark results](assets/BenchmarkDeduplication.txt)
