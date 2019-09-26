package chessminimax

import (
	"reflect"
	"testing"
)

func TestBaseSearcherSetSearcher(
	test *testing.T,
) {
	var innerSearcher MockMoveSearcher
	var searcher BaseSearcher
	searcher.SetSearcher(innerSearcher)

	if !reflect.DeepEqual(
		searcher.searcher,
		innerSearcher,
	) {
		test.Fail()
	}
}

func TestBaseSearcherSetTerminator(
	test *testing.T,
) {
	var terminator MockSearchTerminator
	var searcher BaseSearcher
	searcher.SetTerminator(terminator)

	if !reflect.DeepEqual(
		searcher.terminator,
		terminator,
	) {
		test.Fail()
	}
}
