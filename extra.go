package lazyq

import (
	"iter"
	"slices"
)

// These extra operations on queues are provided for special interactions.

// Clone creates a new queue by cloning the backing slices.
func Clone[E any](q Queue[E]) Queue[E] {
	return Queue[E]{
		maxIdx: q.maxIdx,
		elems:  slices.Clone(q.elems),
	}
}

// Payloads returns an iterator over Queue element payloads.
func Payloads[E any](q Queue[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := range Elements(q) {
			if !yield(e.E) {
				break
			}
		}
	}
}

// Elements returns an iterator over Queue element payloads.
func Elements[E any](q Queue[E]) iter.Seq[Elem[E]] {
	return func(yield func(Elem[E]) bool) {
		for _, e := range ElementIndices(q) {
			if !yield(e) {
				break
			}
		}
	}
}

// ElementIndices returns an iterator over Queue elements and their slices indices.
func ElementIndices[E any](q Queue[E]) iter.Seq2[int, Elem[E]] {
	return func(yield func(int, Elem[E]) bool) {
		n := q.Len()
		for i := q.maxIdx; i < n; i++ {
			if !yield(i, q.elems[i]) {
				return
			}
		}
		for i := range q.maxIdx {
			if !yield(i, q.elems[i]) {
				return
			}
		}
	}
}

// HasMaxElems returns true when q has at least one element in the max slice.
func HasMaxElems[E any](q Queue[E]) bool { return q.hasMaxElements() }

// MaxIndex returns the index of the max element or [Len].
func MaxIndex[E any](q Queue[E]) int { return q.maxIdx }

// First returns the first element (not the top element) ignoring any max slice elements.
func First[E any](q Queue[E]) E { return q.elems[0].E }

// First returns the first max element of q.
func FirstMaxElem[E any](q Queue[E]) E { return q.elems[q.maxIdx].E }

// At returns the element at i, possibly a max slice element.
func At[E any](q Queue[E], i int) E { return q.elems[i].E }

// ReplacePayload replaces the element payload at index i with e.
func ReplacePayload[E any](q Queue[E], i int, e E) { q.elems[i].E = e }

// Grow q to ensure space for another n elements without allocating.
func Grow[E any](q *Queue[E], n int) { q.elems = slices.Grow(q.elems, n) }
