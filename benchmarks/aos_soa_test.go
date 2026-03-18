package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var aosSoASizes = []int{10, 100, 1_000, 10_000, 100_000}

type aosSoABench func(size int) func(*testing.B)

type entityAoS struct {
	ID               uint64
	PosX, PosY, PosZ float64
	VelX, VelY, VelZ float64
	Health           int32
	Team             uint16
	Active           bool
	Name             string
	Material         string
	AIState          string
}

type entitiesSoA struct {
	ID               []uint64
	PosX, PosY, PosZ []float64
	VelX, VelY, VelZ []float64
	Health           []int32
	Team             []uint16
	Active           []bool
	Name             []string
	Material         []string
	AIState          []string
}

type entitySnapshot struct {
	ID               uint64
	PosX, PosY, PosZ float64
	VelX, VelY, VelZ float64
	Health           int32
	Team             uint16
	Name             string
	Material         string
	AIState          string
}

var (
	hotUpdateSink float64
	snapshotSink  []entitySnapshot
)

func TestAoSVsSoA(t *testing.T) {
	for _, size := range sizesReduced {
		aos, soa := buildEntityLayouts(size)

		hotAoS := copySlice(aos)
		hotSoA := cloneEntitiesSoA(soa)

		gotAoS := updatePositionsAoS(hotAoS)
		gotSoA := updatePositionsSoA(hotSoA)

		require.Equalf(t, gotAoS, gotSoA, "hot update accumulator mismatch for size=%d", size)
		requireEquivalentLayouts(t, size, hotAoS, hotSoA)
		require.Equalf(t, buildSnapshotsAoS(aos), buildSnapshotsSoA(soa), "snapshot mismatch for size=%d", size)
	}
}

func BenchmarkAoSVsSoA(b *testing.B) {
	runBenchmark := func(runF aosSoABench) {
		for _, size := range aosSoASizes {
			b.Run(
				fmt.Sprintf("size=%d", size),
				runF(size),
			)
		}
	}

	switch variant {
	case "aos_hot_update":
		runBenchmark(benchmarkAoSHotUpdate)
	case "soa_hot_update":
		runBenchmark(benchmarkSoAHotUpdate)
	case "aos_snapshot":
		runBenchmark(benchmarkAoSSnapshot)
	case "soa_snapshot":
		runBenchmark(benchmarkSoASnapshot)
	default:
		b.Errorf("specify which test to run: -args -variant aos_hot_update|soa_hot_update|aos_snapshot|soa_snapshot")
	}
}

func benchmarkAoSHotUpdate(size int) func(*testing.B) {
	aos, _ := buildEntityLayouts(size)
	return func(b *testing.B) {
		for b.Loop() {
			hotUpdateSink = updatePositionsAoS(aos)
		}
	}
}

func benchmarkSoAHotUpdate(size int) func(*testing.B) {
	_, soa := buildEntityLayouts(size)
	return func(b *testing.B) {
		for b.Loop() {
			hotUpdateSink = updatePositionsSoA(soa)
		}
	}
}

func benchmarkAoSSnapshot(size int) func(*testing.B) {
	aos, _ := buildEntityLayouts(size)
	return func(b *testing.B) {
		for b.Loop() {
			snapshotSink = buildSnapshotsAoS(aos)
		}
	}
}

func benchmarkSoASnapshot(size int) func(*testing.B) {
	_, soa := buildEntityLayouts(size)
	return func(b *testing.B) {
		for b.Loop() {
			snapshotSink = buildSnapshotsSoA(soa)
		}
	}
}

func buildEntityLayouts(size int) ([]entityAoS, entitiesSoA) {
	aos := make([]entityAoS, size)
	soa := entitiesSoA{
		ID:       make([]uint64, size),
		PosX:     make([]float64, size),
		PosY:     make([]float64, size),
		PosZ:     make([]float64, size),
		VelX:     make([]float64, size),
		VelY:     make([]float64, size),
		VelZ:     make([]float64, size),
		Health:   make([]int32, size),
		Team:     make([]uint16, size),
		Active:   make([]bool, size),
		Name:     make([]string, size),
		Material: make([]string, size),
		AIState:  make([]string, size),
	}

	names := []string{"grunt", "ranger", "medic", "tank"}
	materials := []string{"steel", "wood", "cloth"}
	aiStates := []string{"patrol", "hold", "flank", "rush"}

	for i := range aos {
		entity := entityAoS{
			ID:       uint64(i + 1),
			PosX:     float64(i) * 0.25,
			PosY:     float64(i) * 0.5,
			PosZ:     float64(i) * 0.75,
			VelX:     0.25 + float64(i%3)*0.1,
			VelY:     0.5 + float64(i%5)*0.05,
			VelZ:     0.75 + float64(i%7)*0.025,
			Health:   int32((i % 100) + 1),
			Team:     uint16(i % 8),
			Active:   i%3 != 0,
			Name:     names[i%len(names)],
			Material: materials[(i/2)%len(materials)],
			AIState:  aiStates[(i/3)%len(aiStates)],
		}
		aos[i] = entity

		soa.ID[i] = entity.ID
		soa.PosX[i] = entity.PosX
		soa.PosY[i] = entity.PosY
		soa.PosZ[i] = entity.PosZ
		soa.VelX[i] = entity.VelX
		soa.VelY[i] = entity.VelY
		soa.VelZ[i] = entity.VelZ
		soa.Health[i] = entity.Health
		soa.Team[i] = entity.Team
		soa.Active[i] = entity.Active
		soa.Name[i] = entity.Name
		soa.Material[i] = entity.Material
		soa.AIState[i] = entity.AIState
	}

	return aos, soa
}

