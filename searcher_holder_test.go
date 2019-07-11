package chessminimax

import (
	"reflect"
	"testing"
)

func TestSearcherHolderSetSearcher(
	test *testing.T,
) {
	var searcher MockMoveSearcher
	var holder searcherHolder
	holder.SetSearcher(searcher)

	if !reflect.DeepEqual(
		holder.searcher,
		searcher,
	) {
		test.Fail()
	}
}
