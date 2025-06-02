## Range over func

With [Go 1.23 came new feature - range over func](https://go.dev/blog/range-functions), lets check when it makes sense to use that over
pre-allocating a slice and putting values in it.

I'm quite suprised to see that range over func adds extra 3x number of allocations somewhere.
Not sure where, that is to be measured later.

![iteration graph](assets/BenchmarkIterate.png)

[Benchmark results](assets/BenchmarkIterate.txt)
