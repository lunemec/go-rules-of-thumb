## Range over func

With [Go 1.23 came new feature - range over func](https://go.dev/blog/range-functions), lets check when it makes sense to use that over
pre-allocating a slice and putting values in it.

> [!TIP]
> use direct iteration on hot paths  
> use `iter.Seq` when it makes the API or call site cleaner  
> expect about ~10-20% overhead on medium and large loops, and more on tiny ones

![iteration graph](assets/BenchmarkIterate.png)

In this benchmark, direct iteration wins at every tested size. `range over func` settles around `12-13%` overhead on medium and large loops, but the penalty is much higher on tiny loops, so the readability trade-off is real but measurable.

[Benchmark results](assets/BenchmarkIterate.txt)
