## Range over func

With [Go 1.23 came new feature - range over func](https://go.dev/blog/range-functions), lets check when it makes sense to use that over
pre-allocating a slice and putting values in it.

I'm quite suprised to see that range over func adds extra 3x number of allocations somewhere.
Not sure where, that is to be measured later.

Here are some results:

```
BenchmarkRangeFunc/slice(10)_iterations(10)-8            976581       1175 ns/op      880 B/op       11 allocs/op
BenchmarkRangeFunc/iter_func(10)_iterations(10)-8        770889       1556 ns/op      528 B/op       33 allocs/op
BenchmarkRangeFunc/slice(10)_iterations(100)-8           112356      10747 ns/op     8080 B/op      101 allocs/op
BenchmarkRangeFunc/iter_func(10)_iterations(100)-8        80978      14675 ns/op     4848 B/op      303 allocs/op
BenchmarkRangeFunc/slice(10)_iterations(500)-8            22530      54421 ns/op    40080 B/op      501 allocs/op
BenchmarkRangeFunc/iter_func(10)_iterations(500)-8        16900      70454 ns/op    24048 B/op     1503 allocs/op
BenchmarkRangeFunc/slice(10)_iterations(1000)-8           10000     105975 ns/op    80080 B/op     1001 allocs/op
BenchmarkRangeFunc/iter_func(10)_iterations(1000)-8        8494     140141 ns/op    48048 B/op     3003 allocs/op
BenchmarkRangeFunc/slice(100)_iterations(10)-8           123470       9899 ns/op     9856 B/op       11 allocs/op
BenchmarkRangeFunc/iter_func(100)_iterations(10)-8       116560      10279 ns/op      528 B/op       33 allocs/op
BenchmarkRangeFunc/slice(100)_iterations(100)-8           13488      88975 ns/op    90496 B/op      101 allocs/op
BenchmarkRangeFunc/iter_func(100)_iterations(100)-8       12710      94388 ns/op     4848 B/op      303 allocs/op
BenchmarkRangeFunc/slice(100)_iterations(500)-8            2715     443446 ns/op   448897 B/op      501 allocs/op
BenchmarkRangeFunc/iter_func(100)_iterations(500)-8        2498     492394 ns/op    24048 B/op     1503 allocs/op
BenchmarkRangeFunc/slice(100)_iterations(1000)-8           1345     934711 ns/op   896899 B/op     1001 allocs/op
BenchmarkRangeFunc/iter_func(100)_iterations(1000)-8       1252     975041 ns/op    48048 B/op     3003 allocs/op
BenchmarkRangeFunc/slice(500)_iterations(10)-8            25249      47466 ns/op    45056 B/op       11 allocs/op
BenchmarkRangeFunc/iter_func(500)_iterations(10)-8        24255      49390 ns/op      528 B/op       33 allocs/op
BenchmarkRangeFunc/slice(500)_iterations(100)-8            2761     434421 ns/op   413697 B/op      101 allocs/op
BenchmarkRangeFunc/iter_func(500)_iterations(100)-8        2652     455068 ns/op     4848 B/op      303 allocs/op
BenchmarkRangeFunc/slice(500)_iterations(500)-8             556    2159299 ns/op  2052103 B/op      501 allocs/op
BenchmarkRangeFunc/iter_func(500)_iterations(500)-8         532    2263910 ns/op    24048 B/op     1503 allocs/op
BenchmarkRangeFunc/slice(500)_iterations(1000)-8            277    4309678 ns/op  4100116 B/op     1001 allocs/op
BenchmarkRangeFunc/iter_func(500)_iterations(1000)-8        266    4501908 ns/op    48048 B/op     3003 allocs/op
BenchmarkRangeFunc/slice(1000)_iterations(10)-8           12734      94393 ns/op    90112 B/op       11 allocs/op
BenchmarkRangeFunc/iter_func(1000)_iterations(10)-8       12175      98567 ns/op      528 B/op       33 allocs/op
BenchmarkRangeFunc/slice(1000)_iterations(100)-8           1383     862572 ns/op   827395 B/op      101 allocs/op
BenchmarkRangeFunc/iter_func(1000)_iterations(100)-8       1328     907260 ns/op     4848 B/op      303 allocs/op
BenchmarkRangeFunc/slice(1000)_iterations(500)-8            277    4285715 ns/op  4104203 B/op      501 allocs/op
BenchmarkRangeFunc/iter_func(1000)_iterations(500)-8        267    4479069 ns/op    24048 B/op     1503 allocs/op
BenchmarkRangeFunc/slice(1000)_iterations(1000)-8           139    8939082 ns/op  8200234 B/op     1001 allocs/op
BenchmarkRangeFunc/iter_func(1000)_iterations(1000)-8       127    9272349 ns/op    48048 B/op     3003 allocs/op
```
