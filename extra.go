package lazyq

import (
	"iter"
	"slices"
)

// These extra operations on queues are provided for special interactions.
// TODO: Add extra mutation methods here.

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
