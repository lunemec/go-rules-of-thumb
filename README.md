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

> **TL;DR**: use `map` when `len(haystack) > 100 && len(needles) > 100`

![needle in a haystack graph](assets/BenchmarkNeedleInAHaystack.png "Find element in set benchmark")
Depending on size of the _haystack_ (size) and number of _needles_ (iterations), this will differ:

|                                                  | assets/BenchmarkNeedleInAHaystack-map.txt |     | assets/BenchmarkNeedleInAHaystack-slice.txt |     |           |              |
| ------------------------------------------------ | ----------------------------------------- | --- | ------------------------------------------- | --- | --------- | ------------ |
|                                                  | sec/op                                    | CI  | sec/op                                      | CI  | vs base   | P            |
| NeedleInAHaystack/size=10_iterations=10-8        | 2.0195000000000003e-07                    | 0%  | 5.218e-08                                   | 0%  | -74.16%   | p=0.000 n=10 |
| NeedleInAHaystack/size=10_iterations=100-8       | 6.858000000000001e-07                     | 0%  | 3.5580000000000005e-07                      | 0%  | -48.12%   | p=0.000 n=10 |
| NeedleInAHaystack/size=10_iterations=500-8       | 2.849e-06                                 | 0%  | 1.5585e-06                                  | 0%  | -45.30%   | p=0.000 n=10 |
| NeedleInAHaystack/size=10_iterations=1000-8      | 5.527500000000001e-06                     | 0%  | 3.475e-06                                   | 0%  | -37.13%   | p=0.000 n=10 |
| NeedleInAHaystack/size=100_iterations=10-8       | 1.1250000000000002e-06                    | 0%  | 2.54e-07                                    | 0%  | -77.42%   | p=0.000 n=10 |
| NeedleInAHaystack/size=100_iterations=100-8      | 1.6680000000000002e-06                    | 0%  | 2.1815e-06                                  | 0%  | +30.79%   | p=0.000 n=10 |
| NeedleInAHaystack/size=100_iterations=500-8      | 3.8165e-06                                | 0%  | 1.0569e-05                                  | 0%  | +176.93%  | p=0.000 n=10 |
| NeedleInAHaystack/size=100_iterations=1000-8     | 6.5995e-06                                | 0%  | 2.2282500000000002e-05                      | 0%  | +237.64%  | p=0.000 n=10 |
| NeedleInAHaystack/size=500_iterations=10-8       | 6.55e-06                                  | 0%  | 1.1410000000000002e-06                      | 1%  | -82.58%   | p=0.000 n=10 |
| NeedleInAHaystack/size=500_iterations=100-8      | 7.164500000000001e-06                     | 0%  | 1.1229500000000001e-05                      | 1%  | +56.74%   | p=0.000 n=10 |
| NeedleInAHaystack/size=500_iterations=500-8      | 9.369000000000001e-06                     | 0%  | 5.4998000000000004e-05                      | 1%  | +487.02%  | p=0.000 n=10 |
| NeedleInAHaystack/size=500_iterations=1000-8     | 1.20535e-05                               | 0%  | 0.00011133700000000001                      | 1%  | +823.69%  | p=0.000 n=10 |
| NeedleInAHaystack/size=1000_iterations=10-8      | 1.51545e-05                               | 1%  | 2.164e-06                                   | 1%  | -85.72%   | p=0.000 n=10 |
| NeedleInAHaystack/size=1000_iterations=100-8     | 1.58995e-05                               | 0%  | 2.08445e-05                                 | 1%  | +31.10%   | p=0.000 n=10 |
| NeedleInAHaystack/size=1000_iterations=500-8     | 1.7502999999999998e-05                    | 1%  | 0.000105146                                 | 1%  | +500.73%  | p=0.000 n=10 |
| NeedleInAHaystack/size=1000_iterations=1000-8    | 2.06655e-05                               | 0%  | 0.0002072555                                | 1%  | +902.91%  | p=0.000 n=10 |
| NeedleInAHaystack/size=5000_iterations=10-8      | 7.94045e-05                               | 0%  | 1.01825e-05                                 | 1%  | -87.18%   | p=0.000 n=10 |
| NeedleInAHaystack/size=5000_iterations=100-8     | 7.982400000000001e-05                     | 0%  | 0.000102442                                 | 1%  | +28.33%   | p=0.000 n=10 |
| NeedleInAHaystack/size=5000_iterations=500-8     | 8.2366e-05                                | 0%  | 0.0005082955                                | 3%  | +517.12%  | p=0.000 n=10 |
| NeedleInAHaystack/size=5000_iterations=1000-8    | 8.38165e-05                               | 0%  | 0.0010191010000000001                       | 1%  | +1115.87% | p=0.000 n=10 |
| NeedleInAHaystack/size=10000_iterations=10-8     | 0.00016355150000000001                    | 0%  | 2.01805e-05                                 | 1%  | -87.66%   | p=0.000 n=10 |
| NeedleInAHaystack/size=10000_iterations=100-8    | 0.000164117                               | 0%  | 0.00020134900000000003                      | 1%  | +22.69%   | p=0.000 n=10 |
| NeedleInAHaystack/size=10000_iterations=500-8    | 0.00016748800000000003                    | 0%  | 0.0010075785                                | 2%  | +501.58%  | p=0.000 n=10 |
| NeedleInAHaystack/size=10000_iterations=1000-8   | 0.00016740450000000002                    | 0%  | 0.0020248345000000003                       | 4%  | +1109.55% | p=0.000 n=10 |
| NeedleInAHaystack/size=50000_iterations=10-8     | 0.000853139                               | 0%  | 0.00010159900000000001                      | 0%  | -88.09%   | p=0.000 n=10 |
| NeedleInAHaystack/size=50000_iterations=100-8    | 0.000859315                               | 0%  | 0.0010305635000000001                       | 2%  | +19.93%   | p=0.000 n=10 |
| NeedleInAHaystack/size=50000_iterations=500-8    | 0.0008550655000000001                     | 0%  | 0.005007204                                 | 4%  | +485.59%  | p=0.000 n=10 |
| NeedleInAHaystack/size=50000_iterations=1000-8   | 0.0008593165000000001                     | 0%  | 0.010244186500000002                        | 6%  | +1092.13% | p=0.000 n=10 |
| NeedleInAHaystack/size=100000_iterations=10-8    | 0.0016933565000000002                     | 1%  | 0.00020292850000000001                      | 1%  | -88.02%   | p=0.000 n=10 |
| NeedleInAHaystack/size=100000_iterations=100-8   | 0.0016954315                              | 0%  | 0.002017895                                 | 3%  | +19.02%   | p=0.000 n=10 |
| NeedleInAHaystack/size=100000_iterations=500-8   | 0.0016948945000000001                     | 0%  | 0.0100737875                                | 9%  | +494.36%  | p=0.000 n=10 |
| NeedleInAHaystack/size=100000_iterations=1000-8  | 0.001696308                               | 0%  | 0.020756502500000003                        | 7%  | +1123.63% | p=0.000 n=10 |
| NeedleInAHaystack/size=500000_iterations=10-8    | 0.015493154500000002                      | 3%  | 0.0010208700000000001                       | 3%  | -93.41%   | p=0.000 n=10 |
| NeedleInAHaystack/size=500000_iterations=100-8   | 0.015549726000000002                      | 3%  | 0.010098937                                 | 8%  | -35.05%   | p=0.000 n=10 |
| NeedleInAHaystack/size=500000_iterations=500-8   | 0.015741979                               | 2%  | 0.052639357                                 | 7%  | +234.39%  | p=0.000 n=10 |
| NeedleInAHaystack/size=500000_iterations=1000-8  | 0.015602500500000002                      | 4%  | 0.09976974450000001                         | 21% | +539.45%  | p=0.000 n=10 |
| NeedleInAHaystack/size=1000000_iterations=10-8   | 0.0418112905                              | 2%  | 0.0020738275                                | 2%  | -95.04%   | p=0.000 n=10 |
| NeedleInAHaystack/size=1000000_iterations=100-8  | 0.0415576295                              | 1%  | 0.020849485                                 | 8%  | -49.83%   | p=0.000 n=10 |
| NeedleInAHaystack/size=1000000_iterations=500-8  | 0.0417492095                              | 2%  | 0.0975659975                                | 12% | +133.70%  | p=0.000 n=10 |
| NeedleInAHaystack/size=1000000_iterations=1000-8 | 0.041954636000000003                      | 1%  | 0.1997639465                                | 26% | +376.14%  | p=0.000 n=10 |
| geomean                                          | 0.00015047703447270417                    |     | 0.00019927621720009707                      |     | +32.43%   |              |
## Deduplication

