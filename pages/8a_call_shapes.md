## Indirection cost: concrete `*T` vs interface vs generic

When the underlying object is the same concrete struct, how much overhead comes from indirection itself on repeated writes?

This section fixes one concrete type (`*scoringRequest32`), keeps the same `256` prebuilt records for every variant, and varies only how many repeated calls you make through each call shape. It shows two write-side views of the same problem: a tiny-body benchmark to expose the abstraction tax, then a more realistic write body to show how quickly that tax gets diluted by real work.

> [!TIP]
> use direct concrete `*T` as the baseline on hot mutable paths  
> no stable iteration-count cutoff appeared before `1000` repeated passes  
> with a tiny write body, both interface forms land around `15-18%` behind concrete here, and constrained generics around `65%`  
> with a more realistic write body, those same differences shrink into the low single digits  

### Read-write fixed `*scoringRequest32`, tiny body

![indirection tiny write graph](assets/BenchmarkIndirectionCostTinyWrite.png)

This version only increments `AccountID` and returns it. That makes the abstraction tax obvious: both interface forms land roughly `15-18%` behind concrete on geomean, and constrained generic about `65%`. `generic_exact` was the outlier here, so use this chart mainly to see how visible interface and constrained-generic overhead becomes when the method body is almost empty.

[Benchmark results](assets/BenchmarkIndirectionCostTinyWrite.txt)

### Read-write fixed `*scoringRequest32`, realistic body

![indirection write graph](assets/BenchmarkIndirectionCostWrite.png)

This version keeps the fuller write body over the same `256` `*scoringRequest32` records. Here the same call-shape differences are much smaller: preboxed interface dispatch is about `2.3%` behind concrete on geomean, exact generic about `2.6%`, `interface_box_each_call` about `3.4%`, and constrained generic about `5.8%`. For non-trivial work, indirection is usually a small constant tax rather than a late-breakpoint effect, and all variants stay at `0 allocs/op`. Treat the earlier tiny-body `generic_exact` anomaly as compiler-sensitive, not as a rule to generalize.

[Benchmark results](assets/BenchmarkIndirectionCostWrite.txt)

Supplementary read-only comparisons:
[BenchmarkIndirectionCostTinyRead results](assets/BenchmarkIndirectionCostTinyRead.txt)
[BenchmarkIndirectionCostRead results](assets/BenchmarkIndirectionCostRead.txt)

Supplementary mixed size sweep:
[BenchmarkCallShapes results](assets/BenchmarkCallShapes.txt)
