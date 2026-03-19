package main

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

type hotRequest interface {
	HotScore() uint64
}

type hotMutatingRequest interface {
	HotMutate() uint64
}

type tinyRequest interface {
	TinyScore() uint64
}

type tinyMutatingRequest interface {
	TinyMutate() uint64
}

var (
	callShapesSink           uint64
	callShapeIterationCounts = []int{1, 2, 5, 10, 25, 50, 100, 250, 500, 1_000}
	callShapeTestIterations  = []int{1, 10, 100}
)

func TestCallShapesOutputsMatch(t *testing.T) {
	t.Run("payload_sizes", func(t *testing.T) {
		require.Equal(t, uintptr(8), unsafe.Sizeof(scoringRequest8{}))
		require.Equal(t, uintptr(16), unsafe.Sizeof(scoringRequest16{}))
		require.Equal(t, uintptr(24), unsafe.Sizeof(scoringRequest24{}))
		require.Equal(t, uintptr(32), unsafe.Sizeof(scoringRequest32{}))
		require.Equal(t, uintptr(64), unsafe.Sizeof(scoringRequest64{}))
		require.Equal(t, uintptr(96), unsafe.Sizeof(scoringRequest96{}))
		require.Equal(t, uintptr(128), unsafe.Sizeof(scoringRequest128{}))
		require.Equal(t, uintptr(192), unsafe.Sizeof(scoringRequest192{}))
		require.Equal(t, uintptr(256), unsafe.Sizeof(scoringRequest256{}))
		require.Equal(t, uintptr(384), unsafe.Sizeof(scoringRequest384{}))
		require.Equal(t, uintptr(512), unsafe.Sizeof(scoringRequest512{}))
	})

	t.Run("size=8", func(t *testing.T) {
		requireCallShapeOutputsMatch(t, pointerSlice(buildScoringRequests8(paramTestRecords)), sumCallShapes8Concrete, sumCallShapes8GenericExact)
	})
	t.Run("size=16", func(t *testing.T) {
		requireCallShapeOutputsMatch(t, pointerSlice(buildScoringRequests16(paramTestRecords)), sumCallShapes16Concrete, sumCallShapes16GenericExact)
	})
	t.Run("size=24", func(t *testing.T) {
		requireCallShapeOutputsMatch(t, pointerSlice(buildScoringRequests24(paramTestRecords)), sumCallShapes24Concrete, sumCallShapes24GenericExact)
	})
	t.Run("size=32", func(t *testing.T) {
		requireCallShapeOutputsMatch(t, pointerSlice(buildScoringRequests32(paramTestRecords)), sumCallShapes32Concrete, sumCallShapes32GenericExact)
	})
	t.Run("size=64", func(t *testing.T) {
		requireCallShapeOutputsMatch(t, pointerSlice(buildScoringRequests64(paramTestRecords)), sumCallShapes64Concrete, sumCallShapes64GenericExact)
	})
	t.Run("size=96", func(t *testing.T) {
		requireCallShapeOutputsMatch(t, pointerSlice(buildScoringRequests96(paramTestRecords)), sumCallShapes96Concrete, sumCallShapes96GenericExact)
	})
	t.Run("size=128", func(t *testing.T) {
		requireCallShapeOutputsMatch(t, pointerSlice(buildScoringRequests128(paramTestRecords)), sumCallShapes128Concrete, sumCallShapes128GenericExact)
	})
	t.Run("size=192", func(t *testing.T) {
		requireCallShapeOutputsMatch(t, pointerSlice(buildScoringRequests192(paramTestRecords)), sumCallShapes192Concrete, sumCallShapes192GenericExact)
	})
	t.Run("size=256", func(t *testing.T) {
		requireCallShapeOutputsMatch(t, pointerSlice(buildScoringRequests256(paramTestRecords)), sumCallShapes256Concrete, sumCallShapes256GenericExact)
	})
	t.Run("size=384", func(t *testing.T) {
		requireCallShapeOutputsMatch(t, pointerSlice(buildScoringRequests384(paramTestRecords)), sumCallShapes384Concrete, sumCallShapes384GenericExact)
	})
	t.Run("size=512", func(t *testing.T) {
		requireCallShapeOutputsMatch(t, pointerSlice(buildScoringRequests512(paramTestRecords)), sumCallShapes512Concrete, sumCallShapes512GenericExact)
	})
}

func TestIndirectionCostReadOutputsMatch(t *testing.T) {
	require.Equal(t, uintptr(32), unsafe.Sizeof(scoringRequest32{}))

	requests := pointerSlice(buildScoringRequests32(paramTestRecords))
	interfaces := interfaceSlice(requests)

	for _, iterations := range callShapeTestIterations {
		t.Run(fmt.Sprintf("iterations=%d", iterations), func(t *testing.T) {
			expected := sumCallShapes32ConcreteRepeated(requests, iterations)
			require.Equal(t, expected, sumHotRequestsRepeated(interfaces, iterations))
			require.Equal(t, expected, sumHotRequestsBoxEachCallRepeated(requests, iterations))
			require.Equal(t, expected, sumCallShapes32GenericExactRepeated(requests, iterations))
			require.Equal(t, expected, sumRequestsGenericConstraintRepeated(requests, iterations))
		})
	}
}

