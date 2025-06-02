## Append

```go
append([]T, elems...) // append_expand
```

vs

```go
for _, e := range elems {
    arr = append(arr, e) // append_for
}
```

vs

```go
for _, e := range elems {
    arr = append(arr, e) // append_for_prealloc
}
```

vs

```go
for i, e := range elems {
    arr[i] = e // append_for_index (pre-allocated)
}
```

> **TL;DR**: ALWAYS use `append([]T, elems...)` because `for` looping may trigger multiple array re-sizings, whereas `append` will always allocate only once. If you must use `for` loop (extra logic), try to pre-allocate the slice.

Even though regular `append()` has time complexity `O(1)` (amortized constant-time), because every time it needs to allocate more space, it grows the underlying data array by 2x (until 512 elements, after 512 it grows less), simply by having to allocate + copy makes it significantly slower than if you are able to calculate the resulting size and pre-allocating.

| Type                | len(A) | len(B) | ns/op       | B/op       | allocs/op   |     |
| ------------------- | ------ | ------ | ----------- | ---------- | ----------- | --- |
| append_expand       | 10     | 1000   | 878.7 ns/op | 8192 B/op  | 1 allocs/op | ✅  |
| append_for_index    | 10     | 1000   | 1049 ns/op  | 8192 B/op  | 1 allocs/op |
| append_for_prealloc | 10     | 1000   | 1148 ns/op  | 8192 B/op  | 1 allocs/op |
| append_for          | 10     | 1000   | 2115 ns/op  | 19936 B/op | 7 allocs/op |
