package main

import (
	"fmt"
	"testing"
)

type ifSwitchBench func(int) bool

func BenchmarkIfSwitch(b *testing.B) {
	runBenchmark := func(runF ifSwitchBench) {
		for _, size := range sizesReduced {
			b.Run(
				fmt.Sprintf("size=%d", size),
				benchmarkIfSwitch(size, runF),
			)
		}
	}
	switch variant {
	case "if":
		runBenchmark(benchIf1)
	case "switch":
		runBenchmark(benchSwitch1)
	case "if_5":
		runBenchmark(benchIf5)
	case "switch_5":
		runBenchmark(benchSwitch5)
	default:
		b.Errorf("speficy which test to run: -args -test if|if_5|switch|switch_5")
	}
}

func benchmarkIfSwitch(size int, runF ifSwitchBench) func(*testing.B) {
	input := 5

	return func(b *testing.B) {
		for b.Loop() {
			for range size {
				runF(input)
			}
		}
	}
}

func benchIf1(in int) bool {
	return in == 5
}

func benchSwitch1(in int) bool {
	switch in {
	case 5:
		return true
	default:
		return false
	}
}

func benchIf5(in int) bool {
	if in == 1 {
		return false
	} else if in == 2 {
		return false
	} else if in == 3 {
		return false
	} else if in == 4 {
		return false
	} else if in == 5 {
		return true
	}
	return false
}

func benchSwitch5(in int) bool {
	switch in {
	case 1:
		return false
	case 2:
		return false
	case 3:
		return false
	case 4:
		return false
	case 5:
		return true
	default:
		return false
	}
}