When is it more efficient to deduplicate a `slice` as opposed to using a `map[]struct{}` for the same purpose?

> **TL;DR**: use `map` when `len(haystack) > 100`. If you must reduce allocations, use sort + in-place slice
> deduplication. Suprisingly it is fast enough.

![deduplication graph](assets/BenchmarkDeduplication.png "Deduplication variants performance")

```
                              │ BenchmarkDeduplication-map.txt │ BenchmarkDeduplication-slice_sort_inplace.txt │       BenchmarkDeduplication-slice.txt       │
                              │             sec/op             │        sec/op         vs base                 │     sec/op      vs base                      │
Deduplication/size=10-8                          205.55n ±  1%           55.48n ±  2%   -73.01% (p=0.000 n=10)     114.40n ± 2%      -44.34% (p=0.000 n=10)
Deduplication/size=100-8                         1563.0n ±  1%           808.6n ±  1%   -48.26% (p=0.000 n=10)     1262.5n ± 1%      -19.23% (p=0.000 n=10)
Deduplication/size=500-8                          9.117µ ±  1%           4.903µ ±  4%   -46.22% (p=0.000 n=10)     31.617µ ± 8%     +246.81% (p=0.000 n=10)
Deduplication/size=1000-8                         19.56µ ±  0%           10.91µ ±  4%   -44.23% (p=0.000 n=10)     108.98µ ± 2%     +457.19% (p=0.000 n=10)
Deduplication/size=5000-8                         101.6µ ±  1%           187.1µ ±  2%   +84.21% (p=0.000 n=10)     2432.6µ ± 3%    +2294.79% (p=0.000 n=10)
Deduplication/size=10000-8                        220.5µ ±  5%           449.5µ ±  1%  +103.87% (p=0.000 n=10)     9744.1µ ± 2%    +4319.66% (p=0.000 n=10)
Deduplication/size=50000-8                        1.136m ±  8%           2.782m ±  1%  +144.80% (p=0.000 n=10)    235.701m ± 3%   +20640.07% (p=0.000 n=10)
Deduplication/size=100000-8                       2.208m ± 27%           5.999m ±  1%  +171.71% (p=0.000 n=10)    927.130m ± 1%   +41889.99% (p=0.000 n=10)
Deduplication/size=500000-8                       25.18m ± 19%           34.65m ±  4%   +37.58% (p=0.000 n=10)   23061.71m ± 3%   +91471.94% (p=0.000 n=10)
Deduplication/size=1000000-8                      60.85m ± 12%           73.29m ± 19%   +20.46% (p=0.001 n=10)   92737.14m ± 2%  +152307.91% (p=0.000 n=10)
Deduplication/size=5000000-8                      393.3m ±  8%           387.9m ±  1%         ~ (p=0.123 n=10)
```

