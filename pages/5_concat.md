## Strings concatenation

Is it more efficient to `"str1" + var`, `fmt.Sprintf()`, `strings.Join()` or `strings.Builder`? When does it make sense to add `sync.Pool`?

> [!TIP]
> for repeated concatenation, start with `strings.Builder`  
> benchmark `sync.Pool + strings.Builder` once the loop gets very hot or reaches ~1000+ concatenations per operation  
> use `+` for one-off expressions and `fmt.Sprintf` for formatting, not concat speed
>
> in this benchmark, only the `strings.Builder` variants ever win

![concatenation graph](assets/BenchmarkConcat.png)

[Detailed line view](assets/BenchmarkConcat-detail.png)

In this benchmark, only the two `strings.Builder` variants take first place. `sync.Pool + strings.Builder` starts to win more often as the work gets heavier, but plain `strings.Builder` stays competitive across the whole grid, which makes it the safest default.

[Benchmark results](assets/BenchmarkConcat.txt)
