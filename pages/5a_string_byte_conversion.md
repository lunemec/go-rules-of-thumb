## Read-only conversions: `string` vs `[]byte`

How expensive is repeatedly crossing the `string` / `[]byte` boundary when the hot path only needs to scan read-only payloads?

This benchmark measures repeated read-only scans over ASCII payloads. It does not cover mutating the converted slice, retaining converted values, or escape-heavy API boundaries.

> [!TIP]
> keep data in the domain you already have on hot paths
> avoid repeated `string(b)` conversions in loops; here they are about `26-55%` slower and, from `64B` upward, allocate once per conversion
> on this Go `1.26.1` toolchain, read-only `[]byte(s)` stays close to the direct paths and shows `0 allocs/op` here
> if you need a conversion, do it once outside the loop

![string byte conversion graph](assets/BenchmarkStringByteConversion.png)

[Detailed line view](assets/BenchmarkStringByteConversion-detail.png)

Across `16B` to `64 KiB` payloads and `1` to `1000` repeated scans, `string(b)` is the consistently bad path: at `64B x 1000` it is about `55%` slower than direct bytes (`25.2µs` vs `16.2µs`), and at `64 KiB x 1000` it burns about `65 MiB/op` and about `1000 allocs/op`. The direct `string` path and the read-only `[]byte(s)` path stay roughly tied with the direct-byte baseline, so the practical rule from this benchmark is to worry about repeated `string(b)` first and treat `[]byte(s)` as toolchain-sensitive.

[Benchmark results](assets/BenchmarkStringByteConversion.txt)