```
                              │ BenchmarkDeduplication-map.txt │ BenchmarkDeduplication-slice_sort_inplace.txt │    BenchmarkDeduplication-slice.txt     │
                              │              B/op              │         B/op           vs base                │     B/op       vs base                  │
Deduplication/size=10-8                            488.00 ± 0%              80.00 ± 0%  -83.61% (p=0.000 n=10)     200.00 ± 0%  -59.02% (p=0.000 n=10)
Deduplication/size=100-8                           4136.0 ± 0%              896.0 ± 0%  -78.34% (p=0.000 n=10)     1912.0 ± 0%  -53.77% (p=0.000 n=10)
Deduplication/size=500-8                         26.039Ki ± 0%            4.000Ki ± 0%  -84.64% (p=0.000 n=10)   11.992Ki ± 0%  -53.95% (p=0.000 n=10)
Deduplication/size=1000-8                        52.078Ki ± 0%            8.000Ki ± 0%  -84.64% (p=0.000 n=10)   22.617Ki ± 0%  -56.57% (p=0.000 n=10)
Deduplication/size=5000-8                        224.31Ki ± 0%            40.00Ki ± 0%  -82.17% (p=0.000 n=10)   125.24Ki ± 0%  -44.17% (p=0.000 n=10)
Deduplication/size=10000-8                       448.63Ki ± 0%            80.00Ki ± 0%  -82.17% (p=0.000 n=10)   261.24Ki ± 0%  -41.77% (p=0.000 n=10)
Deduplication/size=50000-8                       1938.5Ki ± 0%            392.0Ki ± 0%  -79.78% (p=0.000 n=10)   1525.2Ki ± 0%  -21.32% (p=0.000 n=10)
Deduplication/size=100000-8                      3877.1Ki ± 0%            784.0Ki ± 0%  -79.78% (p=0.000 n=10)   3237.3Ki ± 0%  -16.50% (p=0.000 n=10)
Deduplication/size=500000-8                      25.681Mi ± 0%            3.820Mi ± 0%  -85.12% (p=0.000 n=10)   16.521Mi ± 0%  -35.67% (p=0.000 n=10)
Deduplication/size=1000000-8                     51.346Mi ± 0%            7.633Mi ± 0%  -85.13% (p=0.000 n=10)   32.888Mi ± 0%  -35.95% (p=0.000 n=10)
Deduplication/size=5000000-8                     220.61Mi ± 0%            38.15Mi ± 0%  -82.71% (p=0.000 n=10)
Deduplication/size=10000000-8                    441.22Mi ± 0%            76.30Mi ± 0%  -82.71% (p=0.000 n=10)
```

