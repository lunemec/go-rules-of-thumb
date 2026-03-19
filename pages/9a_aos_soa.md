## Array of Structs vs Struct of Arrays

When is it worth splitting a slice of game-style entities into field-parallel slices?

> [!TIP]
> for hot loops over a few fields, use `AoS` when `len(entities) <= 100`  
> for hot loops over a few fields, use `SoA` when `len(entities) >= 1 000`  
> if you usually work with whole records together, keep `AoS`, especially once `len(entities) >= 10 000`  

This benchmark uses a game-style entity model with hot physics fields (`position`, `velocity`, `active`) and cold metadata (`name`, `material`, `ai state`).
Here, `hot_update` means "update only the physics fields in place" and never read the metadata, while `snapshot_build` means "assemble a fresh output record with all fields for each active entity". It is a whole-record copy workload, not a runtime snapshot.
In this run, `AoS` wins the hot update at `10` and `100` entities, `SoA` takes over from `1 000` upward, and whole-record snapshot building stays close with `AoS` pulling ahead again at `10 000+`.

![aos soa graph](assets/BenchmarkAoSVsSoA.png)

[Benchmark results](assets/BenchmarkAoSVsSoA.txt)
