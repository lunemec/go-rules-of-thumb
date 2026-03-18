## Append

> [!TIP]
> use `append(dst, src...)` as the default  
> if `len(src)` is comparable to or larger than `len(dst)` and this is hot code, preallocate the full result  
> avoid `for` + `append` without preallocation

![append graph](assets/BenchmarkAppend.png)

In this benchmark, `append(dst, src...)` wins most of the grid, especially when the appended slice is small. Once the appended slice gets large relative to the destination, the preallocated indexed copy often pulls ahead, so "append is always fastest" is too strong a rule.

[Detailed line view](assets/BenchmarkAppend-detail.png)

[Benchmark results](assets/BenchmarkAppend.txt)
