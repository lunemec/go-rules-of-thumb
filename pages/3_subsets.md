## Subsets

When checking if **A** is subset of **B** (A ⊆ B), when is it more efficient to iterate both slices in nested loop `A x B` `O(n^2)`, and when does it make sense to use `map`, or `sort` + binary search?  
Meaning of `A ⊆ B` in this test is that _all_ elements of **A** are present in **B**, regardless of position.

> [!TIP]
> use `slice` when `len(A) << len(B)`  
> use `map` when `len(A) > 500 && len(B) > 500`

![subsets graph](assets/BenchmarkSubset.png)

[Benchmark results](assets/BenchmarkSubset.txt)
