## If vs switch

Is there even any difference? In theory, `switch` should be faster (at least for some types) if the
compiler is able to transform it into a jump table.

> [!TIP]
> use whichever is more readable  
> for single-case branches, `if` and `switch` are effectively equal  
> for longer linear chains, `if` is slightly to moderately faster in this benchmark

![if switch graph](assets/BenchmarkIfSwitch.png)

[Benchmark results](assets/BenchmarkIfSwitch.txt)

The 1-case versions are basically a wash here. The 5-case `switch` is consistently slower than the 5-case `if` chain, so the data does not support the idea that `switch` is a free performance win.

It looks like Go does not support jump tables here? The tests I tried compile into same code for both switch/if statements. You can try to hand-roll jump table [similar to the #19791](https://github.com/golang/go/issues/19791).

Read more:

- <https://github.com/golang/go/issues/5496>
- <https://github.com/golang/go/issues/19791>
- <https://github.com/golang/go/issues/10870>
- <https://go-review.googlesource.com/c/go/+/357330>
- <https://go-review.googlesource.com/c/go/+/395714>
