package chessminimax

import (
	"reflect"
	"testing"

	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
)

type MockMoveSearcher struct {
	setSearcher func(
		innerSearcher MoveSearcher,
	)
	setTerminator func(
		terminator terminators.SearchTerminator,
	)
	searchProgress func(deep int) float64
	searchMove     func(
		storage models.PieceStorage,
		color models.Color,
		deep int,
		bounds moves.Bounds,
	) (moves.ScoredMove, error)
}

func (
	searcher MockMoveSearcher,
) SetSearcher(
	innerSearcher MoveSearcher,
) {
	if searcher.setSearcher == nil {
		panic("not implemented")
	}

	searcher.setSearcher(innerSearcher)
}

func (
	searcher MockMoveSearcher,
) SetTerminator(
	terminator terminators.SearchTerminator,
) {
	if searcher.setTerminator == nil {
		panic("not implemented")
	}

	searcher.setTerminator(terminator)
}

func (
	searcher MockMoveSearcher,
) SearchProgress(deep int) float64 {
	if searcher.searchProgress == nil {
		panic("not implemented")
	}

	return searcher.searchProgress(deep)
}

func (
	searcher MockMoveSearcher,
) SearchMove(
	storage models.PieceStorage,
	color models.Color,
	deep int,
	bounds moves.Bounds,
) (moves.ScoredMove, error) {
	if searcher.searchMove == nil {
		panic("not implemented")
	}

	return searcher.searchMove(
		storage,
		color,
		deep,
		bounds,
	)
}

type MockSearchTerminator struct {
	isSearchTerminated func(deep int) bool
	searchProgress     func(deep int) float64
}

func (
	terminator MockSearchTerminator,
) IsSearchTerminated(deep int) bool {
	if terminator.isSearchTerminated == nil {
		panic("not implemented")
	}

	return terminator.isSearchTerminated(deep)
}

func (
	terminator MockSearchTerminator,
) SearchProgress(deep int) float64 {
	if terminator.searchProgress == nil {
		panic("not implemented")
	}

	return terminator.searchProgress(deep)
}

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