func TestIndirectionCostWriteOutputsMatch(t *testing.T) {
	require.Equal(t, uintptr(32), unsafe.Sizeof(scoringRequest32{}))

	base := buildScoringRequests32(paramTestRecords)

	for _, iterations := range callShapeTestIterations {
		t.Run(fmt.Sprintf("iterations=%d", iterations), func(t *testing.T) {
			concreteRequests := pointerSlice(cloneSlice(base))
			expected := sumMutatingCallShapes32ConcreteRepeated(concreteRequests, iterations)

			interfaceDispatchRequests := pointerSlice(cloneSlice(base))
			require.Equal(
				t,
				expected,
				sumHotMutatingRequestsRepeated(interfaceMutatingSlice(interfaceDispatchRequests), iterations),
			)

			interfaceBoxRequests := pointerSlice(cloneSlice(base))
			require.Equal(
				t,
				expected,
				sumHotMutatingRequestsBoxEachCallRepeated(interfaceBoxRequests, iterations),
			)

			genericExactRequests := pointerSlice(cloneSlice(base))
			require.Equal(
				t,
				expected,
				sumMutatingCallShapes32GenericExactRepeated(genericExactRequests, iterations),
			)

			genericConstraintRequests := pointerSlice(cloneSlice(base))
			require.Equal(
				t,
				expected,
				sumMutatingRequestsGenericConstraintRepeated(genericConstraintRequests, iterations),
			)
		})
	}
}

func TestIndirectionCostTinyReadOutputsMatch(t *testing.T) {
	require.Equal(t, uintptr(32), unsafe.Sizeof(scoringRequest32{}))

	requests := pointerSlice(buildScoringRequests32(paramTestRecords))
	interfaces := tinyInterfaceSlice(requests)

	for _, iterations := range callShapeTestIterations {
		t.Run(fmt.Sprintf("iterations=%d", iterations), func(t *testing.T) {
			expected := sumTinyCallShapes32ConcreteRepeated(requests, iterations)
			require.Equal(t, expected, sumTinyRequestsRepeated(interfaces, iterations))
			require.Equal(t, expected, sumTinyRequestsBoxEachCallRepeated(requests, iterations))
			require.Equal(t, expected, sumTinyCallShapes32GenericExactRepeated(requests, iterations))
			require.Equal(t, expected, sumTinyRequestsGenericConstraintRepeated(requests, iterations))
		})
	}
}

func TestIndirectionCostTinyWriteOutputsMatch(t *testing.T) {
	require.Equal(t, uintptr(32), unsafe.Sizeof(scoringRequest32{}))

	base := buildScoringRequests32(paramTestRecords)

	for _, iterations := range callShapeTestIterations {
		t.Run(fmt.Sprintf("iterations=%d", iterations), func(t *testing.T) {
			concreteRequests := pointerSlice(cloneSlice(base))
			expected := sumTinyMutatingCallShapes32ConcreteRepeated(concreteRequests, iterations)

			interfaceDispatchRequests := pointerSlice(cloneSlice(base))
			require.Equal(
				t,
				expected,
				sumTinyMutatingRequestsRepeated(tinyInterfaceMutatingSlice(interfaceDispatchRequests), iterations),
			)

			interfaceBoxRequests := pointerSlice(cloneSlice(base))
			require.Equal(
				t,
				expected,
				sumTinyMutatingRequestsBoxEachCallRepeated(interfaceBoxRequests, iterations),
			)

			genericExactRequests := pointerSlice(cloneSlice(base))
			require.Equal(
				t,
				expected,
				sumTinyMutatingCallShapes32GenericExactRepeated(genericExactRequests, iterations),
			)

			genericConstraintRequests := pointerSlice(cloneSlice(base))
			require.Equal(
				t,
				expected,
				sumTinyMutatingRequestsGenericConstraintRepeated(genericConstraintRequests, iterations),
			)
		})
	}
}

