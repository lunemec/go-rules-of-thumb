package main

import (
	"fmt"
	"testing"
)

type ifSwitchBench func(int) bool

func BenchmarkIfSwitch(b *testing.B) {
	b.Run(
		fmt.Sprintf("if(1)"),
		benchmarkIfSwitch(benchIf1),
	)
	b.Run(
		fmt.Sprintf("switch(1)"),
		benchmarkIfSwitch(benchSwitch1),
	)
	b.Run(
		fmt.Sprintf("if(5)"),
		benchmarkIfSwitch(benchIf5),
	)
	b.Run(
		fmt.Sprintf("switch(5)"),
		benchmarkIfSwitch(benchSwitch5),
	)
}

func benchmarkIfSwitch(runF ifSwitchBench) func(*testing.B) {
	var input = 5

	return func(b *testing.B) {
		for b.Loop() {
			runF(input)
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