```
                              │ BenchmarkDeduplication-map.txt │ BenchmarkDeduplication-slice_sort_inplace.txt │    BenchmarkDeduplication-slice.txt    │
                              │           allocs/op            │      allocs/op        vs base                 │  allocs/op   vs base                   │
Deduplication/size=10-8                             5.000 ± 0%             1.000 ± 0%   -80.00% (p=0.000 n=10)    5.000 ± 0%         ~ (p=1.000 n=10) ¹
Deduplication/size=100-8                            5.000 ± 0%             1.000 ± 0%   -80.00% (p=0.000 n=10)    8.000 ± 0%   +60.00% (p=0.000 n=10)
Deduplication/size=500-8                            5.000 ± 0%             1.000 ± 0%   -80.00% (p=0.000 n=10)   11.000 ± 0%  +120.00% (p=0.000 n=10)
Deduplication/size=1000-8                           7.000 ± 0%             1.000 ± 0%   -85.71% (p=0.000 n=10)   12.000 ± 0%   +71.43% (p=0.000 n=10)
Deduplication/size=5000-8                          19.000 ± 0%             1.000 ± 0%   -94.74% (p=0.000 n=10)   16.000 ± 0%   -15.79% (p=0.000 n=10)
Deduplication/size=10000-8                         35.000 ± 0%             1.000 ± 0%   -97.14% (p=0.000 n=10)   18.000 ± 0%   -48.57% (p=0.000 n=10)
Deduplication/size=50000-8                        131.000 ± 0%             1.000 ± 0%   -99.24% (p=0.000 n=10)   24.000 ± 0%   -81.68% (p=0.000 n=10)
Deduplication/size=100000-8                       259.000 ± 0%             1.000 ± 0%   -99.61% (p=0.000 n=10)   27.000 ± 0%   -89.58% (p=0.000 n=10)
Deduplication/size=500000-8                      2051.000 ± 0%             1.000 ± 0%   -99.95% (p=0.000 n=10)   34.000 ± 3%   -98.34% (p=0.000 n=10)
Deduplication/size=1000000-8                     4099.000 ± 0%             1.000 ± 0%   -99.98% (p=0.000 n=10)   37.000 ± 3%   -99.10% (p=0.000 n=10)
Deduplication/size=5000000-8                    16387.000 ± 0%             1.000 ± 0%   -99.99% (p=0.000 n=10)
Deduplication/size=10000000-8                   32771.000 ± 0%             1.000 ± 0%  -100.00% (p=0.000 n=10)
```
## Subsets

When checking if A is subset of B (A ⊆ B), when is it more efficient to iterate both slices in nested loop `A x B` `O(n^2)`, and when does it make sense to use `map`, or `sort` + binary search?

> **TL;DR**: when use `slice` when `len(A) << len(B)`, use `map` when `len(A) > 500 && len(B) > 500`.

