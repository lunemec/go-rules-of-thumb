## Deduplication

When is it more efficient to deduplicate a `slice` as opposed to using a `map[]struct{}` for the same purpose?

> [!TIP]
> use `map` when `len(haystack) > 100`.  
> if you must reduce allocations, use in-place sort + deduplication  
> if you must preserve original order, use `slice` or other methods

![deduplication graph](assets/BenchmarkDeduplication.png)

[Benchmark results](assets/BenchmarkDeduplication.txt)
