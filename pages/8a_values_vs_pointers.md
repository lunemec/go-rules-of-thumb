## Slice of values vs slice of pointers

Should a read-heavy collection store values (`[]T`) or pointers (`[]*T`)?

This section is about collection layout for read-heavy data, not a blanket rule for API design or parameter passing.
This benchmark has two read patterns: `hot_scan` means walking the collection and reading only a few frequently-used fields, while `snapshot` means building a fresh output slice by copying the full record. It is not a runtime or persistence snapshot.

> [!TIP]
> for wide records with a hot path that only reads a few fields, `[]*T` can win  
> in this benchmark, `[]*T` stays ahead through ~`10 000` records on the hot scan, and `[]T` only pulls ahead around ~`100 000`  
> for snapshot-style reads that copy most of each record, treat the layouts as roughly tied here and benchmark your own workload  
> reach for pointers when you need shared mutation, stable identity, or optional values  

![values vs pointers graph](assets/BenchmarkValuesVsPointers.png)

This benchmark compares the same wide records in `[]T` and in `[]*T` backed by an equivalent contiguous slice, so it isolates pointer indirection without heap-fragmentation noise. The `hot_scan` row answers "what if I mostly read a small hot prefix of each record?", while the `snapshot` row answers "what if I build and copy whole records?". In these results, hot scans favored `[]*T` by about `7-18%` up to `10 000` records, then `[]T` edged ahead by about `4%` at `100 000`; snapshot builds were effectively tied at `10-1 000`, `[]*T` won at `10 000`, and `100 000` was inconclusive, so the real rule of thumb is to match the layout to the read path rather than assume either representation wins in general.

[Benchmark results](assets/BenchmarkValuesVsPointers.txt)

Further reading:
- [CPU Cache-Friendly Data Structures in Go: 10x Speed with Same Algorithm](https://skoredin.pro/blog/golang/cpu-cache-friendly-go)
- [There is no pass-by-reference in Go](https://dave.cheney.net/2017/04/29/there-is-no-pass-by-reference-in-go)