![subsets graph](assets/BenchmarkSubset.png "Subsets variants performance graph")

| Type                 | len(A) | len(B) | ns/op        |     |
| -------------------- | ------ | ------ | ------------ | --- |
| slice                | 10     | 10     | 25.39 ns/op  | ✅  |
| slice_sort_binsearch | 10     | 10     | 193.3 ns/op  |
| map                  | 10     | 10     | 164.7 ns/op  |
| slice                | 10     | 100    | 28.84 ns/op  | ✅  |
| slice_sort_binsearch | 10     | 100    | 2051 ns/op   |
| map                  | 10     | 100    | 2044 ns/op   |
| slice                | 10     | 500    | 28.82 ns/op  | ✅  |
| slice_sort_binsearch | 10     | 500    | 12733 ns/op  |
| map                  | 10     | 500    | 9694 ns/op   |
| slice                | 10     | 1000   | 28.92 ns/op  | ✅  |
| slice_sort_binsearch | 10     | 1000   | 37443 ns/op  |
| map                  | 10     | 1000   | 19409 ns/op  |
| slice                | 100    | 100    | 1550 ns/op   | ✅  |
| slice_sort_binsearch | 100    | 100    | 2661 ns/op   |
| map                  | 100    | 100    | 3441 ns/op   |
| slice                | 100    | 500    | 2040 ns/op   | ✅  |
| slice_sort_binsearch | 100    | 500    | 13988 ns/op  |
| map                  | 100    | 500    | 10596 ns/op  |
| slice                | 100    | 1000   | 2137 ns/op   | ✅  |
| slice_sort_binsearch | 100    | 1000   | 39404 ns/op  |
| map                  | 100    | 1000   | 20172 ns/op  |
| slice                | 500    | 500    | 34328 ns/op  |
| slice_sort_binsearch | 500    | 500    | 20569 ns/op  |
| map                  | 500    | 500    | 16648 ns/op  | ✅  |
| slice                | 500    | 1000   | 36999 ns/op  |
| slice_sort_binsearch | 500    | 1000   | 52155 ns/op  |
| map                  | 500    | 1000   | 25562 ns/op  | ✅  |
| slice                | 1000   | 1000   | 129238 ns/op |
| slice_sort_binsearch | 1000   | 1000   | 75645 ns/op  |
| map                  | 1000   | 1000   | 33859 ns/op  | ✅  |
## Append

```go
append([]T, elems...) // append_expand
```

vs

```go
for _, e := range elems {
    arr = append(arr, e) // append_for
}
```

vs

```go
for _, e := range elems {
    arr = append(arr, e) // append_for_prealloc
}
```

vs

```go
for i, e := range elems {
    arr[i] = e // append_for_index (pre-allocated)
}
```

> **TL;DR**: ALWAYS use `append([]T, elems...)` because `for` looping may trigger multiple array re-sizings, whereas `append` will always allocate only once. If you must use `for` loop (extra logic), try to pre-allocate the slice.

Even though regular `append()` has time complexity `O(1)` (amortized constant-time), because every time it needs to allocate more space, it grows the underlying data array by 2x (until 512 elements, after 512 it grows less), simply by having to allocate + copy makes it significantly slower than if you are able to calculate the resulting size and pre-allocating.

| Type                | len(A) | len(B) | ns/op       | B/op       | allocs/op   |     |
| ------------------- | ------ | ------ | ----------- | ---------- | ----------- | --- |
| append_expand       | 10     | 1000   | 878.7 ns/op | 8192 B/op  | 1 allocs/op | ✅  |
| append_for_index    | 10     | 1000   | 1049 ns/op  | 8192 B/op  | 1 allocs/op |
| append_for_prealloc | 10     | 1000   | 1148 ns/op  | 8192 B/op  | 1 allocs/op |
| append_for          | 10     | 1000   | 2115 ns/op  | 19936 B/op | 7 allocs/op |
## Strings concatenation

Is it more efficient to `"str1" + var`, `fmt.Sprintf()`, `strings.Join()` or `strings.Builder`? When does it make sense to add `sync.Pool`?

