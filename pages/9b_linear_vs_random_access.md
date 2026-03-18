## Locality: linear vs randomized access

When does it pay to keep a traversal contiguous instead of visiting the same records in a fixed random order?

This section is about locality and CPU prefetch behavior over the same `O(n)` work, not about changing algorithmic complexity.

> [!TIP]
> for `64B` records, fixed-random order stays ahead through ~`2 048` records in this run
> linear order takes over around ~`4 096` records and keeps widening from there
> by `65 536` records, linear is about `44%` faster here
> treat the `2 048-4 096` crossover as hardware-sensitive and benchmark on your own CPU

![linear vs random access graph](assets/BenchmarkLinearVsRandomAccess.png)

This benchmark prebuilds one `[]record` plus two index orders: identity and a fixed permutation. Both variants read the same `64B` records exactly once per pass, sum the same hot fields, allocate nothing, and differ only in access order. In this run, the fixed random walk was about `9-19%` faster from `64` through `2 048` records, then linear pulled ahead at `4 096` and widened to about `13%` at `8 192`, `29%` at `32 768`, and `44%` at `65 536`. The practical rule is not that random access is "better"; locality effects can flip at small working sets, but contiguous layout and traversal order become increasingly valuable once the walk grows past the smallest caches.

[Benchmark results](assets/BenchmarkLinearVsRandomAccess.txt)
