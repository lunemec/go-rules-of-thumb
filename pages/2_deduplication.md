## Deduplication

When is it more efficient to deduplicate a `slice` as opposed to using a `map[]struct{}` for the same purpose?

> **TL;DR**: use `map` when `len(haystack) > 100`. If you must reduce allocations, use sort + in-place slice
> deduplication. Suprisingly it is fast enough.

![deduplication graph](assets/BenchmarkDeduplication.png "Deduplication variants performance")

```
                              │ BenchmarkDeduplication-map.txt │ BenchmarkDeduplication-slice_sort_inplace.txt │       BenchmarkDeduplication-slice.txt       │
                              │             sec/op             │        sec/op         vs base                 │     sec/op      vs base                      │
Deduplication/size=10-8                          205.55n ±  1%           55.48n ±  2%   -73.01% (p=0.000 n=10)     114.40n ± 2%      -44.34% (p=0.000 n=10)
Deduplication/size=100-8                         1563.0n ±  1%           808.6n ±  1%   -48.26% (p=0.000 n=10)     1262.5n ± 1%      -19.23% (p=0.000 n=10)
Deduplication/size=500-8                          9.117µ ±  1%           4.903µ ±  4%   -46.22% (p=0.000 n=10)     31.617µ ± 8%     +246.81% (p=0.000 n=10)
Deduplication/size=1000-8                         19.56µ ±  0%           10.91µ ±  4%   -44.23% (p=0.000 n=10)     108.98µ ± 2%     +457.19% (p=0.000 n=10)
Deduplication/size=5000-8                         101.6µ ±  1%           187.1µ ±  2%   +84.21% (p=0.000 n=10)     2432.6µ ± 3%    +2294.79% (p=0.000 n=10)
Deduplication/size=10000-8                        220.5µ ±  5%           449.5µ ±  1%  +103.87% (p=0.000 n=10)     9744.1µ ± 2%    +4319.66% (p=0.000 n=10)
Deduplication/size=50000-8                        1.136m ±  8%           2.782m ±  1%  +144.80% (p=0.000 n=10)    235.701m ± 3%   +20640.07% (p=0.000 n=10)
Deduplication/size=100000-8                       2.208m ± 27%           5.999m ±  1%  +171.71% (p=0.000 n=10)    927.130m ± 1%   +41889.99% (p=0.000 n=10)
Deduplication/size=500000-8                       25.18m ± 19%           34.65m ±  4%   +37.58% (p=0.000 n=10)   23061.71m ± 3%   +91471.94% (p=0.000 n=10)
Deduplication/size=1000000-8                      60.85m ± 12%           73.29m ± 19%   +20.46% (p=0.001 n=10)   92737.14m ± 2%  +152307.91% (p=0.000 n=10)
Deduplication/size=5000000-8                      393.3m ±  8%           387.9m ±  1%         ~ (p=0.123 n=10)
```

