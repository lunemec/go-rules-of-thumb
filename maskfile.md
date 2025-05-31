# Commands

Use [mask](https://github.com/jacobdeichert/mask) to run.

## benchmark

```bash
go test -bench '^BenchmarkDeduplication$' -timeout 30m -count 10 -benchmem ./... -args -test slice > dedup_slice.txt
go test -bench '^BenchmarkDeduplication$' -count 10 -benchmem ./... -args -test slice_comparable >
dedup_slice_comparable.txt
go test -bench '^BenchmarkDeduplication$' -count 10 -benchmem ./... -args -test map > dedup_map.txt
```

## benchstat_csv

```bash
benchstat -format csv dedup_slice.txt dedup_map.txt dedup_slice_comparable.txt > benchstat.csv
```

## graph

```bash
python3.11 plot.py
```