> **TL;DR**: use `strings.Builder` when `len(str) < 100 & N ops < 1000`, use `sync.Pool + strings.Builder` when doing this for every request. For `len(str) > 100` use `+` or `strings.Join`.
>
> Use `fmt.Sprintf` for regular string formatting (not just concatenation).

| Type                 | len(str) | N ops | ns/op         |     |
| -------------------- | -------- | ----- | ------------- | --- |
| plus_sign            | 10       | 10    | 377.2 ns/op   |     |
| sprintf              | 10       | 10    | 1174 ns/op    |
| strings_join         | 10       | 10    | 457.2 ns/op   |     |
| strings_builder      | 10       | 10    | 226.6 ns/op   | ✅  |
| strings_builder_pool | 10       | 10    | 242.4 ns/op   |     |
| plus_sign            | 10       | 100   | 4976 ns/op    |
| sprintf              | 10       | 100   | 13242 ns/op   |
| strings_join         | 10       | 100   | 5832 ns/op    |
| strings_builder      | 10       | 100   | 1275 ns/op    |
| strings_builder_pool | 10       | 100   | 1265 ns/op    | ✅  |
| plus_sign            | 10       | 500   | 62458 ns/op   |
| sprintf              | 10       | 500   | 106616 ns/op  |
| strings_join         | 10       | 500   | 67195 ns/op   |
| strings_builder      | 10       | 500   | 6515 ns/op    | ✅  |
| strings_builder_pool | 10       | 500   | 6670 ns/op    |
| plus_sign            | 10       | 1000  | 209530 ns/op  |
| sprintf              | 10       | 1000  | 308757 ns/op  |
| strings_join         | 10       | 1000  | 219302 ns/op  |
| strings_builder      | 10       | 1000  | 13754 ns/op   |
| strings_builder_pool | 10       | 1000  | 13660 ns/op   | ✅  |
| plus_sign            | 100      | 10    | 533.6 ns/op   |     |
| sprintf              | 100      | 10    | 1342 ns/op    |
| strings_join         | 100      | 10    | 586.6 ns/op   |     |
| strings_builder      | 100      | 10    | 975.8 ns/op   |     |
| strings_builder_pool | 100      | 10    | 1028 ns/op    |
| plus_sign            | 100      | 100   | 6670 ns/op    | ✅  |
| sprintf              | 100      | 100   | 14949 ns/op   |
| strings_join         | 100      | 100   | 7562 ns/op    |
| strings_builder      | 100      | 100   | 9713 ns/op    |
| strings_builder_pool | 100      | 100   | 9918 ns/op    |
| plus_sign            | 100      | 500   | 71459 ns/op   |
| sprintf              | 100      | 500   | 116144 ns/op  |
| strings_join         | 100      | 500   | 75915 ns/op   |
| strings_builder      | 100      | 500   | 41344 ns/op   | ✅  |
| strings_builder_pool | 100      | 500   | 43655 ns/op   |
| plus_sign            | 100      | 1000  | 227323 ns/op  |
| sprintf              | 100      | 1000  | 519672 ns/op  |
| strings_join         | 100      | 1000  | 316674 ns/op  |
| strings_builder      | 100      | 1000  | 93747 ns/op   | ✅  |
| strings_builder_pool | 100      | 1000  | 99357 ns/op   |
| plus_sign            | 500      | 10    | 1900 ns/op    | ✅  |
| sprintf              | 500      | 10    | 2741 ns/op    |
| strings_join         | 500      | 10    | 2066 ns/op    |
| strings_builder      | 500      | 10    | 7532 ns/op    |
| strings_builder_pool | 500      | 10    | 6380 ns/op    |
| plus_sign            | 500      | 100   | 20481 ns/op   | ✅  |
| sprintf              | 500      | 100   | 30771 ns/op   |
| strings_join         | 500      | 100   | 21401 ns/op   |
| strings_builder      | 500      | 100   | 68556 ns/op   |
| strings_builder_pool | 500      | 100   | 69828 ns/op   |
| plus_sign            | 500      | 500   | 139848 ns/op  | ✅  |
| sprintf              | 500      | 500   | 193818 ns/op  |
| strings_join         | 500      | 500   | 143040 ns/op  |
| strings_builder      | 500      | 500   | 313665 ns/op  |
| strings_builder_pool | 500      | 500   | 322586 ns/op  |
| plus_sign            | 500      | 1000  | 370027 ns/op  | ✅  |
| sprintf              | 500      | 1000  | 515227 ns/op  |
| strings_join         | 500      | 1000  | 379155 ns/op  |
| strings_builder      | 500      | 1000  | 647395 ns/op  |
| strings_builder_pool | 500      | 1000  | 573904 ns/op  |
| plus_sign            | 1000     | 10    | 3220 ns/op    | ✅  |
| sprintf              | 1000     | 10    | 4613 ns/op    |
| strings_join         | 1000     | 10    | 3342 ns/op    |
| strings_builder      | 1000     | 10    | 13307 ns/op   |
| strings_builder_pool | 1000     | 10    | 13747 ns/op   |
| plus_sign            | 1000     | 100   | 40945 ns/op   |
| sprintf              | 1000     | 100   | 49810 ns/op   |
| strings_join         | 1000     | 100   | 36396 ns/op   | ✅  |
| strings_builder      | 1000     | 100   | 142254 ns/op  |
| strings_builder_pool | 1000     | 100   | 149184 ns/op  |
| plus_sign            | 1000     | 500   | 224290 ns/op  | ✅  |
| sprintf              | 1000     | 500   | 296963 ns/op  |
| strings_join         | 1000     | 500   | 233688 ns/op  |
| strings_builder      | 1000     | 500   | 783015 ns/op  |
| strings_builder_pool | 1000     | 500   | 657683 ns/op  |
| plus_sign            | 1000     | 1000  | 557672 ns/op  | ✅  |
| sprintf              | 1000     | 1000  | 715151 ns/op  |
| strings_join         | 1000     | 1000  | 571112 ns/op  |
| strings_builder      | 1000     | 1000  | 1326209 ns/op |
| strings_builder_pool | 1000     | 1000  | 1106394 ns/op |
## If vs switch

