## Deduplication

When is it more efficient to deduplicate a `slice` as opposed to using a `map[]struct{}` for the same purpose?

> [!TIP]
> if order does not matter, use in-place sort + dedup up to ~1000 items  
> use `map` from roughly ~5000 items upward  
> if you must preserve original order, use `map`

![deduplication graph](assets/BenchmarkDeduplication.png)

In this benchmark, in-place sort + dedup is the fastest option from `10` through `1000` items, while `map` takes over from `5000` onward. The plain slice scan is never the fastest path here; it is a simplicity choice for very small inputs, not a performance choice.

[Benchmark results](assets/BenchmarkDeduplication.txt)
