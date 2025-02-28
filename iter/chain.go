package iter

import "github.com/BooleanCat/go-functional/option"

// ChainIter implements `Chain`. See `Chain`'s documentation.
type ChainIter[T any] struct {
	iterators     []Iterator[T]
	iteratorIndex int
}

// Chain instantiates a `ChainIter` that will yield all items in the provided iterators to exhaustion, from left to right.
func Chain[T any](iterators ...Iterator[T]) *ChainIter[T] {
	return &ChainIter[T]{iterators, 0}
}

// Next implements the Iterator interface for `Chain`.
func (iter *ChainIter[T]) Next() option.Option[T] {
	for {
		// If the index is equal to the number of iterators then we have exhausted them all.
		if iter.iteratorIndex == len(iter.iterators) {
			return option.None[T]()
		}

		// Otherwise get the currently active iterator and ask it for a value
		currentIterator := iter.iterators[iter.iteratorIndex]
		value, ok := currentIterator.Next().Value()

		// If there is a value then emit it
		if ok {
			return option.Some(value)
		}

		// Otherwise the iterator has been exhausted, so increase the iterator index for the next
		// iteration of the loop.
		iter.iteratorIndex = iter.iteratorIndex + 1
	}

}

var _ Iterator[struct{}] = new(ChainIter[struct{}])