Is there even any difference? In theory, `switch` should be faster (at least for some types) if the
compiler is able to transform it into a jump table.

> **TL;DR**: Use which ever one is more readable.

| Type   | N statements | ns/op        |     |
| ------ | ------------ | ------------ | --- |
| if     | 1            | 0.9470 ns/op |
| switch | 1            | 0.9486 ns/op |
| if     | 5            | 1.270 ns/op  |
| switch | 5            | 1.578 ns/op  |

It looks like Go doesn't support jump tables yet? The tests I tried compile into same code for both switch/if statements. You can try to hand-roll jump table [similar to the #19791](https://github.com/golang/go/issues/19791).

Read more:

- <https://github.com/golang/go/issues/5496>
- <https://github.com/golang/go/issues/19791>
- <https://github.com/golang/go/issues/10870>
- <https://go-review.googlesource.com/c/go/+/357330>
- <https://go-review.googlesource.com/c/go/+/395714>
## Negative space programming (a.k.a. asserts everywhere)

What is the cost of adding `assert`? Does it make any significant impact?

> **TL;DR**: Use asserts whenever possible to improve reliability of your software. The cost is almost non-existent.

| Type         | N statements | ns/op        |     |
| ------------ | ------------ | ------------ | --- |
| no assert    | 1            | 0.3453 ns/op |
| assert       | 1            | 0.4979 ns/op |
| assert       | 5            | 1.791 ns/op  |
| defer assert | 1            | 2.411 ns/op  |

Read more:

- <https://spinroot.com/gerard/pdf/P10.pdf>
- <https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md#safety>
- <https://alfasin.com/2017/12/21/negative-space-and-how-does-it-apply-to-coding/>
## Pass by reference vs copy

When should you pass a reference (pointer), and when should you use pass by value?

> **TL;DR**: Pass by reference if you want to mutate the data, otherwise pass a copy.

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

I'm quite suprised to see that range over func adds extra 3x number of allocations somewhere.
Not sure where, that is to be measured later.

Here are some results:

