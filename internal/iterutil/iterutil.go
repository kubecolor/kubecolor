package iterutil

import "iter"

type NextPair[T any] struct {
	Curr    T
	Next    T
	HasNext bool
}

type NextPair2[K, V any] struct {
	Curr1   V
	Next1   V
	Curr0   K
	Next0   K
	HasNext bool
}

func WindowedPair[T any](seq iter.Seq[T]) iter.Seq[NextPair[T]] {
	return func(yield func(NextPair[T]) bool) {
		var pair NextPair[T]
		for val := range seq {
			if !pair.HasNext {
				pair.Next = val
				pair.HasNext = true
				// is yielded in the next iter
				continue
			}

			pair.Curr = pair.Next
			pair.Next = val
			if !yield(pair) {
				return
			}
		}

		if pair.HasNext {
			// must check "pair.HasNext" in case sequence was empty
			pair.Curr = pair.Next
			var zero T
			pair.Next = zero
			pair.HasNext = false
			if !yield(pair) {
				return
			}
		}
	}
}

func WindowedPair2[K, V any](seq iter.Seq2[K, V]) iter.Seq[NextPair2[K, V]] {
	return func(yield func(NextPair2[K, V]) bool) {
		var pair NextPair2[K, V]
		for val0, val1 := range seq {
			if !pair.HasNext {
				pair.Next0 = val0
				pair.Next1 = val1
				pair.HasNext = true
				// is yielded in the next iter
				continue
			}

			pair.Curr0 = pair.Next0
			pair.Curr1 = pair.Next1
			pair.Next0 = val0
			pair.Next1 = val1
			if !yield(pair) {
				return
			}
		}

		if pair.HasNext {
			// must check "pair.HasNext" in case sequence was empty
			pair.Curr0 = pair.Next0
			pair.Curr1 = pair.Next1
			var zero0 K
			var zero1 V
			pair.Next0 = zero0
			pair.Next1 = zero1
			pair.HasNext = false
			if !yield(pair) {
				return
			}
		}
	}
}
