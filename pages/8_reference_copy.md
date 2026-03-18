## Read-only parameter passing: `T` vs `*T`

Should a read-only call boundary take a large payload as `T` or `*T`?

This benchmark is the representative parameter-passing case for this repo. It stands in for either a plain function parameter or a method receiver with the same read-only call shape.
Here, "hot fields" just means the small set of fields the callee actually reads on every call; the benchmark is not modeling whole-record copies.

> [!TIP]
> use `T` up to about `16B`
> around `24-32B`, benchmark your own workload
> on this benchmark, prefer `*T` from about `32B` upward for read-only hot paths
> keep `T` when you specifically want value semantics or isolation

![param value vs pointer graph](assets/BenchmarkParamValueVsPointer.png)

Each benchmark operation runs `256` `//go:noinline` read-only calls over aligned mixed-field structs from `8B` to `512B`, and each call reads only a few scalar fields into a sink accumulator. In these results, `T` is about `6.6%` faster at `8B`, `16B` is effectively a wash, `*T` is about `11%` faster at `24B`, and the gap grows from about `31%` at `32B` to about `195%` at `512B`, so the measured crossover for this call shape is around `24-32B`. This does not measure mutation, whole-record copying, interface dispatch, slice layout, or GC-heavy escaping.

[Benchmark results](assets/BenchmarkParamValueVsPointer.txt)

Further reading:
- [There is no pass-by-reference in Go](https://dave.cheney.net/2017/04/29/there-is-no-pass-by-reference-in-go)
