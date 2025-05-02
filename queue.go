package lazyq

import (
	"math"
)

var maxPriority = float32(math.Inf(+1))

// Elem is a entry in the queue with a float32 prioirty and a generic payload.
type Elem[E any] struct {
	E        E
	Priority float32
}

// Queue implements a monotomic binary max-heap task queue with an O(1) push operation called AppendMax.
//
// The heap is represented efficiently as a single slice with a heapified part, followed by a unformatted slice
// of max-priority elements.
//
// Queue is not a general purpose queue. It only serves use cases where
// items are always added with the max priority, and the task list grows
// monotomically in size, and only the top element's priority can be changed.
//
// See the additional methods in this package to support out of band updates
// to the queue.
//
// Example:
//
//	1 2 3 4 5 6 7 | +inf +inf +inf...
//
// Operations:
//
//	AppendMax: Append an item with priority +inf in O(1).
//	Next:      Return the highest priority item in O(1).
//	Decrease:  Decrease the priority of the next element in O(logN).
type Queue[E any] struct {
	maxIdx int // Index of the first max element, or q.Len().
	elems  []Elem[E]
}

// AppendMax appends a new element with the max priority in O(1).
func (q *Queue[E]) AppendMax(e E) { q.elems = append(q.elems, Elem[E]{e, maxPriority}) }

// Next returns the maximum element.
//
// Next heapifies one of the max elements if needed.
//
// Next panics when q.Len() == 0.
func (q *Queue[E]) Next() E {
	// heapify a max element if available.
	if q.maxIdx == 0 || q.elems[0].Priority < maxPriority && q.hasMaxElements() {
		if q.maxIdx > 0 {
			q.up(q.maxIdx)
		}
		q.maxIdx++
	}
	return q.elems[0].E
}

// Decrease decreases the priority of the next element by setting it to p.
//
// Decrease panics if q.Len() == 0.
func (q *Queue[E]) Decrease(p float32) {
	if math.IsNaN(float64(p)) {
		return
	}
	q.elems[0].Priority = p
	if q.maxIdx == 0 {
		// This happens when [AppendMax] is directly followed by [Decrease]
		// from an empty Queue. Here we just do what [Next] would have:
		// increment maxIdx and return. No need to call down0 on one element.
		q.maxIdx++
		return
	}
	q.down0()
}

// Len returns the length of the queue.
func (q Queue[E]) Len() int { return len(q.elems) }

func (q *Queue[E]) hasMaxElements() bool { return q.maxIdx < q.Len() }

// Code for up and down is borrowed from go/src/heap.go.

func (q Queue[E]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !q.higherPriority(j, i) {
			break
		}
		q.swap(i, j)
		j = i
	}
}

func (q Queue[E]) down0() {
	i := 0
	n := q.maxIdx // The heap ends at maxIdx-1.
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

func (q Queue[E]) swap(i, j int) { q.elems[i], q.elems[j] = q.elems[j], q.elems[i] }

// Higher priority nodes first.
func (q Queue[E]) higherPriority(i, j int) bool { return q.elems[i].Priority > q.elems[j].Priority }
