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