func BenchmarkCallShapes(b *testing.B) {
	requests8 := pointerSlice(buildScoringRequests8(paramBenchmarkRecords))
	requests16 := pointerSlice(buildScoringRequests16(paramBenchmarkRecords))
	requests24 := pointerSlice(buildScoringRequests24(paramBenchmarkRecords))
	requests32 := pointerSlice(buildScoringRequests32(paramBenchmarkRecords))
	requests64 := pointerSlice(buildScoringRequests64(paramBenchmarkRecords))
	requests96 := pointerSlice(buildScoringRequests96(paramBenchmarkRecords))
	requests128 := pointerSlice(buildScoringRequests128(paramBenchmarkRecords))
	requests192 := pointerSlice(buildScoringRequests192(paramBenchmarkRecords))
	requests256 := pointerSlice(buildScoringRequests256(paramBenchmarkRecords))
	requests384 := pointerSlice(buildScoringRequests384(paramBenchmarkRecords))
	requests512 := pointerSlice(buildScoringRequests512(paramBenchmarkRecords))

	cases := []struct {
		size                   int
		concreteBench          func(*testing.B)
		interfaceBench         func(*testing.B)
		interfaceBoxingBench   func(*testing.B)
		genericExactBench      func(*testing.B)
		genericConstraintBench func(*testing.B)
	}{
		{
			size: int(unsafe.Sizeof(scoringRequest8{})),
			concreteBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes8Concrete(requests8)
				}
			},
			interfaceBench: benchmarkCallShapesInterface(requests8),
			interfaceBoxingBench: benchmarkCallShapesInterfaceBoxEachCall(
				requests8,
			),
			genericExactBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes8GenericExact(requests8)
				}
			},
			genericConstraintBench: benchmarkCallShapesGenericConstraint(requests8),
		},
		{
			size: int(unsafe.Sizeof(scoringRequest16{})),
			concreteBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes16Concrete(requests16)
				}
			},
			interfaceBench: benchmarkCallShapesInterface(requests16),
			interfaceBoxingBench: benchmarkCallShapesInterfaceBoxEachCall(
				requests16,
			),
			genericExactBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes16GenericExact(requests16)
				}
			},
			genericConstraintBench: benchmarkCallShapesGenericConstraint(requests16),
		},
		{
			size: int(unsafe.Sizeof(scoringRequest24{})),
			concreteBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes24Concrete(requests24)
				}
			},
			interfaceBench: benchmarkCallShapesInterface(requests24),
			interfaceBoxingBench: benchmarkCallShapesInterfaceBoxEachCall(
				requests24,
			),
			genericExactBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes24GenericExact(requests24)
				}
			},
			genericConstraintBench: benchmarkCallShapesGenericConstraint(requests24),
		},
		{
			size: int(unsafe.Sizeof(scoringRequest32{})),
			concreteBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes32Concrete(requests32)
				}
			},
			interfaceBench: benchmarkCallShapesInterface(requests32),
			interfaceBoxingBench: benchmarkCallShapesInterfaceBoxEachCall(
				requests32,
			),
			genericExactBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes32GenericExact(requests32)
				}
			},
			genericConstraintBench: benchmarkCallShapesGenericConstraint(requests32),
		},
		{
			size: int(unsafe.Sizeof(scoringRequest64{})),
			concreteBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes64Concrete(requests64)
				}
			},
			interfaceBench: benchmarkCallShapesInterface(requests64),
			interfaceBoxingBench: benchmarkCallShapesInterfaceBoxEachCall(
				requests64,
			),
			genericExactBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes64GenericExact(requests64)
				}
			},
			genericConstraintBench: benchmarkCallShapesGenericConstraint(requests64),
		},
		{
			size: int(unsafe.Sizeof(scoringRequest96{})),
			concreteBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes96Concrete(requests96)
				}
			},
			interfaceBench: benchmarkCallShapesInterface(requests96),
			interfaceBoxingBench: benchmarkCallShapesInterfaceBoxEachCall(
				requests96,
			),
			genericExactBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes96GenericExact(requests96)
				}
			},
			genericConstraintBench: benchmarkCallShapesGenericConstraint(requests96),
		},
		{
			size: int(unsafe.Sizeof(scoringRequest128{})),
			concreteBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes128Concrete(requests128)
				}
			},
			interfaceBench: benchmarkCallShapesInterface(requests128),
			interfaceBoxingBench: benchmarkCallShapesInterfaceBoxEachCall(
				requests128,
			),
			genericExactBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes128GenericExact(requests128)
				}
			},
			genericConstraintBench: benchmarkCallShapesGenericConstraint(requests128),
		},
		{
			size: int(unsafe.Sizeof(scoringRequest192{})),
			concreteBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes192Concrete(requests192)
				}
			},
			interfaceBench: benchmarkCallShapesInterface(requests192),
			interfaceBoxingBench: benchmarkCallShapesInterfaceBoxEachCall(
				requests192,
			),
			genericExactBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes192GenericExact(requests192)
				}
			},
			genericConstraintBench: benchmarkCallShapesGenericConstraint(requests192),
		},
		{
			size: int(unsafe.Sizeof(scoringRequest256{})),
			concreteBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes256Concrete(requests256)
				}
			},
			interfaceBench: benchmarkCallShapesInterface(requests256),
			interfaceBoxingBench: benchmarkCallShapesInterfaceBoxEachCall(
				requests256,
			),
			genericExactBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes256GenericExact(requests256)
				}
			},
			genericConstraintBench: benchmarkCallShapesGenericConstraint(requests256),
		},
		{
			size: int(unsafe.Sizeof(scoringRequest384{})),
			concreteBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes384Concrete(requests384)
				}
			},
			interfaceBench: benchmarkCallShapesInterface(requests384),
			interfaceBoxingBench: benchmarkCallShapesInterfaceBoxEachCall(
				requests384,
			),
			genericExactBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes384GenericExact(requests384)
				}
			},
			genericConstraintBench: benchmarkCallShapesGenericConstraint(requests384),
		},
		{
			size: int(unsafe.Sizeof(scoringRequest512{})),
			concreteBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes512Concrete(requests512)
				}
			},
			interfaceBench: benchmarkCallShapesInterface(requests512),
			interfaceBoxingBench: benchmarkCallShapesInterfaceBoxEachCall(
				requests512,
			),
			genericExactBench: func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes512GenericExact(requests512)
				}
			},
			genericConstraintBench: benchmarkCallShapesGenericConstraint(requests512),
		},
	}

	for _, tc := range cases {
		run := tc.concreteBench
		switch variant {
		case "concrete_ptr":
		case "interface_dispatch":
			run = tc.interfaceBench
		case "interface_box_each_call":
			run = tc.interfaceBoxingBench
		case "generic_exact":
			run = tc.genericExactBench
		case "generic_constraint":
			run = tc.genericConstraintBench
		default:
			b.Fatalf("specify which benchmark to run: -args -variant concrete_ptr|interface_dispatch|interface_box_each_call|generic_exact|generic_constraint")
		}

		b.Run(fmt.Sprintf("size=%d", tc.size), run)
	}
}

