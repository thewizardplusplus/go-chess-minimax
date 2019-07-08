package chessminimax

import (
	"reflect"
	"testing"
)

func TestSearcherHolderSetSearcher(
	test *testing.T,
) {
	var searcher MockBoundedMoveSearcher
	var holder searcherHolder
	holder.SetSearcher(searcher)

	if !reflect.DeepEqual(
		holder.searcher,
		searcher,
	) {
		test.Fail()
	}
}
