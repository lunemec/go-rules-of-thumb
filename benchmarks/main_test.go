package main

import "flag"

func init() {
	flag.StringVar(&variant, "variant", "", "specifies which variant to run")
}

var (
	variant string
	sizes   = []int{10, 100, 500, 1_000, 5_000, 10_000, 50_000, 100_000, 500_000, 1_000_000, 5_000_000, 10_000_000}
	needles = []int{10, 100, 500, 1_000}
)
