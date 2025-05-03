package lazyq

import (
	"math"
	"math/rand/v2"
	"testing"
)

type baselineQueue[E any] []Elem[E]

func (q *baselineQueue[E]) AppendMax(e E) {
	n := q.Len()
	*q = append(*q, Elem[E]{e, maxPriority})
	q.up(n)
}

func (q baselineQueue[E]) Next() E { return q[0].E }

func (q baselineQueue[E]) Decrease(p float32) {
	if math.IsNaN(float64(p)) {
		return
	}
	q[0].Priority = p
	if q.Len() > 1 {
		q.down0()
	}
}

func (q baselineQueue[E]) Len() int { return len(q) }

// Code for up and down is borrowed from go/src/heap.go.

func (q baselineQueue[E]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !q.higherPriority(j, i) {
			break
		}
		q.swap(i, j)
		j = i
	}
}

func (q baselineQueue[E]) down0() {
	i := 0
	n := q.Len()
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && q.higherPriority(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !q.higherPriority(j, i) {
			break
		}
		q.swap(i, j)
		i = j
	}
}

func (q baselineQueue[E]) swap(i, j int) { q[i], q[j] = q[j], q[i] }

func (q baselineQueue[E]) higherPriority(i, j int) bool { return q[i].Priority > q[j].Priority }

type queueInterface[E any] interface {
	Len() int
	AppendMax(E)
	Decrease(float32)
	Next() E
}

// readFactor determines the expected relative rate of element reads to writes in benchmarks.
const readFactor = 1.0 / 100

const meanOps = 1000

func randomReadsAndWrites[E any](tb testing.TB, r *rand.Rand, q queueInterface[E]) {
	tb.Helper()
	nw := int(r.ExpFloat64() * meanOps)
	var zero E
	for range nw {
		q.AppendMax(zero)
	}
	nr := int(r.ExpFloat64() * meanOps * readFactor)
	for range nr {
		q.Next()
		q.Decrease(randomPriority(tb, r))
	}
}

func BenchmarkBaselineQueue(b *testing.B) {
	r := rand.New(rand.NewPCG(1337, 420))
	for b.Loop() {
		var q baselineQueue[struct{}]
		q.AppendMax(struct{}{})
		randomReadsAndWrites(b, r, &q)
	}
}

func BenchmarkLazyQueue(b *testing.B) {
	r := rand.New(rand.NewPCG(1337, 420))
	for b.Loop() {
		var q Queue[struct{}]
		q.AppendMax(struct{}{})
		randomReadsAndWrites(b, r, &q)
	}
}
