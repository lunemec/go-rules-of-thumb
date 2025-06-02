## Append

> [!TIP]
> ALWAYS use `append([]T, elems...)` because `for` looping may trigger multiple array re-sizings, whereas `append` will always allocate only once  
> if you must use `for` loop (extra logic), try to pre-allocate the slice

![append graph](assets/BenchmarkAppend.png)

Even though regular `append()` has time complexity `O(1)` (amortized constant-time), because every time it needs to allocate more space, it grows the underlying data array by 2x (until 512 elements, after 512 it grows less), simply by having to allocate + copy makes it significantly slower than if you are able to calculate the resulting size and pre-allocating.

[Benchmark results](assets/BenchmarkAppend.txt)
