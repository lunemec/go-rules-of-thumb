## Strings concatenation

Is it more efficient to `"str1" + var`, `fmt.Sprintf()`, `strings.Join()` or `strings.Builder`? When does it make sense to add `sync.Pool`?

> **TL;DR**: use `strings.Builder` when `len(str) < 100 & N ops < 1000`, use `sync.Pool + strings.Builder` when doing this for every request. For `len(str) > 100` use `+` or `strings.Join`.
>
> Use `fmt.Sprintf` for regular string formatting (not just concatenation).

| Type                 | len(str) | N ops | ns/op         |     |
| -------------------- | -------- | ----- | ------------- | --- |
| plus_sign            | 10       | 10    | 377.2 ns/op   |     |
| sprintf              | 10       | 10    | 1174 ns/op    |
| strings_join         | 10       | 10    | 457.2 ns/op   |     |
| strings_builder      | 10       | 10    | 226.6 ns/op   | ✅  |
| strings_builder_pool | 10       | 10    | 242.4 ns/op   |     |
| plus_sign            | 10       | 100   | 4976 ns/op    |
| sprintf              | 10       | 100   | 13242 ns/op   |
| strings_join         | 10       | 100   | 5832 ns/op    |
| strings_builder      | 10       | 100   | 1275 ns/op    |
| strings_builder_pool | 10       | 100   | 1265 ns/op    | ✅  |
| plus_sign            | 10       | 500   | 62458 ns/op   |
| sprintf              | 10       | 500   | 106616 ns/op  |
| strings_join         | 10       | 500   | 67195 ns/op   |
| strings_builder      | 10       | 500   | 6515 ns/op    | ✅  |
| strings_builder_pool | 10       | 500   | 6670 ns/op    |
| plus_sign            | 10       | 1000  | 209530 ns/op  |
| sprintf              | 10       | 1000  | 308757 ns/op  |
| strings_join         | 10       | 1000  | 219302 ns/op  |
| strings_builder      | 10       | 1000  | 13754 ns/op   |
| strings_builder_pool | 10       | 1000  | 13660 ns/op   | ✅  |
| plus_sign            | 100      | 10    | 533.6 ns/op   |     |
| sprintf              | 100      | 10    | 1342 ns/op    |
| strings_join         | 100      | 10    | 586.6 ns/op   |     |
| strings_builder      | 100      | 10    | 975.8 ns/op   |     |
| strings_builder_pool | 100      | 10    | 1028 ns/op    |
| plus_sign            | 100      | 100   | 6670 ns/op    | ✅  |
| sprintf              | 100      | 100   | 14949 ns/op   |
| strings_join         | 100      | 100   | 7562 ns/op    |
| strings_builder      | 100      | 100   | 9713 ns/op    |
| strings_builder_pool | 100      | 100   | 9918 ns/op    |
| plus_sign            | 100      | 500   | 71459 ns/op   |
| sprintf              | 100      | 500   | 116144 ns/op  |
| strings_join         | 100      | 500   | 75915 ns/op   |
| strings_builder      | 100      | 500   | 41344 ns/op   | ✅  |
| strings_builder_pool | 100      | 500   | 43655 ns/op   |
| plus_sign            | 100      | 1000  | 227323 ns/op  |
| sprintf              | 100      | 1000  | 519672 ns/op  |
| strings_join         | 100      | 1000  | 316674 ns/op  |
| strings_builder      | 100      | 1000  | 93747 ns/op   | ✅  |
| strings_builder_pool | 100      | 1000  | 99357 ns/op   |
| plus_sign            | 500      | 10    | 1900 ns/op    | ✅  |
| sprintf              | 500      | 10    | 2741 ns/op    |
| strings_join         | 500      | 10    | 2066 ns/op    |
| strings_builder      | 500      | 10    | 7532 ns/op    |
| strings_builder_pool | 500      | 10    | 6380 ns/op    |
| plus_sign            | 500      | 100   | 20481 ns/op   | ✅  |
| sprintf              | 500      | 100   | 30771 ns/op   |
| strings_join         | 500      | 100   | 21401 ns/op   |
| strings_builder      | 500      | 100   | 68556 ns/op   |
| strings_builder_pool | 500      | 100   | 69828 ns/op   |
| plus_sign            | 500      | 500   | 139848 ns/op  | ✅  |
| sprintf              | 500      | 500   | 193818 ns/op  |
| strings_join         | 500      | 500   | 143040 ns/op  |
| strings_builder      | 500      | 500   | 313665 ns/op  |
| strings_builder_pool | 500      | 500   | 322586 ns/op  |
| plus_sign            | 500      | 1000  | 370027 ns/op  | ✅  |
| sprintf              | 500      | 1000  | 515227 ns/op  |
| strings_join         | 500      | 1000  | 379155 ns/op  |
| strings_builder      | 500      | 1000  | 647395 ns/op  |
| strings_builder_pool | 500      | 1000  | 573904 ns/op  |
| plus_sign            | 1000     | 10    | 3220 ns/op    | ✅  |
| sprintf              | 1000     | 10    | 4613 ns/op    |
| strings_join         | 1000     | 10    | 3342 ns/op    |
| strings_builder      | 1000     | 10    | 13307 ns/op   |
| strings_builder_pool | 1000     | 10    | 13747 ns/op   |
| plus_sign            | 1000     | 100   | 40945 ns/op   |
| sprintf              | 1000     | 100   | 49810 ns/op   |
| strings_join         | 1000     | 100   | 36396 ns/op   | ✅  |
| strings_builder      | 1000     | 100   | 142254 ns/op  |
| strings_builder_pool | 1000     | 100   | 149184 ns/op  |
| plus_sign            | 1000     | 500   | 224290 ns/op  | ✅  |
| sprintf              | 1000     | 500   | 296963 ns/op  |
| strings_join         | 1000     | 500   | 233688 ns/op  |
| strings_builder      | 1000     | 500   | 783015 ns/op  |
| strings_builder_pool | 1000     | 500   | 657683 ns/op  |
| plus_sign            | 1000     | 1000  | 557672 ns/op  | ✅  |
| sprintf              | 1000     | 1000  | 715151 ns/op  |
| strings_join         | 1000     | 1000  | 571112 ns/op  |
| strings_builder      | 1000     | 1000  | 1326209 ns/op |
| strings_builder_pool | 1000     | 1000  | 1106394 ns/op |