```
                              │ BenchmarkDeduplication-map.txt │ BenchmarkDeduplication-slice_sort_inplace.txt │    BenchmarkDeduplication-slice.txt     │
                              │              B/op              │         B/op           vs base                │     B/op       vs base                  │
Deduplication/size=10-8                            488.00 ± 0%              80.00 ± 0%  -83.61% (p=0.000 n=10)     200.00 ± 0%  -59.02% (p=0.000 n=10)
Deduplication/size=100-8                           4136.0 ± 0%              896.0 ± 0%  -78.34% (p=0.000 n=10)     1912.0 ± 0%  -53.77% (p=0.000 n=10)
Deduplication/size=500-8                         26.039Ki ± 0%            4.000Ki ± 0%  -84.64% (p=0.000 n=10)   11.992Ki ± 0%  -53.95% (p=0.000 n=10)
Deduplication/size=1000-8                        52.078Ki ± 0%            8.000Ki ± 0%  -84.64% (p=0.000 n=10)   22.617Ki ± 0%  -56.57% (p=0.000 n=10)
Deduplication/size=5000-8                        224.31Ki ± 0%            40.00Ki ± 0%  -82.17% (p=0.000 n=10)   125.24Ki ± 0%  -44.17% (p=0.000 n=10)
Deduplication/size=10000-8                       448.63Ki ± 0%            80.00Ki ± 0%  -82.17% (p=0.000 n=10)   261.24Ki ± 0%  -41.77% (p=0.000 n=10)
Deduplication/size=50000-8                       1938.5Ki ± 0%            392.0Ki ± 0%  -79.78% (p=0.000 n=10)   1525.2Ki ± 0%  -21.32% (p=0.000 n=10)
Deduplication/size=100000-8                      3877.1Ki ± 0%            784.0Ki ± 0%  -79.78% (p=0.000 n=10)   3237.3Ki ± 0%  -16.50% (p=0.000 n=10)
Deduplication/size=500000-8                      25.681Mi ± 0%            3.820Mi ± 0%  -85.12% (p=0.000 n=10)   16.521Mi ± 0%  -35.67% (p=0.000 n=10)
Deduplication/size=1000000-8                     51.346Mi ± 0%            7.633Mi ± 0%  -85.13% (p=0.000 n=10)   32.888Mi ± 0%  -35.95% (p=0.000 n=10)
Deduplication/size=5000000-8                     220.61Mi ± 0%            38.15Mi ± 0%  -82.71% (p=0.000 n=10)
Deduplication/size=10000000-8                    441.22Mi ± 0%            76.30Mi ± 0%  -82.71% (p=0.000 n=10)
```

```
                              │ BenchmarkDeduplication-map.txt │ BenchmarkDeduplication-slice_sort_inplace.txt │    BenchmarkDeduplication-slice.txt    │
                              │           allocs/op            │      allocs/op        vs base                 │  allocs/op   vs base                   │
Deduplication/size=10-8                             5.000 ± 0%             1.000 ± 0%   -80.00% (p=0.000 n=10)    5.000 ± 0%         ~ (p=1.000 n=10) ¹
Deduplication/size=100-8                            5.000 ± 0%             1.000 ± 0%   -80.00% (p=0.000 n=10)    8.000 ± 0%   +60.00% (p=0.000 n=10)
Deduplication/size=500-8                            5.000 ± 0%             1.000 ± 0%   -80.00% (p=0.000 n=10)   11.000 ± 0%  +120.00% (p=0.000 n=10)
Deduplication/size=1000-8                           7.000 ± 0%             1.000 ± 0%   -85.71% (p=0.000 n=10)   12.000 ± 0%   +71.43% (p=0.000 n=10)
Deduplication/size=5000-8                          19.000 ± 0%             1.000 ± 0%   -94.74% (p=0.000 n=10)   16.000 ± 0%   -15.79% (p=0.000 n=10)
Deduplication/size=10000-8                         35.000 ± 0%             1.000 ± 0%   -97.14% (p=0.000 n=10)   18.000 ± 0%   -48.57% (p=0.000 n=10)
Deduplication/size=50000-8                        131.000 ± 0%             1.000 ± 0%   -99.24% (p=0.000 n=10)   24.000 ± 0%   -81.68% (p=0.000 n=10)
Deduplication/size=100000-8                       259.000 ± 0%             1.000 ± 0%   -99.61% (p=0.000 n=10)   27.000 ± 0%   -89.58% (p=0.000 n=10)
Deduplication/size=500000-8                      2051.000 ± 0%             1.000 ± 0%   -99.95% (p=0.000 n=10)   34.000 ± 3%   -98.34% (p=0.000 n=10)
Deduplication/size=1000000-8                     4099.000 ± 0%             1.000 ± 0%   -99.98% (p=0.000 n=10)   37.000 ± 3%   -99.10% (p=0.000 n=10)
Deduplication/size=5000000-8                    16387.000 ± 0%             1.000 ± 0%   -99.99% (p=0.000 n=10)
Deduplication/size=10000000-8                   32771.000 ± 0%             1.000 ± 0%  -100.00% (p=0.000 n=10)
```
