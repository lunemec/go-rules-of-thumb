# AGENTS.md

Purpose: working rules for agents contributing to `go-rules-of-thumb`.

## Project intent
- This repository is a benchmark-backed catalog of practical Go performance heuristics.
- The main output is a readable `README.md` that turns benchmark results into cautious "rules of thumb".
- The repo is about measured tradeoffs, not dogma. Keep the disclaimer and the "benchmark your own workload" framing intact.

## Repo map
- `README.md`: published document assembled from `pages/*` in lexical order; regenerate it via `mask readme`.
- `pages/*.md`: source fragments/stubs for the README; numeric prefixes control section order, and all documentation edits belong here.
- `benchmarks/*_test.go`: benchmark and correctness coverage for each topic.
- `assets/`: generated benchmark outputs, `benchstat` summaries, CSVs, and charts referenced from the README.
- `maskfile.md`: canonical developer workflow for tests, benchmark runs, graph generation, and README rebuilds.
- `scripts/plot.py`, `scripts/plot2.py`, `scripts/plots.json`: plotting helpers and benchmark-specific chart config.
- `pages/zz_notes.md`: benchmark environment notes; update it when published numbers are refreshed on different hardware or a different Go version.

## Working rules
- Read `README.md` and `maskfile.md` before changing benchmark logic, documentation structure, or asset generation.
- Do not treat `README.md` as the source of truth. Make doc edits in the `pages/*.md` stubs, then regenerate `README.md` with `mask readme`.
- Keep recommendations benchmark-scoped and probabilistic. Avoid blanket claims like "always faster" unless the measurements truly support that wording.
- New or updated sections should keep the same pattern: problem statement, short actionable tip block, chart, and a plain-language interpretation of the crossover.
- In the tip block, prefer approximate crossover guidance such as `~50`, `~1000`, or "hundreds of lookups" when the benchmark supports that level of precision.
- Keep the descriptive prose below the tip as short as possible: one tight paragraph that explains the measured crossover without rephrasing the whole section.
- Keep benchmark environment metadata centralized in `pages/zz_notes.md`; do not repeat machine or Go-version details inside individual benchmark sections unless the user explicitly asks for per-section provenance.
- Keep asset names aligned with benchmark names and variants so `benchstat`, plotting, and README links continue to work.
- Keep benchmark names, variant flags, and linked asset filenames stable unless you are updating every dependent page and command in the same change.
- The module path is `github.com/lunemec/go-rule-of-thumb/benchmarks` even though the directory name is `go-rules-of-thumb`; do not rename it casually.

## Change flow
1. Add or update the relevant benchmark and correctness test under `benchmarks/`.
2. Add or update the corresponding commands in `maskfile.md` if the workflow changed.
3. Regenerate the raw outputs, summary tables, and plots in `assets/` for the affected benchmark.
4. Add or update the matching `pages/*.md` stub.
5. Rebuild `README.md` with `mask readme`.
6. Verify that all image and benchmark-result links still resolve.

## Validation
- Docs-only page edits: update the relevant `pages/*.md` stub, then run `mask readme`
- Benchmark logic changes: `go test -v -timeout 60m -count 1 -race ./benchmarks/...`
- Topic-specific performance updates: run the relevant benchmark variants, then `benchstat`, then the graph command from `maskfile.md`
- If plotting inputs or naming conventions change, verify the plotting scripts still generate the expected files under `assets/`

## Out of scope unless asked
- Recasting microbenchmarks as universal Go guidance.
- Hand-editing generated benchmark outputs or charts.
- Large naming refactors across benchmark families, assets, and pages without updating the full pipeline.
