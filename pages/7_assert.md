## Negative space programming (a.k.a. asserts everywhere)

What is the cost of adding `assert`? Does it make any significant impact?

> [!TIP]
> use `assert` freely outside hot loops  
> in hot loops, plain `assert` is near-free below ~10 checks and noticeable around ~100+ checks  
> avoid `defer`-based asserts in hot loops

![assert graph](assets/BenchmarkAssert.png)

The absolute times are still small, but the relative cost shows up clearly in a tight loop: plain `assert` is about `1.01x` at `1`-`10` checks and about `2.1x` by `100`-`1000` checks, while `defer`-based asserts are worse. That makes direct asserts fine for most code, but worth avoiding in very hot inner loops.

[Benchmark results](assets/BenchmarkAssert.txt)

Read more:

- <https://spinroot.com/gerard/pdf/P10.pdf>
- <https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md#safety>
- <https://alfasin.com/2017/12/21/negative-space-and-how-does-it-apply-to-coding/>