func BenchmarkIndirectionCostRead(b *testing.B) {
	requests32 := pointerSlice(buildScoringRequests32(paramBenchmarkRecords))
	interfaces32 := interfaceSlice(requests32)
	for _, iterations := range callShapeIterationCounts {
		run := func(b *testing.B) {
			for b.Loop() {
				callShapesSink = sumCallShapes32ConcreteRepeated(requests32, iterations)
			}
		}
		switch variant {
		case "concrete_ptr":
		case "interface_dispatch":
			run = func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumHotRequestsRepeated(interfaces32, iterations)
				}
			}
		case "interface_box_each_call":
			run = func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumHotRequestsBoxEachCallRepeated(requests32, iterations)
				}
			}
		case "generic_exact":
			run = func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumCallShapes32GenericExactRepeated(requests32, iterations)
				}
			}
		case "generic_constraint":
			run = func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumRequestsGenericConstraintRepeated(requests32, iterations)
				}
			}
		default:
			b.Fatalf("specify which benchmark to run: -args -variant concrete_ptr|interface_dispatch|interface_box_each_call|generic_exact|generic_constraint")
		}

		b.Run(fmt.Sprintf("iterations=%d", iterations), run)
	}
}

func BenchmarkIndirectionCostWrite(b *testing.B) {
	for _, iterations := range callShapeIterationCounts {
		run := func(b *testing.B) {
			requests32 := pointerSlice(buildScoringRequests32(paramBenchmarkRecords))
			for b.Loop() {
				callShapesSink = sumMutatingCallShapes32ConcreteRepeated(requests32, iterations)
			}
		}

		switch variant {
		case "concrete_ptr":
		case "interface_dispatch":
			run = func(b *testing.B) {
				requests32 := pointerSlice(buildScoringRequests32(paramBenchmarkRecords))
				interfaces32 := interfaceMutatingSlice(requests32)
				for b.Loop() {
					callShapesSink = sumHotMutatingRequestsRepeated(interfaces32, iterations)
				}
			}
		case "interface_box_each_call":
			run = func(b *testing.B) {
				requests32 := pointerSlice(buildScoringRequests32(paramBenchmarkRecords))
				for b.Loop() {
					callShapesSink = sumHotMutatingRequestsBoxEachCallRepeated(requests32, iterations)
				}
			}
		case "generic_exact":
			run = func(b *testing.B) {
				requests32 := pointerSlice(buildScoringRequests32(paramBenchmarkRecords))
				for b.Loop() {
					callShapesSink = sumMutatingCallShapes32GenericExactRepeated(requests32, iterations)
				}
			}
		case "generic_constraint":
			run = func(b *testing.B) {
				requests32 := pointerSlice(buildScoringRequests32(paramBenchmarkRecords))
				for b.Loop() {
					callShapesSink = sumMutatingRequestsGenericConstraintRepeated(requests32, iterations)
				}
			}
		default:
			b.Fatalf("specify which benchmark to run: -args -variant concrete_ptr|interface_dispatch|interface_box_each_call|generic_exact|generic_constraint")
		}

		b.Run(fmt.Sprintf("iterations=%d", iterations), run)
	}
}

func BenchmarkIndirectionCostTinyRead(b *testing.B) {
	requests32 := pointerSlice(buildScoringRequests32(paramBenchmarkRecords))
	interfaces32 := tinyInterfaceSlice(requests32)
	for _, iterations := range callShapeIterationCounts {
		run := func(b *testing.B) {
			for b.Loop() {
				callShapesSink = sumTinyCallShapes32ConcreteRepeated(requests32, iterations)
			}
		}
		switch variant {
		case "concrete_ptr":
		case "interface_dispatch":
			run = func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumTinyRequestsRepeated(interfaces32, iterations)
				}
			}
		case "interface_box_each_call":
			run = func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumTinyRequestsBoxEachCallRepeated(requests32, iterations)
				}
			}
		case "generic_exact":
			run = func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumTinyCallShapes32GenericExactRepeated(requests32, iterations)
				}
			}
		case "generic_constraint":
			run = func(b *testing.B) {
				for b.Loop() {
					callShapesSink = sumTinyRequestsGenericConstraintRepeated(requests32, iterations)
				}
			}
		default:
			b.Fatalf("specify which benchmark to run: -args -variant concrete_ptr|interface_dispatch|interface_box_each_call|generic_exact|generic_constraint")
		}

		b.Run(fmt.Sprintf("iterations=%d", iterations), run)
	}
}

