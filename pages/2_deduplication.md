## Deduplication

When is it more efficient to deduplicate a `slice` as opposed to using a `map[]struct{}` for the same purpose?

> [!TIP]
> if order does not matter, use in-place sort + dedup through at least ~5000 items  
> around `10000` items, benchmark `map` against sort + dedup on your workload  
> if you must preserve original order, use `map`  

![deduplication graph](assets/BenchmarkDeduplication.png)

In this benchmark, in-place sort + dedup is the fastest option from `10` through `5000` items, and `10000` items is effectively a wash with a slight edge to `map`. The plain slice scan is never the fastest path here; it is a simplicity choice for very small inputs, not a performance choice.

[Benchmark results](assets/BenchmarkDeduplication.txt)
