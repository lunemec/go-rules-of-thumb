## Deduplication

When is it more efficient to deduplicate a `slice` as opposed to using a `map[]struct{}` for the same purpose?

```go
// slice
var result []int
for _, val := range haystack {
    if !slices.Contains(result, val) {
        result = append(result, val)
    }
}
return result
```

```go
// in-place sort + deduplicate
// "borrowed" from https://go.dev/wiki/SliceTricks#in-place-deduplicate-comparable, thanks!
// Note sort + slices.Compact is the same thing.
sort.Ints(haystack)
j := 0
for i := 1; i < len(haystack); i++ {
    if haystack[j] == haystack[i] {
        continue
    }
    j++
    // preserve the original data
    // in[i], in[j] = in[j], in[i]
    // only set what is required
    haystack[j] = haystack[i]
}
return haystack[:j+1]
```

```go
// map
result := make([]int, 0, len(haystack))
seen := make(map[int]struct{}, len(haystack))

for _, item := range haystack {
    if _, ok := seen[item]; ok {
        continue
    }

    seen[item] = struct{}{}
    result = append(result, item)
}

return result
```

> **TL;DR**: use `map` when `len(haystack) > 100`.  
> If you must reduce allocations, use in-place sort + deduplication. Suprisingly it is fast enough.

![deduplication graph](assets/BenchmarkDeduplication.png)

[Benchmark results](assets/BenchmarkDeduplication.txt)