func BenchmarkIndirectionCostTinyWrite(b *testing.B) {
	for _, iterations := range callShapeIterationCounts {
		run := func(b *testing.B) {
			requests32 := pointerSlice(buildScoringRequests32(paramBenchmarkRecords))
			for b.Loop() {
				callShapesSink = sumTinyMutatingCallShapes32ConcreteRepeated(requests32, iterations)
			}
		}

		switch variant {
		case "concrete_ptr":
		case "interface_dispatch":
			run = func(b *testing.B) {
				requests32 := pointerSlice(buildScoringRequests32(paramBenchmarkRecords))
				interfaces32 := tinyInterfaceMutatingSlice(requests32)
				for b.Loop() {
					callShapesSink = sumTinyMutatingRequestsRepeated(interfaces32, iterations)
				}
			}
		case "interface_box_each_call":
			run = func(b *testing.B) {
				requests32 := pointerSlice(buildScoringRequests32(paramBenchmarkRecords))
				for b.Loop() {
					callShapesSink = sumTinyMutatingRequestsBoxEachCallRepeated(requests32, iterations)
				}
			}
		case "generic_exact":
			run = func(b *testing.B) {
				requests32 := pointerSlice(buildScoringRequests32(paramBenchmarkRecords))
				for b.Loop() {
					callShapesSink = sumTinyMutatingCallShapes32GenericExactRepeated(requests32, iterations)
				}
			}
		case "generic_constraint":
			run = func(b *testing.B) {
				requests32 := pointerSlice(buildScoringRequests32(paramBenchmarkRecords))
				for b.Loop() {
					callShapesSink = sumTinyMutatingRequestsGenericConstraintRepeated(requests32, iterations)
				}
			}
		default:
			b.Fatalf("specify which benchmark to run: -args -variant concrete_ptr|interface_dispatch|interface_box_each_call|generic_exact|generic_constraint")
		}

		b.Run(fmt.Sprintf("iterations=%d", iterations), run)
	}
}

func requireCallShapeOutputsMatch[T hotRequest](t *testing.T, requests []T, sumConcrete func([]T) uint64, sumGenericExact func([]T) uint64) {
	t.Helper()

	expected := sumConcrete(requests)
	require.Equal(t, expected, sumHotRequests(interfaceSlice(requests)))
	require.Equal(t, expected, sumHotRequestsBoxEachCall(requests))
	require.Equal(t, expected, sumGenericExact(requests))
	require.Equal(t, expected, sumRequestsGenericConstraint(requests))
}

func pointerSlice[T any](requests []T) []*T {
	pointers := make([]*T, len(requests))
	for i := range requests {
		pointers[i] = &requests[i]
	}
	return pointers
}

func cloneSlice[T any](requests []T) []T {
	cloned := make([]T, len(requests))
	copy(cloned, requests)
	return cloned
}

func interfaceSlice[T hotRequest](requests []T) []hotRequest {
	interfaces := make([]hotRequest, len(requests))
	for i := range requests {
		interfaces[i] = requests[i]
	}
	return interfaces
}

func interfaceMutatingSlice[T hotMutatingRequest](requests []T) []hotMutatingRequest {
	interfaces := make([]hotMutatingRequest, len(requests))
	for i := range requests {
		interfaces[i] = requests[i]
	}
	return interfaces
}

func tinyInterfaceSlice[T tinyRequest](requests []T) []tinyRequest {
	interfaces := make([]tinyRequest, len(requests))
	for i := range requests {
		interfaces[i] = requests[i]
	}
	return interfaces
}

func tinyInterfaceMutatingSlice[T tinyMutatingRequest](requests []T) []tinyMutatingRequest {
	interfaces := make([]tinyMutatingRequest, len(requests))
	for i := range requests {
		interfaces[i] = requests[i]
	}
	return interfaces
}

func benchmarkCallShapesInterface[T hotRequest](requests []T) func(*testing.B) {
	interfaces := interfaceSlice(requests)
	return func(b *testing.B) {
		for b.Loop() {
			callShapesSink = sumHotRequests(interfaces)
		}
	}
}

func benchmarkCallShapesInterfaceBoxEachCall[T hotRequest](requests []T) func(*testing.B) {
	return func(b *testing.B) {
		for b.Loop() {
			callShapesSink = sumHotRequestsBoxEachCall(requests)
		}
	}
}

func benchmarkCallShapesGenericConstraint[T hotRequest](requests []T) func(*testing.B) {
	return func(b *testing.B) {
		for b.Loop() {
			callShapesSink = sumRequestsGenericConstraint(requests)
		}
	}
}

func sumHotRequests(requests []hotRequest) uint64 {
	var acc uint64
	for i := range requests {
		acc += requests[i].HotScore()
	}
	return acc
}

func sumHotRequestsBoxEachCall[T hotRequest](requests []T) uint64 {
	var acc uint64
	for i := range requests {
		var request hotRequest = requests[i]
		acc += request.HotScore()
	}
	return acc
}

func sumRequestsGenericConstraint[T hotRequest](requests []T) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequestGenericConstraint(requests[i])
	}
	return acc
}

func sumHotRequestsRepeated(requests []hotRequest, iterations int) uint64 {
	var acc uint64
	for range iterations {
		for i := range requests {
			acc += requests[i].HotScore()
		}
	}
	return acc
}

