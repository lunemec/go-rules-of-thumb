## If vs switch

Is there even any difference? In theory, `switch` should be faster (at least for some types) if the
compiler is able to transform it into a jump table.

> **TL;DR**: Use which ever one is more readable.

![if switch graph](assets/BenchmarkIfSwitch.png)

[Benchmark results](assets/BenchmarkIfSwitch.txt)

It looks like Go doesn't support jump tables yet? The tests I tried compile into same code for both switch/if statements. You can try to hand-roll jump table [similar to the #19791](https://github.com/golang/go/issues/19791).

Read more:

- <https://github.com/golang/go/issues/5496>
- <https://github.com/golang/go/issues/19791>
- <https://github.com/golang/go/issues/10870>
- <https://go-review.googlesource.com/c/go/+/357330>
- <https://go-review.googlesource.com/c/go/+/395714>
