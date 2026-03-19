## Range over values vs index access

How much does `for _, v := range values` cost when each iteration copies a whole record, and when is it worth switching to index access or even `[]*T`?

> [!TIP]
> for hot loops over `[]T`, prefer `for i := range values` once records are more than tiny
> indexed iteration wins `26/30` cells in this benchmark
> here, `range` over values is about `2.3x` slower at `16B`, about `5.4x` slower at `128B`, and about `11.8x` slower at `512B` on the `100`-record row
> `[]*T` only helps in a few very large `256-512B` cases and is usually still slower than indexing into `[]T`

![range copy graph](assets/BenchmarkRangeCopy.png)

[Detailed line view](assets/BenchmarkRangeCopy-detail.png)

This benchmark scans hot fields from `16B` to `512B` records across lengths `10` to `100 000`. Indexed iteration is the stable winner, while `range` over values pays a full-record copy each trip and widens dramatically with size. `[]*T` avoids the value copy but adds indirection, so it only takes a few largest-record cells rather than becoming a general replacement for indexed iteration.

[Benchmark results](assets/BenchmarkRangeCopy.txt)