func sumHotMutatingRequestsRepeated(requests []hotMutatingRequest, iterations int) uint64 {
	var acc uint64
	for range iterations {
		for i := range requests {
			acc += requests[i].HotMutate()
		}
	}
	return acc
}

func sumTinyRequestsRepeated(requests []tinyRequest, iterations int) uint64 {
	var acc uint64
	for range iterations {
		for i := range requests {
			acc += requests[i].TinyScore()
		}
	}
	return acc
}

func sumTinyMutatingRequestsRepeated(requests []tinyMutatingRequest, iterations int) uint64 {
	var acc uint64
	for range iterations {
		for i := range requests {
			acc += requests[i].TinyMutate()
		}
	}
	return acc
}

func sumHotRequestsBoxEachCallRepeated[T hotRequest](requests []T, iterations int) uint64 {
	var acc uint64
	for range iterations {
		for i := range requests {
			var request hotRequest = requests[i]
			acc += request.HotScore()
		}
	}
	return acc
}

func sumHotMutatingRequestsBoxEachCallRepeated[T hotMutatingRequest](requests []T, iterations int) uint64 {
	var acc uint64
	for range iterations {
		for i := range requests {
			var request hotMutatingRequest = requests[i]
			acc += request.HotMutate()
		}
	}
	return acc
}

func sumTinyRequestsBoxEachCallRepeated[T tinyRequest](requests []T, iterations int) uint64 {
	var acc uint64
	for range iterations {
		for i := range requests {
			var request tinyRequest = requests[i]
			acc += request.TinyScore()
		}
	}
	return acc
}

func sumTinyMutatingRequestsBoxEachCallRepeated[T tinyMutatingRequest](requests []T, iterations int) uint64 {
	var acc uint64
	for range iterations {
		for i := range requests {
			var request tinyMutatingRequest = requests[i]
			acc += request.TinyMutate()
		}
	}
	return acc
}

func sumRequestsGenericConstraintRepeated[T hotRequest](requests []T, iterations int) uint64 {
	var acc uint64
	for range iterations {
		for i := range requests {
			acc += scoreRequestGenericConstraint(requests[i])
		}
	}
	return acc
}

func sumTinyRequestsGenericConstraintRepeated[T tinyRequest](requests []T, iterations int) uint64 {
	var acc uint64
	for range iterations {
		for i := range requests {
			acc += scoreTinyRequestGenericConstraint(requests[i])
		}
	}
	return acc
}

func sumTinyMutatingRequestsGenericConstraintRepeated[T tinyMutatingRequest](requests []T, iterations int) uint64 {
	var acc uint64
	for range iterations {
		for i := range requests {
			acc += mutateTinyRequestGenericConstraint(requests[i])
		}
	}
	return acc
}

func sumMutatingRequestsGenericConstraintRepeated[T hotMutatingRequest](requests []T, iterations int) uint64 {
	var acc uint64
	for range iterations {
		for i := range requests {
			acc += mutateRequestGenericConstraint(requests[i])
		}
	}
	return acc
}

func sumTinyCallShapes32ConcreteRepeated(requests []*scoringRequest32, iterations int) uint64 {
	var acc uint64
	for range iterations {
		for i := range requests {
			acc += scoreRequest32TinyConcrete(requests[i])
		}
	}
	return acc
}

func sumTinyMutatingCallShapes32ConcreteRepeated(requests []*scoringRequest32, iterations int) uint64 {
	var acc uint64
	for range iterations {
		for i := range requests {
			acc += mutateRequest32TinyConcrete(requests[i])
		}
	}
	return acc
}

func sumCallShapes32ConcreteRepeated(requests []*scoringRequest32, iterations int) uint64 {
	var acc uint64
	for range iterations {
		for i := range requests {
			acc += scoreRequest32Concrete(requests[i])
		}
	}
	return acc
}

func sumMutatingCallShapes32ConcreteRepeated(requests []*scoringRequest32, iterations int) uint64 {
	var acc uint64
	for range iterations {
		for i := range requests {
			acc += mutateRequest32Concrete(requests[i])
		}
	}
	return acc
}

func sumTinyCallShapes32GenericExactRepeated(requests []*scoringRequest32, iterations int) uint64 {
	var acc uint64
	for range iterations {
		for i := range requests {
			acc += scoreRequest32TinyGenericExact(requests[i])
		}
	}
	return acc
}

func sumTinyMutatingCallShapes32GenericExactRepeated(requests []*scoringRequest32, iterations int) uint64 {
	var acc uint64
	for range iterations {
		for i := range requests {
			acc += mutateRequest32TinyGenericExact(requests[i])
		}
	}
	return acc
}

func sumCallShapes32GenericExactRepeated(requests []*scoringRequest32, iterations int) uint64 {
	var acc uint64
	for range iterations {
		for i := range requests {
			acc += scoreRequest32GenericExact(requests[i])
		}
	}
	return acc
}

func sumMutatingCallShapes32GenericExactRepeated(requests []*scoringRequest32, iterations int) uint64 {
	var acc uint64
	for range iterations {
		for i := range requests {
			acc += mutateRequest32GenericExact(requests[i])
		}
	}
	return acc
}

