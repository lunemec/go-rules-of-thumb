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
> use `slice` if `len(neeldes) <= 10`  
> use `map` when `len(haystack) > 100 && len(needles) > 100`

Depending on size of the _haystack_ (size) and number of _needles_ (iterations), this will differ:
![needle in a haystack graph](assets/BenchmarkNeedleInAHaystack.png)

[Benchmark results](assets/BenchmarkNeedleInAHaystack.txt)
## Deduplication

When is it more efficient to deduplicate a `slice` as opposed to using a `map[]struct{}` for the same purpose?

> [!TIP]
> use `map` when `len(haystack) > 100`.  
> if you must reduce allocations, use in-place sort + deduplication

![deduplication graph](assets/BenchmarkDeduplication.png)

[Benchmark results](assets/BenchmarkDeduplication.txt)
## Subsets

When checking if A is subset of B (A ⊆ B), when is it more efficient to iterate both slices in nested loop `A x B` `O(n^2)`, and when does it make sense to use `map`, or `sort` + binary search?

> [!TIP]
> use `slice` when `len(A) << len(B)`  
> use `map` when `len(A) > 500 && len(B) > 500`

![subsets graph](assets/BenchmarkSubset.png)

[Benchmark results](assets/BenchmarkSubset.txt)
## Append

> [!TIP]
> ALWAYS use `append([]T, elems...)` because `for` looping may trigger multiple array re-sizings, whereas `append` will always allocate only once  
> if you must use `for` loop (extra logic), try to pre-allocate the slice

![append graph](assets/BenchmarkAppend.png)

Even though regular `append()` has time complexity `O(1)` (amortized constant-time), because every time it needs to allocate more space, it grows the underlying data array by 2x (until 512 elements, after 512 it grows less), simply by having to allocate + copy makes it significantly slower than if you are able to calculate the resulting size and pre-allocating.

[Benchmark results](assets/BenchmarkAppend.txt)
## Strings concatenation

Is it more efficient to `"str1" + var`, `fmt.Sprintf()`, `strings.Join()` or `strings.Builder`? When does it make sense to add `sync.Pool`?

> [!TIP]
> use `strings.Builder` when `len(str) < 100 & N ops < 1000`  
> use `sync.Pool + strings.Builder` when doing this for every request  
> for `len(str) > 100` use `+` or `strings.Join`
>
> use `fmt.Sprintf` for regular string formatting (not just concatenation)

![concatenation graph](assets/BenchmarkConcat.png)

[Benchmark results](assets/BenchmarkConcat.txt)
## If vs switch

Is there even any difference? In theory, `switch` should be faster (at least for some types) if the
compiler is able to transform it into a jump table.

> [!TIP]
> use which ever one is more readable  
> but `switch` is tiny bit slower

![if switch graph](assets/BenchmarkIfSwitch.png)

[Benchmark results](assets/BenchmarkIfSwitch.txt)

It looks like Go doesn't support jump tables yet? The tests I tried compile into same code for both switch/if statements. You can try to hand-roll jump table [similar to the #19791](https://github.com/golang/go/issues/19791).

Read more:

- <https://github.com/golang/go/issues/5496>
- <https://github.com/golang/go/issues/19791>
- <https://github.com/golang/go/issues/10870>
- <https://go-review.googlesource.com/c/go/+/357330>
- <https://go-review.googlesource.com/c/go/+/395714>
## Negative space programming (a.k.a. asserts everywhere)

What is the cost of adding `assert`? Does it make any significant impact?

> [!TIP]
> use `assert` whenever possible to improve reliability of your software

![assert graph](assets/BenchmarkAssert.png)

[Benchmark results](assets/BenchmarkAssert.txt)

Read more:

- <https://spinroot.com/gerard/pdf/P10.pdf>
- <https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md#safety>
- <https://alfasin.com/2017/12/21/negative-space-and-how-does-it-apply-to-coding/>
## Pass by reference vs copy

When should you pass a reference (pointer), and when should you use pass by value?

> [!TIP]
> pass by reference if you want to mutate the data, otherwise pass a copy

Performance-wise, this one is almost impossible to give general advice for. If your struct (or nested structs)
are very big (it depends on the types of fields too), copying will become slower.
But if you have many more pointers, you increase GC pressure and your program will
spend more time on waiting on memory pointer lookup.

References (pointers) vs copied values is way more complicated,
and there is tons of resources on this topic, great one is
[this article](https://dave.cheney.net/2017/04/29/there-is-no-pass-by-reference-in-go) by Dave Cheney.
## Range over func

With [Go 1.23 came new feature - range over func](https://go.dev/blog/range-functions), lets check when it makes sense to use that over
pre-allocating a slice and putting values in it.

> [!TIP]
> use `iter.Seq`` for better readability for ~20% time cost  
> direct iteration is always faster

![iteration graph](assets/BenchmarkIterate.png)

[Benchmark results](assets/BenchmarkIterate.txt)
## Notes

- More "Rules of thumb" will be added over time.
- All benchmarks were conducted on a **Macbook Pro M1 (2020) 16GB RAM**, using **Go 1.24.3**.
