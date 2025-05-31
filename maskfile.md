# Commands

Use [mask](https://github.com/jacobdeichert/mask) to run.

## bench

Runs all benchmarks and regenerates graphs.

```bash
$MASK bench_one "BenchmarkDeduplication" "slice"
$MASK bench_one "BenchmarkDeduplication" "slice_sort_inplace"
$MASK bench_one "BenchmarkDeduplication" "map"
$MASK benchstat "BenchmarkDeduplication"
$MASK graph "BenchmarkDeduplication"
```

## bench_one (benchname) (variant)

```bash
echo "Running: $benchname $variant"
go test -bench "^$benchname\$" -timeout 30m -count 10 -benchmem ./... -args -variant "$variant" > "$benchname-$variant.txt"
```

## benchstat (benchname)

```bash
benchstat "$benchname"*
benchstat -format csv "$benchname"* > "$benchname.csv"
```

## graph (benchname)

```bash
python3.11 plot.py "$benchname.csv" "$benchname.png"
echo "Wrote: $benchname.png"
```
