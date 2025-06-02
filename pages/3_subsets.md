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
