## Hot/cold split: inline `T` vs `*T` field

Should a mostly-cold sub-struct live inline as `T`, or behind `*T` inside a larger record?

This section is the benchmark-backed answer to "should this field be `T` or `*T` inside a struct?" Here, "hot" means the small set of fields the fast path reads every time, "cold" means bulky fields that are usually present but ignored by that path, and `snapshot` means "build and copy a full output record", not "take a runtime snapshot". The result depends on access pattern: hot-prefix scans and whole-record copies want different layouts.

> [!TIP]
> keep the field inline when callers usually read or copy the whole record  
> split to `*Cold` only when a hot path scans large collections and mostly ignores the cold tail  
> in this benchmark, the split layout starts to win around `5_000` records on the hot-only scan and is clearly better by `10_000+`  
> for full-record snapshots, inline stays better across the whole measured range

![hot cold split graph](assets/BenchmarkHotColdSplit.png)

This benchmark stores the same records either as one wide inline struct or as a small hot struct pointing at a contiguous cold backing slice. The `hot_scan` row answers "what if the loop only reads the hot prefix and never touches the cold tail?", while the `snapshot` row answers "what if the code needs to assemble and copy the whole record?". In these results, inline is slightly better through `1_000` records on the hot scan, `Hot + *Cold` starts to edge ahead around `5_000`, stays about `7%` faster at `10_000-50_000`, and is about `56%` faster at `100_000`. The snapshot path goes the other way: inline is effectively tied at `1` record and then stays about `5-13%` faster from `10` upward because the full record is already contiguous.

[Benchmark results](assets/BenchmarkHotColdSplit.txt)
