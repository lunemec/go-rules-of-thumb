## Strings concatenation

Is it more efficient to `"str1" + var`, `fmt.Sprintf()`, `strings.Join()` or `strings.Builder`? When does it make sense to add `sync.Pool`?

> [!TIP]
> for repeated concatenation, start with `strings.Builder`  
> treat `sync.Pool + strings.Builder` as a niche large-case optimization and benchmark it on your real workload  
> use `+` for one-off expressions and `fmt.Sprintf` for formatting, not concat speed
>
> in these benchmarks, plain `strings.Builder` is the safest default and `sync.Pool` does not reduce allocation totals

![concatenation graph](assets/BenchmarkConcat.png)

[Detailed line view](assets/BenchmarkConcat-detail.png)

On the main reduced matrix, `strings.Builder` wins the geomean and stays ahead in most cases. `sync.Pool + strings.Builder` does take a few isolated cells, but not enough to support a general threshold rule, so it still belongs in the "benchmark this exact workload" bucket.

[Benchmark results](assets/BenchmarkConcat.txt)

`strings.Builder.Reset()` currently discards the backing buffer, so pooling the builder does not preserve capacity here and does not lower `B/op`. This benchmark is measuring `sync.Pool` overhead around the builder, not reusable builder storage.

For the removed large-size region, a separate large-case-only matrix keeps the run practical while still checking whether pooling can help once the total bytes per operation get very large.

![large-case concatenation graph](assets/BenchmarkConcatLarge.png)

[Detailed large-case line view](assets/BenchmarkConcatLarge-detail.png)

In that large-case matrix, `sync.Pool + strings.Builder` does win several `500` to `5000` byte cases at `500+` concatenations per operation, but it still increases or matches allocations, so it should stay an opt-in benchmark target rather than a default recommendation.

[Large-case benchmark results](assets/BenchmarkConcatLarge.txt)
