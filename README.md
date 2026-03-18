# Rules of thumb for Go

> In English, the phrase "rule of thumb" refers to an approximate method for doing something, based on practical experience rather than theory.

As a software engineer, you likely have a good understanding of data structures and the `Big O` complexities associated with different usage patterns. However, determining the most suitable data structure for your specific use case can be a challenging decision.

Often, you'll find yourself in a situation where you need to weigh the benefits of creating a new data structure optimized for your access pattern. The question then arises: when is it worthwhile to invest the effort in building a custom data structure?

The decision isn't as simple as solely relying on `Big O` notation, which primarily reflects time complexity. Real-world performance depends on various factors, such as memory locality, the number of allocations, pointer chasing, and more.

## General rules of thumb that always apply

Always KISS (keep it simple, stupid).

1. make it work
2. make it right
3. make it fast

## Disclaimer

❗These rules are not a dogma! Please don't link to this document saying "you should use this because rules-of-thumb says so". Always measure and benchmark your own code with your own data.

Examples here are what is called micro-optimization, before diving into these, profile your code, find real bottlenecks, and fix low hanging fruit there first.
## Needle in a haystack

When is it more efficient to convert a _slice_ into a _map_ for locating an element `x` within the set `A` (x ∈ A)?

> [!TIP]
> use `slice` for one-off checks and up to ~50 lookups  
> switch to `map` when you are doing hundreds of lookups on the same haystack  
> between ~50 and ~100 lookups, benchmark your real workload

Depending on size of the _haystack_ (size) and number of _needles_ (iterations), this will differ:
![needle in a haystack graph](assets/BenchmarkNeedleInAHaystack.png)

[Detailed line view](assets/BenchmarkNeedleInAHaystack-detail.png)

In this benchmark, `slice` wins every `10`-lookup case and still wins much of the `50`-lookup region. `map` takes over most of the grid once lookups reach `100+`, but the exact crossover still depends on both the haystack size and how many lookups you amortize the map build across.

[Benchmark results](assets/BenchmarkNeedleInAHaystack.txt)
## Deduplication

When is it more efficient to deduplicate a `slice` as opposed to using a `map[]struct{}` for the same purpose?

> [!TIP]
> if order does not matter, use in-place sort + dedup up to ~1000 items  
> use `map` from roughly ~5000 items upward  
> if you must preserve original order, use `map`

![deduplication graph](assets/BenchmarkDeduplication.png)

In this benchmark, in-place sort + dedup is the fastest option from `10` through `1000` items, while `map` takes over from `5000` onward. The plain slice scan is never the fastest path here; it is a simplicity choice for very small inputs, not a performance choice.

[Benchmark results](assets/BenchmarkDeduplication.txt)
## Subsets

When checking if **A** is subset of **B** (A ⊆ B), when is it more efficient to iterate both slices in nested loop `A x B` `O(n^2)`, and when does it make sense to use `map`, or `sort` + binary search?  
Meaning of `A ⊆ B` in this test is that _all_ elements of **A** are present in **B**, regardless of position.

> [!TIP]
> if `len(A) <= 100`, start with nested loops  
> if `len(A) >= 500 && len(B) >= 1000`, use `map`  
> use `sort + binary search` only in the middle, or when `B` is already sorted

![subsets graph](assets/BenchmarkSubset.png)

[Detailed line view](assets/BenchmarkSubset-detail.png)

The measured crossover is mostly driven by the size of `A`: small subsets keep the nested loop competitive for surprisingly long, while `map` dominates once both sides are non-trivial. `sort + binary search` only wins a narrow middle band in this benchmark, so it is best treated as a special-case option rather than a default.

[Benchmark results](assets/BenchmarkSubset.txt)
## Append

> [!TIP]
> use `append(dst, src...)` as the default  
> if `len(src)` is comparable to or larger than `len(dst)` and this is hot code, preallocate the full result  
> avoid `for` + `append` without preallocation

![append graph](assets/BenchmarkAppend.png)

In this benchmark, `append(dst, src...)` wins most of the grid, especially when the appended slice is small. Once the appended slice gets large relative to the destination, the preallocated indexed copy often pulls ahead, so "append is always fastest" is too strong a rule.

[Detailed line view](assets/BenchmarkAppend-detail.png)

[Benchmark results](assets/BenchmarkAppend.txt)
## Strings concatenation

Is it more efficient to `"str1" + var`, `fmt.Sprintf()`, `strings.Join()` or `strings.Builder`? When does it make sense to add `sync.Pool`?

> [!TIP]
> for repeated concatenation, start with `strings.Builder`  
> benchmark `sync.Pool + strings.Builder` once the loop gets very hot or reaches ~1000+ concatenations per operation  
> use `+` for one-off expressions and `fmt.Sprintf` for formatting, not concat speed
>
> in this benchmark, only the `strings.Builder` variants ever win

![concatenation graph](assets/BenchmarkConcat.png)

[Detailed line view](assets/BenchmarkConcat-detail.png)

In this benchmark, only the two `strings.Builder` variants take first place. `sync.Pool + strings.Builder` starts to win more often as the work gets heavier, but plain `strings.Builder` stays competitive across the whole grid, which makes it the safest default.

[Benchmark results](assets/BenchmarkConcat.txt)
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
## Pass by reference vs copy

When should you pass a reference (pointer), and when should you use pass by value?

> [!TIP]
> use pointers when you need mutation  
> for read-only data, start with the simpler API and measure  
> this benchmark does not show a reliable universal size cutoff

In this benchmark, passing a pointer wins for all tested struct sizes, and the gap grows as the copied array gets larger. That is still a narrow microbenchmark, so the safe rule is not "always use pointers", but "measure once copying large values shows up in a profile".

References (pointers) vs copied values are still way more complicated than one synthetic test can capture, and there is tons of resources on this topic. A great one is
[this article](https://dave.cheney.net/2017/04/29/there-is-no-pass-by-reference-in-go) by Dave Cheney.
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
## Notes

- More "Rules of thumb" will be added over time.
- All benchmarks were conducted on a **Macbook Pro M1 (2020) 16GB RAM**, using **Go 1.24.3**.
