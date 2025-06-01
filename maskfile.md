# Commands

Use [mask](https://github.com/jacobdeichert/mask) to run.

## bench

Runs all benchmarks and regenerates graphs.

```bash
$MASK bench deduplication
$MASK bench needleInHaystack
$MASK bench subset
$MASK bench append
```

###  deduplication

```bash
$MASK bench_one "BenchmarkDeduplication" "slice"
$MASK bench_one "BenchmarkDeduplication" "slice_sort_inplace"
$MASK bench_one "BenchmarkDeduplication" "map"
$MASK benchstat "BenchmarkDeduplication"
$MASK graph "BenchmarkDeduplication"
```

###  needleInHaystack

```bash
$MASK bench_one "BenchmarkNeedleInAHaystack" "slice"
$MASK bench_one "BenchmarkNeedleInAHaystack" "map"
$MASK benchstat "BenchmarkNeedleInAHaystack"
$MASK graph "BenchmarkNeedleInAHaystack"
```

###  subset

```bash
$MASK bench_one "BenchmarkSubset" "slice"
$MASK bench_one "BenchmarkSubset" "slice_sort_binsearch"
$MASK bench_one "BenchmarkSubset" "map"
$MASK benchstat "BenchmarkSubset"
$MASK graph "BenchmarkSubset"
```

###  append

```bash
$MASK bench_one "BenchmarkAppend" "expand"
$MASK bench_one "BenchmarkAppend" "for"
$MASK bench_one "BenchmarkAppend" "for_prealloc"
$MASK bench_one "BenchmarkAppend" "for_index"
$MASK benchstat "BenchmarkAppend"
$MASK graph "BenchmarkAppend"
```

## bench_one (benchname) (variant)

```bash
echo "Running: $benchname $variant"
go test -bench "^$benchname\$" -timeout 30m -count 10 -benchmem ./benchmarks/... -args -variant "$variant" > "assets/$benchname-$variant.txt"
```

## benchstat (benchname)

```bash
benchstat "assets/$benchname"*
benchstat -format csv "assets/$benchname"* > "assets/$benchname.csv"
```

## graph (benchname)

```bash
python3.11 scripts/plot.py "assets/$benchname.csv" "assets/$benchname.png"
echo "Wrote: assets/$benchname.png"
```