func sumCallShapes8Concrete(requests []*scoringRequest8) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest8Concrete(requests[i])
	}
	return acc
}

func sumCallShapes8GenericExact(requests []*scoringRequest8) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest8GenericExact(requests[i])
	}
	return acc
}

func sumCallShapes16Concrete(requests []*scoringRequest16) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest16Concrete(requests[i])
	}
	return acc
}

func sumCallShapes16GenericExact(requests []*scoringRequest16) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest16GenericExact(requests[i])
	}
	return acc
}

func sumCallShapes24Concrete(requests []*scoringRequest24) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest24Concrete(requests[i])
	}
	return acc
}

func sumCallShapes24GenericExact(requests []*scoringRequest24) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest24GenericExact(requests[i])
	}
	return acc
}

func sumCallShapes32Concrete(requests []*scoringRequest32) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest32Concrete(requests[i])
	}
	return acc
}

func sumCallShapes32GenericExact(requests []*scoringRequest32) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest32GenericExact(requests[i])
	}
	return acc
}

func sumCallShapes64Concrete(requests []*scoringRequest64) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest64Concrete(requests[i])
	}
	return acc
}

func sumCallShapes64GenericExact(requests []*scoringRequest64) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest64GenericExact(requests[i])
	}
	return acc
}

func sumCallShapes96Concrete(requests []*scoringRequest96) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest96Concrete(requests[i])
	}
	return acc
}

func sumCallShapes96GenericExact(requests []*scoringRequest96) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest96GenericExact(requests[i])
	}
	return acc
}

func sumCallShapes128Concrete(requests []*scoringRequest128) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest128Concrete(requests[i])
	}
	return acc
}

func sumCallShapes128GenericExact(requests []*scoringRequest128) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest128GenericExact(requests[i])
	}
	return acc
}

func sumCallShapes192Concrete(requests []*scoringRequest192) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest192Concrete(requests[i])
	}
	return acc
}

func sumCallShapes192GenericExact(requests []*scoringRequest192) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest192GenericExact(requests[i])
	}
	return acc
}

func sumCallShapes256Concrete(requests []*scoringRequest256) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest256Concrete(requests[i])
	}
	return acc
}

func sumCallShapes256GenericExact(requests []*scoringRequest256) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest256GenericExact(requests[i])
	}
	return acc
}

func sumCallShapes384Concrete(requests []*scoringRequest384) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest384Concrete(requests[i])
	}
	return acc
}

func sumCallShapes384GenericExact(requests []*scoringRequest384) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest384GenericExact(requests[i])
	}
	return acc
}

func sumCallShapes512Concrete(requests []*scoringRequest512) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest512Concrete(requests[i])
	}
	return acc
}

func sumCallShapes512GenericExact(requests []*scoringRequest512) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest512GenericExact(requests[i])
	}
	return acc
}

//go:noinline
func scoreRequestGenericConstraint[T hotRequest](request T) uint64 {
	return request.HotScore()
}

//go:noinline
func mutateRequestGenericConstraint[T hotMutatingRequest](request T) uint64 {
	return request.HotMutate()
}

//go:noinline
func scoreTinyRequestGenericConstraint[T tinyRequest](request T) uint64 {
	return request.TinyScore()
}

//go:noinline
func mutateTinyRequestGenericConstraint[T tinyMutatingRequest](request T) uint64 {
	return request.TinyMutate()
}

//go:noinline
func scoreRequest8Concrete(request *scoringRequest8) uint64 {
	return request.Core.AccountID
}

//go:noinline
func (request *scoringRequest8) HotScore() uint64 {
	return request.Core.AccountID
}

//go:noinline
func scoreRequest8GenericExact[T *scoringRequest8](request T) uint64 {
	concrete := (*scoringRequest8)(request)
	return concrete.Core.AccountID
}

//go:noinline
func scoreRequest16Concrete(request *scoringRequest16) uint64 {
	acc := request.Core.AccountID
	acc += uint64(request.Core.Score)
	acc += uint64(request.Core.Attempts)
	acc += uint64(request.Core.Flags)
	return acc
}

//go:noinline
func (request *scoringRequest16) HotScore() uint64 {
	acc := request.Core.AccountID
	acc += uint64(request.Core.Score)
	acc += uint64(request.Core.Attempts)
	acc += uint64(request.Core.Flags)
	return acc
}

//go:noinline
func scoreRequest16GenericExact[T *scoringRequest16](request T) uint64 {
	concrete := (*scoringRequest16)(request)
	acc := concrete.Core.AccountID
	acc += uint64(concrete.Core.Score)
	acc += uint64(concrete.Core.Attempts)
	acc += uint64(concrete.Core.Flags)
	return acc
}

//go:noinline
func scoreRequest24Concrete(request *scoringRequest24) uint64 {
	acc := request.Core.AccountID ^ request.Core.RequestID
	acc += uint64(request.Core.Score)
	acc += uint64(request.Core.Attempts)
	return acc
}

//go:noinline
func (request *scoringRequest24) HotScore() uint64 {
	acc := request.Core.AccountID ^ request.Core.RequestID
	acc += uint64(request.Core.Score)
	acc += uint64(request.Core.Attempts)
	return acc
}

