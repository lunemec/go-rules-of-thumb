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
$MASK bench valuesvspointers
$MASK bench hotcoldsplit
$MASK bench linearvsrandomaccess
$MASK bench callshapes
$MASK bench indirectioncostread
$MASK bench indirectioncostwrite
$MASK bench indirectioncosttinyread
$MASK bench indirectioncosttinywrite
$MASK bench paramvaluevspointer
$MASK bench returnvaluevspointer
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

### concatenate-large

```bash
$MASK bench_one "BenchmarkConcatLarge" "plus"
$MASK bench_one "BenchmarkConcatLarge" "sprintf"
$MASK bench_one "BenchmarkConcatLarge" "join"
$MASK bench_one "BenchmarkConcatLarge" "builder"
$MASK bench_one "BenchmarkConcatLarge" "builder_pool"
$MASK benchstat "BenchmarkConcatLarge"
$MASK graph "BenchmarkConcatLarge"
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

### valuesvspointers

```bash
$MASK bench_one "BenchmarkValuesVsPointers" "values_hot_scan"
$MASK bench_one "BenchmarkValuesVsPointers" "pointers_hot_scan"
$MASK bench_one "BenchmarkValuesVsPointers" "values_snapshot"
$MASK bench_one "BenchmarkValuesVsPointers" "pointers_snapshot"
$MASK benchstat "BenchmarkValuesVsPointers"
$MASK graph "BenchmarkValuesVsPointers"
```

### hotcoldsplit

```bash
$MASK bench_one "BenchmarkHotColdSplit" "inline_hot_scan"
$MASK bench_one "BenchmarkHotColdSplit" "split_hot_scan"
$MASK bench_one "BenchmarkHotColdSplit" "inline_snapshot"
$MASK bench_one "BenchmarkHotColdSplit" "split_snapshot"
$MASK benchstat "BenchmarkHotColdSplit"
$MASK graph "BenchmarkHotColdSplit"
```

### linearvsrandomaccess

```bash
$MASK bench_one "BenchmarkLinearVsRandomAccess" "linear"
$MASK bench_one "BenchmarkLinearVsRandomAccess" "random"
$MASK benchstat "BenchmarkLinearVsRandomAccess"
$MASK graph "BenchmarkLinearVsRandomAccess"
```

### callshapes

```bash
$MASK bench_one "BenchmarkCallShapes" "concrete_ptr"
$MASK bench_one "BenchmarkCallShapes" "interface_dispatch"
$MASK bench_one "BenchmarkCallShapes" "interface_box_each_call"
$MASK bench_one "BenchmarkCallShapes" "generic_exact"
$MASK bench_one "BenchmarkCallShapes" "generic_constraint"
$MASK benchstat "BenchmarkCallShapes"
$MASK graph "BenchmarkCallShapes"
```

### indirectioncostread

```bash
echo "Running: BenchmarkIndirectionCostRead concrete_ptr"
go test -run '^$' -bench "^BenchmarkIndirectionCostRead\$" -benchtime 500ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "concrete_ptr" > "assets/BenchmarkIndirectionCostRead-concrete_ptr.txt"
echo "Running: BenchmarkIndirectionCostRead interface_dispatch"
go test -run '^$' -bench "^BenchmarkIndirectionCostRead\$" -benchtime 500ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "interface_dispatch" > "assets/BenchmarkIndirectionCostRead-interface_dispatch.txt"
echo "Running: BenchmarkIndirectionCostRead interface_box_each_call"
go test -run '^$' -bench "^BenchmarkIndirectionCostRead\$" -benchtime 500ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "interface_box_each_call" > "assets/BenchmarkIndirectionCostRead-interface_box_each_call.txt"
echo "Running: BenchmarkIndirectionCostRead generic_exact"
go test -run '^$' -bench "^BenchmarkIndirectionCostRead\$" -benchtime 500ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "generic_exact" > "assets/BenchmarkIndirectionCostRead-generic_exact.txt"
echo "Running: BenchmarkIndirectionCostRead generic_constraint"
go test -run '^$' -bench "^BenchmarkIndirectionCostRead\$" -benchtime 500ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "generic_constraint" > "assets/BenchmarkIndirectionCostRead-generic_constraint.txt"
$MASK benchstat "BenchmarkIndirectionCostRead"
$MASK graph "BenchmarkIndirectionCostRead"
```

### indirectioncostwrite

```bash
echo "Running: BenchmarkIndirectionCostWrite concrete_ptr"
go test -run '^$' -bench "^BenchmarkIndirectionCostWrite\$" -benchtime 500ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "concrete_ptr" > "assets/BenchmarkIndirectionCostWrite-concrete_ptr.txt"
echo "Running: BenchmarkIndirectionCostWrite interface_dispatch"
go test -run '^$' -bench "^BenchmarkIndirectionCostWrite\$" -benchtime 500ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "interface_dispatch" > "assets/BenchmarkIndirectionCostWrite-interface_dispatch.txt"
echo "Running: BenchmarkIndirectionCostWrite interface_box_each_call"
go test -run '^$' -bench "^BenchmarkIndirectionCostWrite\$" -benchtime 500ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "interface_box_each_call" > "assets/BenchmarkIndirectionCostWrite-interface_box_each_call.txt"
echo "Running: BenchmarkIndirectionCostWrite generic_exact"
go test -run '^$' -bench "^BenchmarkIndirectionCostWrite\$" -benchtime 500ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "generic_exact" > "assets/BenchmarkIndirectionCostWrite-generic_exact.txt"
echo "Running: BenchmarkIndirectionCostWrite generic_constraint"
go test -run '^$' -bench "^BenchmarkIndirectionCostWrite\$" -benchtime 500ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "generic_constraint" > "assets/BenchmarkIndirectionCostWrite-generic_constraint.txt"
$MASK benchstat "BenchmarkIndirectionCostWrite"
$MASK graph "BenchmarkIndirectionCostWrite"
```

### indirectioncosttinyread

```bash
echo "Running: BenchmarkIndirectionCostTinyRead concrete_ptr"
go test -run '^$' -bench "^BenchmarkIndirectionCostTinyRead\$" -benchtime 500ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "concrete_ptr" > "assets/BenchmarkIndirectionCostTinyRead-concrete_ptr.txt"
echo "Running: BenchmarkIndirectionCostTinyRead interface_dispatch"
go test -run '^$' -bench "^BenchmarkIndirectionCostTinyRead\$" -benchtime 500ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "interface_dispatch" > "assets/BenchmarkIndirectionCostTinyRead-interface_dispatch.txt"
echo "Running: BenchmarkIndirectionCostTinyRead interface_box_each_call"
go test -run '^$' -bench "^BenchmarkIndirectionCostTinyRead\$" -benchtime 500ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "interface_box_each_call" > "assets/BenchmarkIndirectionCostTinyRead-interface_box_each_call.txt"
echo "Running: BenchmarkIndirectionCostTinyRead generic_exact"
go test -run '^$' -bench "^BenchmarkIndirectionCostTinyRead\$" -benchtime 500ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "generic_exact" > "assets/BenchmarkIndirectionCostTinyRead-generic_exact.txt"
echo "Running: BenchmarkIndirectionCostTinyRead generic_constraint"
go test -run '^$' -bench "^BenchmarkIndirectionCostTinyRead\$" -benchtime 500ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "generic_constraint" > "assets/BenchmarkIndirectionCostTinyRead-generic_constraint.txt"
$MASK benchstat "BenchmarkIndirectionCostTinyRead"
$MASK graph "BenchmarkIndirectionCostTinyRead"
```

### indirectioncosttinywrite

```bash
echo "Running: BenchmarkIndirectionCostTinyWrite concrete_ptr"
go test -run '^$' -bench "^BenchmarkIndirectionCostTinyWrite\$" -benchtime 500ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "concrete_ptr" > "assets/BenchmarkIndirectionCostTinyWrite-concrete_ptr.txt"
echo "Running: BenchmarkIndirectionCostTinyWrite interface_dispatch"
go test -run '^$' -bench "^BenchmarkIndirectionCostTinyWrite\$" -benchtime 500ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "interface_dispatch" > "assets/BenchmarkIndirectionCostTinyWrite-interface_dispatch.txt"
echo "Running: BenchmarkIndirectionCostTinyWrite interface_box_each_call"
go test -run '^$' -bench "^BenchmarkIndirectionCostTinyWrite\$" -benchtime 500ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "interface_box_each_call" > "assets/BenchmarkIndirectionCostTinyWrite-interface_box_each_call.txt"
echo "Running: BenchmarkIndirectionCostTinyWrite generic_exact"
go test -run '^$' -bench "^BenchmarkIndirectionCostTinyWrite\$" -benchtime 500ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "generic_exact" > "assets/BenchmarkIndirectionCostTinyWrite-generic_exact.txt"
echo "Running: BenchmarkIndirectionCostTinyWrite generic_constraint"
go test -run '^$' -bench "^BenchmarkIndirectionCostTinyWrite\$" -benchtime 500ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "generic_constraint" > "assets/BenchmarkIndirectionCostTinyWrite-generic_constraint.txt"
$MASK benchstat "BenchmarkIndirectionCostTinyWrite"
$MASK graph "BenchmarkIndirectionCostTinyWrite"
```

### paramvaluevspointer

```bash
$MASK bench_one "BenchmarkParamValueVsPointer" "value_param"
$MASK bench_one "BenchmarkParamValueVsPointer" "pointer_param"
$MASK benchstat "BenchmarkParamValueVsPointer"
$MASK graph "BenchmarkParamValueVsPointer"
```

### returnvaluevspointer

```bash
$MASK bench_one "BenchmarkReturnValueVsPointer" "return_value"
$MASK bench_one "BenchmarkReturnValueVsPointer" "return_pointer"
$MASK benchstat "BenchmarkReturnValueVsPointer"
$MASK graph "BenchmarkReturnValueVsPointer"
```

## bench_one (benchname) (variant)

```bash
echo "Running: $benchname $variant"
go test -run '^$' -bench "^$benchname\$" -benchtime 100ms -timeout 60m -count 6 -benchmem ./benchmarks/... -args -variant "$variant" > "assets/$benchname-$variant.txt"
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
