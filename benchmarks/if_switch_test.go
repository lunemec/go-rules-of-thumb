package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type ifSwitchSliceRunner func([]int) int

var (
	ifSwitchCases     = []int{3, 5, 9}
	ifSwitchWorkloads = []string{"first", "middle", "last", "miss", "cycle", "random"}
	ifSwitchSink      int
)

const ifSwitchInputsPerOp = 256

func TestIfSwitchMatches(t *testing.T) {
	for _, cases := range ifSwitchCases {
		t.Run(fmt.Sprintf("cases=%d", cases), func(t *testing.T) {
			for in := -1; in <= cases+1; in++ {
				want := expectedIfSwitchResult(cases, in)

				require.Equal(t, want, dispatchIf(cases, in), "if input=%d", in)
				require.Equal(t, want, dispatchSwitch(cases, in), "switch input=%d", in)
			}
		})
	}
}

func BenchmarkIfSwitch(b *testing.B) {
	runBenchmark := func(selectRunner func(int) ifSwitchSliceRunner) {
		for _, cases := range ifSwitchCases {
			runF := selectRunner(cases)
			for _, workload := range ifSwitchWorkloads {
				inputs := buildIfSwitchInputs(cases, workload)
				b.Run(
					fmt.Sprintf("cases=%d workload=%s", cases, workload),
					benchmarkIfSwitch(inputs, runF),
				)
			}
		}
	}

	switch variant {
	case "if":
		runBenchmark(selectIfSwitchIfRunner)
	case "switch":
		runBenchmark(selectIfSwitchSwitchRunner)
	default:
		b.Errorf("specify which test to run: -args -variant if|switch")
	}
}

func benchmarkIfSwitch(inputs []int, runF ifSwitchSliceRunner) func(*testing.B) {
	return func(b *testing.B) {
		for b.Loop() {
			ifSwitchSink = runF(inputs)
		}
	}
}

func buildIfSwitchInputs(cases int, workload string) []int {
	first, middle, last, miss := ifSwitchRepresentativeInputs(cases)
	inputs := make([]int, ifSwitchInputsPerOp)

	switch workload {
	case "first":
		for i := range inputs {
			inputs[i] = first
		}
	case "middle":
		for i := range inputs {
			inputs[i] = middle
		}
	case "last":
		for i := range inputs {
			inputs[i] = last
		}
	case "miss":
		for i := range inputs {
			inputs[i] = miss
		}
	case "cycle":
		pattern := [...]int{first, middle, last, miss}
		for i := range inputs {
			inputs[i] = pattern[i%len(pattern)]
		}
	case "random":
		pattern := [...]int{first, middle, last, miss}
		rng := newDeterministicRand(stableSeed("BenchmarkIfSwitch", cases, workload))
		for i := range inputs {
			inputs[i] = pattern[rng.Intn(len(pattern))]
		}
	default:
		panic(fmt.Sprintf("unknown workload %q", workload))
	}

	return inputs
}

func ifSwitchRepresentativeInputs(cases int) (first, middle, last, miss int) {
	return 1, (cases / 2) + 1, cases, 0
}

func expectedIfSwitchResult(cases, in int) int {
	if in >= 1 && in <= cases {
		return in
	}
	return 0
}

func selectIfSwitchIfRunner(cases int) ifSwitchSliceRunner {
	switch cases {
	case 3:
		return runIf3Inputs
	case 5:
		return runIf5Inputs
	case 9:
		return runIf9Inputs
	default:
		panic(fmt.Sprintf("unsupported if cases %d", cases))
	}
}

func selectIfSwitchSwitchRunner(cases int) ifSwitchSliceRunner {
	switch cases {
	case 3:
		return runSwitch3Inputs
	case 5:
		return runSwitch5Inputs
	case 9:
		return runSwitch9Inputs
	default:
		panic(fmt.Sprintf("unsupported switch cases %d", cases))
	}
}

func dispatchIf(cases, in int) int {
	switch cases {
	case 3:
		return benchIf3(in)
	case 5:
		return benchIf5(in)
	case 9:
		return benchIf9(in)
	default:
		panic(fmt.Sprintf("unsupported if cases %d", cases))
	}
}

func dispatchSwitch(cases, in int) int {
	switch cases {
	case 3:
		return benchSwitch3(in)
	case 5:
		return benchSwitch5(in)
	case 9:
		return benchSwitch9(in)
	default:
		panic(fmt.Sprintf("unsupported switch cases %d", cases))
	}
}

func runIf3Inputs(inputs []int) int {
	acc := 0
	for _, in := range inputs {
		acc += benchIf3(in)
	}
	return acc
}

func runSwitch3Inputs(inputs []int) int {
	acc := 0
	for _, in := range inputs {
		acc += benchSwitch3(in)
	}
	return acc
}

func runIf5Inputs(inputs []int) int {
	acc := 0
	for _, in := range inputs {
		acc += benchIf5(in)
	}
	return acc
}

func runSwitch5Inputs(inputs []int) int {
	acc := 0
	for _, in := range inputs {
		acc += benchSwitch5(in)
	}
	return acc
}

func runIf9Inputs(inputs []int) int {
	acc := 0
	for _, in := range inputs {
		acc += benchIf9(in)
	}
	return acc
}

func runSwitch9Inputs(inputs []int) int {
	acc := 0
	for _, in := range inputs {
		acc += benchSwitch9(in)
	}
	return acc
}

func benchIf3(in int) int {
	if in == 1 {
		return 1
	} else if in == 2 {
		return 2
	} else if in == 3 {
		return 3
	}
	return 0
}

func benchSwitch3(in int) int {
	switch in {
	case 1:
		return 1
	case 2:
		return 2
	case 3:
		return 3
	default:
		return 0
	}
}

func benchIf5(in int) int {
	if in == 1 {
		return 1
	} else if in == 2 {
		return 2
	} else if in == 3 {
		return 3
	} else if in == 4 {
		return 4
	} else if in == 5 {
		return 5
	}
	return 0
}

func benchSwitch5(in int) int {
	switch in {
	case 1:
		return 1
	case 2:
		return 2
	case 3:
		return 3
	case 4:
		return 4
	case 5:
		return 5
	default:
		return 0
	}
}

func benchIf9(in int) int {
	if in == 1 {
		return 1
	} else if in == 2 {
		return 2
	} else if in == 3 {
		return 3
	} else if in == 4 {
		return 4
	} else if in == 5 {
		return 5
	} else if in == 6 {
		return 6
	} else if in == 7 {
		return 7
	} else if in == 8 {
		return 8
	} else if in == 9 {
		return 9
	}
	return 0
}

func benchSwitch9(in int) int {
	switch in {
	case 1:
		return 1
	case 2:
		return 2
	case 3:
		return 3
	case 4:
		return 4
	case 5:
		return 5
	case 6:
		return 6
	case 7:
		return 7
	case 8:
		return 8
	case 9:
		return 9
	default:
		return 0
	}
}
