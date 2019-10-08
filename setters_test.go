package chessminimax

import (
	"reflect"
	"testing"
)

func TestSearcherSetterSetSearcher(
	test *testing.T,
) {
	var searcher MockMoveSearcher
	var setter SearcherSetter
	setter.SetSearcher(searcher)

	if !reflect.DeepEqual(
		setter.searcher,
		searcher,
	) {
		test.Fail()
	}
}

func TestTerminatorSetterSetTerminator(
	test *testing.T,
) {
	var terminator MockSearchTerminator
	var setter TerminatorSetter
	setter.SetTerminator(terminator)

	if !reflect.DeepEqual(
		setter.terminator,
		terminator,
	) {
		test.Fail()
	}
}
