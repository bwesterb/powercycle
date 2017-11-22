package powercycle_test

import (
	"github.com/bwesterb/powercycle"
	"testing"
)

func TestOneElementCycle(t *testing.T) {
	cycle := powercycle.New(1)
	if cycle.Apply(0) != 0 {
		t.Fatal("One-element cycle does not map 0 to 0.")
	}
}

func testCycle(n uint64, t *testing.T) {
	var x, i uint64
	cycle := powercycle.New(n)
	seen := make([]bool, n)
	for i = 0; i < n; i++ {
		x = cycle.Apply(x)
		if seen[x] {
			t.Fatalf("%d appears twice in cycle %s", x, cycle)
		}
		seen[x] = true
	}
}

func testHugeCycle(n, m uint64, t *testing.T) {
	var x, i uint64
	cycle := powercycle.New(n)
	seen := make(map[uint64]bool)
	for i = 0; i < m; i++ {
		x = cycle.Apply(x)
		if x >= n {
			t.Fatalf("%d is too big for cycle %s", x, cycle)
		}
		if _, oops := seen[x]; oops {
			t.Fatalf("%d appears twice in cycle %s", x, cycle)
		}
		seen[x] = true
	}
}

func TestSmallCycles(t *testing.T) {
	var n uint64
	for n = 2; n < 100; n++ {
		testCycle(n, t)
	}
}

func TestMediumSizedCycle(t *testing.T) {
	testCycle(1000, t)
	testCycle(5000, t)
	testCycle(10000, t)
	testCycle(50000, t)
	testCycle(100000, t)
	testCycle(500000, t)
}

func TestBigCycle(t *testing.T) {
	testCycle(10000000, t)
}

func TestHugeCycles(t *testing.T) {
	testHugeCycle(1000000000, 1000000, t)
	testHugeCycle(100000000000, 1000000, t)
}

func benchmarkApply(n uint64, b *testing.B) {
	cycle := powercycle.New(n)
	b.ResetTimer()
	var x uint64
	for i := 0; i < b.N; i++ {
		x = cycle.Apply(x)
	}
}

func benchmarkNew(n uint64, b *testing.B) {
	for i := 0; i < b.N; i++ {
		powercycle.New(n)
	}
}

func BenchmarkNew10(b *testing.B) {
	benchmarkNew(10, b)
}
func BenchmarkNew1000(b *testing.B) {
	benchmarkNew(1000, b)
}
func BenchmarkNew1000000(b *testing.B) {
	benchmarkNew(1000000, b)
}
func BenchmarkNew1000000000(b *testing.B) {
	benchmarkNew(1000000000, b)
}
func BenchmarkNew1000000000000(b *testing.B) {
	benchmarkNew(1000000000000, b)
}

func BenchmarkApply10(b *testing.B) {
	benchmarkApply(10, b)
}
func BenchmarkApply1000(b *testing.B) {
	benchmarkApply(1000, b)
}
func BenchmarkApply1000000(b *testing.B) {
	benchmarkApply(1000000, b)
}
func BenchmarkApply1000000000(b *testing.B) {
	benchmarkApply(1000000000, b)
}
func BenchmarkApply1000000000000(b *testing.B) {
	benchmarkApply(1000000000000, b)
}
func BenchmarkApply1000000000000000(b *testing.B) {
	benchmarkApply(1000000000000000, b)
}
