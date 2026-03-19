package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	stringBytePayloadSizes = []int{16, 64, 256, 1_024, 4_096, 16_384, 65_536}
	stringByteIterations   = []int{1, 10, 100, 1_000}
)

type stringByteConversionBench func(iterations int, payload []byte, asString string) uint64

var stringByteConversionSink uint64

func TestStringByteConversionOutputsMatch(t *testing.T) {
	testSizes := []int{16, 64, 256, 1_024}
	testIterations := []int{1, 3, 10}

	for _, size := range testSizes {
		payload := testingBytePayload(size, stableSeed("TestStringByteConversion", size))
		asString := string(payload)

		for _, iterations := range testIterations {
			t.Run(fmt.Sprintf("size=%d iterations=%d", size, iterations), func(t *testing.T) {
				bytesDirect := benchBytesDirect(iterations, payload, asString)
				stringDirect := benchStringDirect(iterations, payload, asString)
				stringToBytes := benchStringToBytes(iterations, payload, asString)
				bytesToString := benchBytesToString(iterations, payload, asString)

				require.Equal(t, bytesDirect, stringDirect)
				require.Equal(t, bytesDirect, stringToBytes)
				require.Equal(t, bytesDirect, bytesToString)
			})
		}
	}
}

func BenchmarkStringByteConversion(b *testing.B) {
	runBenchmark := func(runF stringByteConversionBench) {
		for _, size := range stringBytePayloadSizes {
			for _, iterations := range stringByteIterations {
				b.Run(
					fmt.Sprintf("size=%d iterations=%d", size, iterations),
					benchmarkStringByteConversion(size, iterations, runF),
				)
			}
		}
	}

	switch variant {
	case "bytes_direct":
		runBenchmark(benchBytesDirect)
	case "string_direct":
		runBenchmark(benchStringDirect)
	case "string_to_bytes":
		runBenchmark(benchStringToBytes)
	case "bytes_to_string":
		runBenchmark(benchBytesToString)
	default:
		b.Fatalf("specify which benchmark to run: -args -variant bytes_direct|string_direct|string_to_bytes|bytes_to_string")
	}
}

func benchmarkStringByteConversion(size int, iterations int, runF stringByteConversionBench) func(*testing.B) {
	payload := testingBytePayload(size, stableSeed("BenchmarkStringByteConversion", size, iterations))
	asString := string(payload)

	return func(b *testing.B) {
		for b.Loop() {
			stringByteConversionSink = runF(iterations, payload, asString)
		}
	}
}

func testingBytePayload(size int, seed uint64) []byte {
	rng := newDeterministicRand(seed)
	payload := make([]byte, size)
	for i := range payload {
		payload[i] = byte('a' + rng.Intn(26))
	}
	return payload
}

func benchBytesDirect(iterations int, payload []byte, _ string) uint64 {
	var acc uint64
	for range iterations {
		acc += consumeBytes(payload)
	}
	return acc
}

func benchStringDirect(iterations int, _ []byte, asString string) uint64 {
	var acc uint64
	for range iterations {
		acc += consumeString(asString)
	}
	return acc
}

func benchStringToBytes(iterations int, _ []byte, asString string) uint64 {
	var acc uint64
	for range iterations {
		acc += consumeBytes([]byte(asString))
	}
	return acc
}

func benchBytesToString(iterations int, payload []byte, _ string) uint64 {
	var acc uint64
	for range iterations {
		acc += consumeString(string(payload))
	}
	return acc
}

func consumeBytes(payload []byte) uint64 {
	var acc uint64
	for i := range payload {
		acc += uint64(payload[i])
	}
	return acc
}

func consumeString(payload string) uint64 {
	var acc uint64
	for i := 0; i < len(payload); i++ {
		acc += uint64(payload[i])
	}
	return acc
}
