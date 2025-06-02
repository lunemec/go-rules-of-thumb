package main

import "flag"

func init() {
	flag.StringVar(&variant, "variant", "", "specifies which variant to run")
}

var (
	variant string
	sizes   = []int{10, 100, 500, 1_000, 5_000, 10_000, 50_000, 100_000, 500_000, 1_000_000}
	// sizesReduced is used for tests where it does not make sense to do that many iterations.
	sizesReduced = []int{1, 10, 100, 1_000}
	needles      = []int{10, 100, 500, 1_000}
)
