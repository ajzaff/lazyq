package lazyq

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
)

func checkHeapProperties[E any](t testing.TB, testName string, h []Elem[E]) {
	n := len(h)
	for i, e := range h {
		left := 2*i + 1
		right := 2*i + 2

		if i > 0 && e.Priority == maxPriority {
			t.Errorf("%s: lazy queue property violated at %d: max priority found beyond 0th element", testName, i)
		}
		if left < n {
			e1 := h[left]
			if e.Priority < e1.Priority {
				t.Errorf("%s: heap property violated at (%d, %d): e < left", testName, i, left)
			}
		}
		if right < n {
			e1 := h[right]
			if e.Priority < e1.Priority {
				t.Errorf("%s: heap property violated at (%d, %d): e < right", testName, i, right)
			}
		}

		if right >= n {
			break
		}
	}
}

func checkQueueConditions[E any](t testing.TB, testName string, q Queue[E]) {
	t.Helper()

	defer func() {
		if err := recover(); err != nil {
			t.Errorf("%s: recovered from panic while asserting queue conditions: %v", testName, err)
		}
	}()

	n := q.Len()

	if n < q.maxIdx {
		t.Errorf("%s: lazy queue property violated (%d < %d): q.Len() < maxIdx", testName, n, q.maxIdx)
	}

	checkHeapProperties(t, testName, q.elems[:q.maxIdx])

	// Check unsorted slice.
	for i := q.maxIdx; i < n; i++ {
		e := q.elems[i]
		if e.Priority != maxPriority {
			t.Errorf("%s: lazy queue property violated at %d: e.Priority != +inf", testName, i)
		}
	}

	// Check all elements again for instances of NaN.
	for i, e := range q.elems {
		if math.IsNaN(float64(e.Priority)) {
			t.Errorf("%s: lazy queue property violated at %d: found NaN priority", testName, i)
		}
	}
}

func TestZero(t *testing.T) {
	var q Queue[struct{}]

	checkQueueConditions(t, "TestZero()", q)
}

func TestZeroInsertMax(t *testing.T) {
	var q Queue[struct{}]

	q.AppendMax(struct{}{})
	q.AppendMax(struct{}{})
	q.Next()
	q.Decrease(7)
	q.Next()

	checkQueueConditions(t, "TestZero()", q)
}

func TestZeroInsertNaN(t *testing.T) {
	var q Queue[struct{}]

	q.AppendMax(struct{}{})
	q.Next()
	q.Decrease(float32(math.NaN()))

	checkQueueConditions(t, "TestZero()", q)
}

func TestDecreaseMaxPriorityHasNoImpact(t *testing.T) {
	var q Queue[struct{}]

	q.AppendMax(struct{}{})
	q.Decrease(maxPriority)

	checkQueueConditions(t, "TestDecreaseMaxPriorityHasNoImpact()", q)
}

func TestFuzzQueue838c3636634247da(t *testing.T) {
	fuzzFunc(t, "aa0nna%n", 1337, 420)
}

func FuzzQueue(f *testing.F) {
	f.Add("", uint64(1337), uint64(420))  // No op
	f.Add("a", uint64(1337), uint64(420)) // AppendMax(struct{})
	f.Add("d", uint64(1337), uint64(420)) // Decrease(randomPriority())
	f.Add("n", uint64(1337), uint64(420)) // Next()
	f.Fuzz(fuzzFunc)
}

func randomPriority(t *testing.T, r *rand.Rand) float32 {
	t.Helper()
	switch r.IntN(7) {
	case 0:
		return -maxPriority
	case 1:
		return -0
	case 2:
		return +0
	case 3:
		return 2*r.Float32() - 1
	case 4:
		return math.Float32frombits(r.Uint32())
	case 5:
		return float32(math.NaN())
	default: // 6
		return maxPriority
	}
}

func fuzzFunc(t *testing.T, cmds string, seed1, seed2 uint64) {
	r := rand.New(rand.NewPCG(seed1, seed2))
	testName := fmt.Sprintf("FuzzQueue(%q)", string(cmds))
	var q Queue[struct{}]
	for _, cmd := range []byte(cmds) {
		if cmd == 'a' {
			q.AppendMax(struct{}{})
			continue
		}
		if q.Len() == 0 {
			continue
		}
		switch cmd {
		case 'd':
			q.Decrease(randomPriority(t, r))
		case 'n':
			q.Next()
		}
	}
	checkQueueConditions(t, testName, q)
}
