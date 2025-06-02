## Needle in a haystack

When is it more efficient to convert a _slice_ into a _map_ for locating an element `x` within the set `A` (x ∈ A)?

> [!IMPORTANT]
> use `slice` if `len(neeldes) <= 10`

> [!HINT]
> use `map` when `len(haystack) > 100 && len(needles) > 100`

> **TL;DR**: use `slice` if `len(neeldes) <= 10`  
> use `map` when `len(haystack) > 100 && len(needles) > 100`

Depending on size of the _haystack_ (size) and number of _needles_ (iterations), this will differ:
![needle in a haystack graph](assets/BenchmarkNeedleInAHaystack.png)

[Benchmark results](assets/BenchmarkNeedleInAHaystack.txt)