```
BenchmarkRangeFunc/slice(10)_iterations(10)-8            976581       1175 ns/op      880 B/op       11 allocs/op
BenchmarkRangeFunc/iter_func(10)_iterations(10)-8        770889       1556 ns/op      528 B/op       33 allocs/op
BenchmarkRangeFunc/slice(10)_iterations(100)-8           112356      10747 ns/op     8080 B/op      101 allocs/op
BenchmarkRangeFunc/iter_func(10)_iterations(100)-8        80978      14675 ns/op     4848 B/op      303 allocs/op
BenchmarkRangeFunc/slice(10)_iterations(500)-8            22530      54421 ns/op    40080 B/op      501 allocs/op
BenchmarkRangeFunc/iter_func(10)_iterations(500)-8        16900      70454 ns/op    24048 B/op     1503 allocs/op
BenchmarkRangeFunc/slice(10)_iterations(1000)-8           10000     105975 ns/op    80080 B/op     1001 allocs/op
BenchmarkRangeFunc/iter_func(10)_iterations(1000)-8        8494     140141 ns/op    48048 B/op     3003 allocs/op
BenchmarkRangeFunc/slice(100)_iterations(10)-8           123470       9899 ns/op     9856 B/op       11 allocs/op
BenchmarkRangeFunc/iter_func(100)_iterations(10)-8       116560      10279 ns/op      528 B/op       33 allocs/op
BenchmarkRangeFunc/slice(100)_iterations(100)-8           13488      88975 ns/op    90496 B/op      101 allocs/op
BenchmarkRangeFunc/iter_func(100)_iterations(100)-8       12710      94388 ns/op     4848 B/op      303 allocs/op
BenchmarkRangeFunc/slice(100)_iterations(500)-8            2715     443446 ns/op   448897 B/op      501 allocs/op
BenchmarkRangeFunc/iter_func(100)_iterations(500)-8        2498     492394 ns/op    24048 B/op     1503 allocs/op
BenchmarkRangeFunc/slice(100)_iterations(1000)-8           1345     934711 ns/op   896899 B/op     1001 allocs/op
BenchmarkRangeFunc/iter_func(100)_iterations(1000)-8       1252     975041 ns/op    48048 B/op     3003 allocs/op
BenchmarkRangeFunc/slice(500)_iterations(10)-8            25249      47466 ns/op    45056 B/op       11 allocs/op
BenchmarkRangeFunc/iter_func(500)_iterations(10)-8        24255      49390 ns/op      528 B/op       33 allocs/op
BenchmarkRangeFunc/slice(500)_iterations(100)-8            2761     434421 ns/op   413697 B/op      101 allocs/op
BenchmarkRangeFunc/iter_func(500)_iterations(100)-8        2652     455068 ns/op     4848 B/op      303 allocs/op
BenchmarkRangeFunc/slice(500)_iterations(500)-8             556    2159299 ns/op  2052103 B/op      501 allocs/op
BenchmarkRangeFunc/iter_func(500)_iterations(500)-8         532    2263910 ns/op    24048 B/op     1503 allocs/op
BenchmarkRangeFunc/slice(500)_iterations(1000)-8            277    4309678 ns/op  4100116 B/op     1001 allocs/op
BenchmarkRangeFunc/iter_func(500)_iterations(1000)-8        266    4501908 ns/op    48048 B/op     3003 allocs/op
BenchmarkRangeFunc/slice(1000)_iterations(10)-8           12734      94393 ns/op    90112 B/op       11 allocs/op
BenchmarkRangeFunc/iter_func(1000)_iterations(10)-8       12175      98567 ns/op      528 B/op       33 allocs/op
BenchmarkRangeFunc/slice(1000)_iterations(100)-8           1383     862572 ns/op   827395 B/op      101 allocs/op
BenchmarkRangeFunc/iter_func(1000)_iterations(100)-8       1328     907260 ns/op     4848 B/op      303 allocs/op
BenchmarkRangeFunc/slice(1000)_iterations(500)-8            277    4285715 ns/op  4104203 B/op      501 allocs/op
BenchmarkRangeFunc/iter_func(1000)_iterations(500)-8        267    4479069 ns/op    24048 B/op     1503 allocs/op
BenchmarkRangeFunc/slice(1000)_iterations(1000)-8           139    8939082 ns/op  8200234 B/op     1001 allocs/op
BenchmarkRangeFunc/iter_func(1000)_iterations(1000)-8       127    9272349 ns/op    48048 B/op     3003 allocs/op
```
## Notes

- More "Rules of thumb" will be added over time.
- All benchmarks were conducted on a **Macbook Pro M1 (2020) 16GB RAM**, using **Go 1.24.3**.
