## Range over func

With [Go 1.23 came new feature - range over func](https://go.dev/blog/range-functions), lets check when it makes sense to use that over
pre-allocating a slice and putting values in it.

> [!TIP]
> use direct iteration on hot paths  
> use `iter.Seq` when it makes the API or call site cleaner  
> expect `range over func` to stay close on tiny loops and cost about ~15-25% on larger ones  
> materializing a slice is noticeably more expensive because it also pays the slice build cost  

![iteration graph](assets/BenchmarkIterate.png)

In this benchmark, direct iteration still wins at every tested size. `range over func` is effectively a wash on the tiniest loops, then settles around `15-25%` overhead once the loop gets large enough for iterator machinery to show up. Prebuilding a slice is consistently the slowest option here because it pays both the generation work and the slice materialization cost.

[Benchmark results](assets/BenchmarkIterate.txt)
