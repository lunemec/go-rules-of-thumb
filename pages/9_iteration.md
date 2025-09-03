## Range over func

With [Go 1.23 came new feature - range over func](https://go.dev/blog/range-functions), lets check when it makes sense to use that over
pre-allocating a slice and putting values in it.

> [!TIP]
> use `iter.Seq` for better readability for ~20% time cost  
> direct iteration is always faster

![iteration graph](assets/BenchmarkIterate.png)

[Benchmark results](assets/BenchmarkIterate.txt)
