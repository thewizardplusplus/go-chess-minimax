package chessminimax

import (
	"reflect"
	"testing"
)

func TestBaseSearcherSetSearcher(
	test *testing.T,
) {
	var innerSearcher MockMoveSearcher
	var searcher baseSearcher
	searcher.SetSearcher(innerSearcher)

	if !reflect.DeepEqual(
		searcher.searcher,
		innerSearcher,
	) {
		test.Fail()
	}
}