//go:noinline
func scoreRequest24GenericExact[T *scoringRequest24](request T) uint64 {
	concrete := (*scoringRequest24)(request)
	acc := concrete.Core.AccountID ^ concrete.Core.RequestID
	acc += uint64(concrete.Core.Score)
	acc += uint64(concrete.Core.Attempts)
	return acc
}

//go:noinline
func scoreRequest32Concrete(request *scoringRequest32) uint64 {
	return scoreHotFields(request.Core, requestLatency32{})
}

//go:noinline
func scoreRequest32TinyConcrete(request *scoringRequest32) uint64 {
	return request.Core.AccountID
}

//go:noinline
func mutateRequest32Concrete(request *scoringRequest32) uint64 {
	request.Core.Score += 1
	request.Core.Attempts += 1
	request.Core.Flags ^= 1
	request.Core.RequestID += uint64(request.Core.Score)
	return scoreHotFields(request.Core, requestLatency32{})
}

//go:noinline
func mutateRequest32TinyConcrete(request *scoringRequest32) uint64 {
	request.Core.AccountID += 1
	return request.Core.AccountID
}

//go:noinline
func (request *scoringRequest32) HotScore() uint64 {
	return scoreHotFields(request.Core, requestLatency32{})
}

//go:noinline
func (request *scoringRequest32) TinyScore() uint64 {
	return request.Core.AccountID
}

//go:noinline
func (request *scoringRequest32) HotMutate() uint64 {
	request.Core.Score += 1
	request.Core.Attempts += 1
	request.Core.Flags ^= 1
	request.Core.RequestID += uint64(request.Core.Score)
	return scoreHotFields(request.Core, requestLatency32{})
}

//go:noinline
func (request *scoringRequest32) TinyMutate() uint64 {
	request.Core.AccountID += 1
	return request.Core.AccountID
}

//go:noinline
func scoreRequest32GenericExact[T *scoringRequest32](request T) uint64 {
	concrete := (*scoringRequest32)(request)
	return scoreHotFields(concrete.Core, requestLatency32{})
}

//go:noinline
func scoreRequest32TinyGenericExact[T *scoringRequest32](request T) uint64 {
	concrete := (*scoringRequest32)(request)
	return concrete.Core.AccountID
}

//go:noinline
func mutateRequest32GenericExact[T *scoringRequest32](request T) uint64 {
	concrete := (*scoringRequest32)(request)
	concrete.Core.Score += 1
	concrete.Core.Attempts += 1
	concrete.Core.Flags ^= 1
	concrete.Core.RequestID += uint64(concrete.Core.Score)
	return scoreHotFields(concrete.Core, requestLatency32{})
}

//go:noinline
func mutateRequest32TinyGenericExact[T *scoringRequest32](request T) uint64 {
	concrete := (*scoringRequest32)(request)
	concrete.Core.AccountID += 1
	return concrete.Core.AccountID
}

//go:noinline
func scoreRequest64Concrete(request *scoringRequest64) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func (request *scoringRequest64) HotScore() uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func scoreRequest64GenericExact[T *scoringRequest64](request T) uint64 {
	concrete := (*scoringRequest64)(request)
	return scoreHotFields(concrete.Core, concrete.Latency)
}

//go:noinline
func scoreRequest96Concrete(request *scoringRequest96) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func (request *scoringRequest96) HotScore() uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func scoreRequest96GenericExact[T *scoringRequest96](request T) uint64 {
	concrete := (*scoringRequest96)(request)
	return scoreHotFields(concrete.Core, concrete.Latency)
}

//go:noinline
func scoreRequest128Concrete(request *scoringRequest128) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func (request *scoringRequest128) HotScore() uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func scoreRequest128GenericExact[T *scoringRequest128](request T) uint64 {
	concrete := (*scoringRequest128)(request)
	return scoreHotFields(concrete.Core, concrete.Latency)
}

//go:noinline
func scoreRequest192Concrete(request *scoringRequest192) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func (request *scoringRequest192) HotScore() uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func scoreRequest192GenericExact[T *scoringRequest192](request T) uint64 {
	concrete := (*scoringRequest192)(request)
	return scoreHotFields(concrete.Core, concrete.Latency)
}

//go:noinline
func scoreRequest256Concrete(request *scoringRequest256) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func (request *scoringRequest256) HotScore() uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func scoreRequest256GenericExact[T *scoringRequest256](request T) uint64 {
	concrete := (*scoringRequest256)(request)
	return scoreHotFields(concrete.Core, concrete.Latency)
}

//go:noinline
func scoreRequest384Concrete(request *scoringRequest384) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func (request *scoringRequest384) HotScore() uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func scoreRequest384GenericExact[T *scoringRequest384](request T) uint64 {
	concrete := (*scoringRequest384)(request)
	return scoreHotFields(concrete.Core, concrete.Latency)
}

//go:noinline
func scoreRequest512Concrete(request *scoringRequest512) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func (request *scoringRequest512) HotScore() uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func scoreRequest512GenericExact[T *scoringRequest512](request T) uint64 {
	concrete := (*scoringRequest512)(request)
	return scoreHotFields(concrete.Core, concrete.Latency)
}