func cloneEntitiesSoA(in entitiesSoA) entitiesSoA {
	return entitiesSoA{
		ID:       copySlice(in.ID),
		PosX:     copySlice(in.PosX),
		PosY:     copySlice(in.PosY),
		PosZ:     copySlice(in.PosZ),
		VelX:     copySlice(in.VelX),
		VelY:     copySlice(in.VelY),
		VelZ:     copySlice(in.VelZ),
		Health:   copySlice(in.Health),
		Team:     copySlice(in.Team),
		Active:   copySlice(in.Active),
		Name:     copySlice(in.Name),
		Material: copySlice(in.Material),
		AIState:  copySlice(in.AIState),
	}
}

func updatePositionsAoS(entities []entityAoS) float64 {
	var acc float64

	for i := range entities {
		if !entities[i].Active {
			continue
		}

		entities[i].PosX += entities[i].VelX
		entities[i].PosY += entities[i].VelY
		entities[i].PosZ += entities[i].VelZ

		acc += entities[i].PosX + entities[i].PosY + entities[i].PosZ
	}

	return acc
}

func updatePositionsSoA(entities entitiesSoA) float64 {
	var acc float64

	for i := range entities.ID {
		if !entities.Active[i] {
			continue
		}

		entities.PosX[i] += entities.VelX[i]
		entities.PosY[i] += entities.VelY[i]
		entities.PosZ[i] += entities.VelZ[i]

		acc += entities.PosX[i] + entities.PosY[i] + entities.PosZ[i]
	}

	return acc
}

func buildSnapshotsAoS(entities []entityAoS) []entitySnapshot {
	out := make([]entitySnapshot, 0, len(entities))

	for i := range entities {
		if !entities[i].Active {
			continue
		}

		out = append(out, entitySnapshot{
			ID:       entities[i].ID,
			PosX:     entities[i].PosX,
			PosY:     entities[i].PosY,
			PosZ:     entities[i].PosZ,
			VelX:     entities[i].VelX,
			VelY:     entities[i].VelY,
			VelZ:     entities[i].VelZ,
			Health:   entities[i].Health,
			Team:     entities[i].Team,
			Name:     entities[i].Name,
			Material: entities[i].Material,
			AIState:  entities[i].AIState,
		})
	}

	return out
}

func buildSnapshotsSoA(entities entitiesSoA) []entitySnapshot {
	out := make([]entitySnapshot, 0, len(entities.ID))

	for i := range entities.ID {
		if !entities.Active[i] {
			continue
		}

		out = append(out, entitySnapshot{
			ID:       entities.ID[i],
			PosX:     entities.PosX[i],
			PosY:     entities.PosY[i],
			PosZ:     entities.PosZ[i],
			VelX:     entities.VelX[i],
			VelY:     entities.VelY[i],
			VelZ:     entities.VelZ[i],
			Health:   entities.Health[i],
			Team:     entities.Team[i],
			Name:     entities.Name[i],
			Material: entities.Material[i],
			AIState:  entities.AIState[i],
		})
	}

	return out
}

func requireEquivalentLayouts(t *testing.T, size int, aos []entityAoS, soa entitiesSoA) {
	t.Helper()

	require.Len(t, aos, len(soa.ID))

	for i := range aos {
		require.Equalf(t, aos[i].ID, soa.ID[i], "size=%d index=%d id mismatch", size, i)
		require.Equalf(t, aos[i].PosX, soa.PosX[i], "size=%d index=%d posx mismatch", size, i)
		require.Equalf(t, aos[i].PosY, soa.PosY[i], "size=%d index=%d posy mismatch", size, i)
		require.Equalf(t, aos[i].PosZ, soa.PosZ[i], "size=%d index=%d posz mismatch", size, i)
		require.Equalf(t, aos[i].VelX, soa.VelX[i], "size=%d index=%d velx mismatch", size, i)
		require.Equalf(t, aos[i].VelY, soa.VelY[i], "size=%d index=%d vely mismatch", size, i)
		require.Equalf(t, aos[i].VelZ, soa.VelZ[i], "size=%d index=%d velz mismatch", size, i)
		require.Equalf(t, aos[i].Health, soa.Health[i], "size=%d index=%d health mismatch", size, i)
		require.Equalf(t, aos[i].Team, soa.Team[i], "size=%d index=%d team mismatch", size, i)
		require.Equalf(t, aos[i].Active, soa.Active[i], "size=%d index=%d active mismatch", size, i)
		require.Equalf(t, aos[i].Name, soa.Name[i], "size=%d index=%d name mismatch", size, i)
		require.Equalf(t, aos[i].Material, soa.Material[i], "size=%d index=%d material mismatch", size, i)
		require.Equalf(t, aos[i].AIState, soa.AIState[i], "size=%d index=%d ai state mismatch", size, i)
	}
}
