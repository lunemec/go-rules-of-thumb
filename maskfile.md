# Commands

Use [mask](https://github.com/jacobdeichert/mask) to run.

## readme

Creates readme from templates with all the benchmark data refreshed.

```bash
cat pages/* > README.md
```

## test

Runs test cases.

```bash
go test -v -timeout 60m -count 1 -race  ./benchmarks/...
```

## bench

Runs all benchmarks and regenerates graphs.

```bash
$MASK bench deduplication
$MASK bench needleInHaystack
$MASK bench subset
$MASK bench append
$MASK bench assert
$MASK bench ifswitch
$MASK bench iterate
$MASK bench concatenate
$MASK bench aossoa
```

### deduplication

```bash
$MASK bench_one "BenchmarkDeduplication" "slice"
$MASK bench_one "BenchmarkDeduplication" "slice_sort_inplace"
$MASK bench_one "BenchmarkDeduplication" "map"
$MASK benchstat "BenchmarkDeduplication"
$MASK graph "BenchmarkDeduplication"
```

### needleInHaystack

```bash
$MASK bench_one "BenchmarkNeedleInAHaystack" "slice"
$MASK bench_one "BenchmarkNeedleInAHaystack" "map"
$MASK benchstat "BenchmarkNeedleInAHaystack"
$MASK graph "BenchmarkNeedleInAHaystack"
```

### subset

```bash
$MASK bench_one "BenchmarkSubset" "slice"
$MASK bench_one "BenchmarkSubset" "slice_sort_binsearch"
$MASK bench_one "BenchmarkSubset" "map"
$MASK benchstat "BenchmarkSubset"
$MASK graph "BenchmarkSubset"
```

### append

```bash
$MASK bench_one "BenchmarkAppend" "expand"
$MASK bench_one "BenchmarkAppend" "for"
$MASK bench_one "BenchmarkAppend" "for_prealloc"
$MASK bench_one "BenchmarkAppend" "for_index"
$MASK benchstat "BenchmarkAppend"
$MASK graph "BenchmarkAppend"
```

### assert

```bash
$MASK bench_one "BenchmarkAssert" "no_assert"
$MASK bench_one "BenchmarkAssert" "assert"
$MASK bench_one "BenchmarkAssert" "defer_assert"
$MASK benchstat "BenchmarkAssert"
$MASK graph "BenchmarkAssert"
```

### ifswitch

```bash
$MASK bench_one "BenchmarkIfSwitch" "if"
$MASK bench_one "BenchmarkIfSwitch" "switch"
$MASK bench_one "BenchmarkIfSwitch" "if_5"
$MASK bench_one "BenchmarkIfSwitch" "switch_5"
$MASK benchstat "BenchmarkIfSwitch"
$MASK graph "BenchmarkIfSwitch"
```

### iterate

```bash
$MASK bench_one "BenchmarkIterate" "slice_iterate"
$MASK bench_one "BenchmarkIterate" "range_func"
$MASK bench_one "BenchmarkIterate" "direct"
$MASK benchstat "BenchmarkIterate"
$MASK graph "BenchmarkIterate"
```

### concatenate

```bash
$MASK bench_one "BenchmarkConcat" "plus"
$MASK bench_one "BenchmarkConcat" "sprintf"
$MASK bench_one "BenchmarkConcat" "join"
$MASK bench_one "BenchmarkConcat" "builder"
$MASK bench_one "BenchmarkConcat" "builder_pool"
$MASK benchstat "BenchmarkConcat"
$MASK graph "BenchmarkConcat"
```

### aossoa

```bash
$MASK bench_one "BenchmarkAoSVsSoA" "aos_hot_update"
$MASK bench_one "BenchmarkAoSVsSoA" "soa_hot_update"
$MASK bench_one "BenchmarkAoSVsSoA" "aos_snapshot"
$MASK bench_one "BenchmarkAoSVsSoA" "soa_snapshot"
$MASK benchstat "BenchmarkAoSVsSoA"
$MASK graph "BenchmarkAoSVsSoA"
```

### pointer_copy

```bash
$MASK bench_one "BenchmarkCopyVsPointer" "copy"
$MASK bench_one "BenchmarkCopyVsPointer" "pointer"
$MASK benchstat "BenchmarkCopyVsPointer"
```

## bench_one (benchname) (variant)

```bash
echo "Running: $benchname $variant"
go test -bench "^$benchname\$" -timeout 60m -count 10 -benchmem ./benchmarks/... -args -variant "$variant" > "assets/$benchname-$variant.txt"
```

## benchstat (benchname)

```bash
cd assets
benchstat "$benchname"-*.txt > "$benchname.txt"
benchstat -format csv "$benchname"-*.txt > "$benchname.csv"
```

## graph (benchname)

```bash
python3.11 scripts/plot.py "assets/$benchname.csv" "assets/$benchname.png"
echo "Wrote: assets/$benchname.png"
```
