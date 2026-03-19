## Owned return values: `T` vs `*T`

Should a function that builds and returns a fresh owned result use `T` or `*T`?

This section is about owned return values and the escape/allocation behavior of returning a freshly built result. It is not a blanket rule for APIs that need shared mutable identity, optional values, or polymorphic nil signaling.

> [!TIP]
> use `T` through at least `512B` for freshly built owned results in this benchmark  
> no cutoff appeared before `512B`  
> `*T` loses here because it allocates (`256 allocs/op` vs `0`)  
> shared mutable identity and optional/nil results are separate API-design concerns  

![return value vs pointer graph](assets/BenchmarkReturnValueVsPointer.png)

This benchmark builds `256` fresh results per operation from scalar seed inputs. `return T` wins at every tested size from `8B` to `512B`; `return *T` never catches up because it escapes and allocates (`256 allocs/op`, `2KiB/op` to `128KiB/op`) while `return T` stays at `0 allocs/op`.

[Benchmark results](assets/BenchmarkReturnValueVsPointer.txt)

Further reading:
- [There is no pass-by-reference in Go](https://dave.cheney.net/2017/04/29/there-is-no-pass-by-reference-in-go)
