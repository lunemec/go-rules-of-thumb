## If vs switch

Is there even any difference? In theory, `switch` should be faster (at least for some types) if the
compiler is able to transform it into a jump table.

> [!TIP]
> use whichever is more readable  
> when hits are usually the first case, `if` stays competitive and often wins  
> on this benchmark, `switch` usually wins once misses or later/mixed hits are common, especially by `9` cases  

![if switch graph](assets/BenchmarkIfSwitch.png)

[Detailed line view](assets/BenchmarkIfSwitch-detail.png)

[Benchmark results](assets/BenchmarkIfSwitch.txt)

This benchmark compares dense integer equality chains of `3`, `5`, and `9` cases across six deterministic input patterns: `first`, `middle`, `last`, `miss`, `cycle`, and `random`.

There is no single winner. `if` only clearly leads on `first` hits and stays roughly tied on the shorter chains, while `switch` pulls ahead on most `miss`, mixed, and later-hit workloads. By `9` cases, `switch` wins every workload except `first`, often by a wide margin; by `5` cases, the result is already split between `if` on `first` and `switch` on `miss`, `cycle`, and `random`. So the old “longer linear chains favor `if`” conclusion does not hold on these measurements.

This benchmark shows measured behavior for these dense integer chains on this toolchain and hardware. It does not establish which compiler lowering strategy produced the result.

Read more:

- <https://github.com/golang/go/issues/5496>
- <https://github.com/golang/go/issues/19791>
- <https://github.com/golang/go/issues/10870>
- <https://go-review.googlesource.com/c/go/+/357330>
- <https://go-review.googlesource.com/c/go/+/395714>
