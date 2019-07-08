package chessminimax

import (
	"reflect"
	"testing"
)

func TestNewSearcherHolder(test *testing.T) {
	var searcher MockBoundedMoveSearcher
	holder := newSearcherHolder(searcher)

	if !reflect.DeepEqual(
		holder.searcher,
		searcher,
	) {
		test.Fail()
	}
}

func TestSearcherHolderSetSearcher(
	test *testing.T,
) {
	var searcher MockBoundedMoveSearcher
	holder := newSearcherHolder(nil)
	holder.SetSearcher(searcher)

	if !reflect.DeepEqual(
		holder.searcher,
		searcher,
	) {
		test.Fail()
	}
}
