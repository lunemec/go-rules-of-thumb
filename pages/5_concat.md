## Strings concatenation

Is it more efficient to `"str1" + var`, `fmt.Sprintf()`, `strings.Join()` or `strings.Builder`? When does it make sense to add `sync.Pool`?

> **TL;DR**: use `strings.Builder` when `len(str) < 100 & N ops < 1000`, use `sync.Pool + strings.Builder` when doing this for every request. For `len(str) > 100` use `+` or `strings.Join`.
>
> Use `fmt.Sprintf` for regular string formatting (not just concatenation).

![concatenation graph](assets/BenchmarkConcat.png)

[Benchmark results](assets/BenchmarkConcat.txt)
