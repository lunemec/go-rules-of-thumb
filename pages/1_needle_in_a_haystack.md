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
