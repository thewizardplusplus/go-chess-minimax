package chessminimax

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-chess-models"
)

type MockSearchTerminator struct{}

func (
	terminator MockSearchTerminator,
) IsSearchTerminate(deep int) bool {
	panic("not implemented")
}

type MockBoardEvaluator struct{}

func (
	evaluator MockBoardEvaluator,
) EvaluateBoard(
	storage models.PieceStorage,
	color models.Color,
) float64 {
	panic("not implemented")
}

type MockMoveGenerator struct{}

func (
	generator MockMoveGenerator,
) MovesForColor(
	storage models.PieceStorage,
	color models.Color,
) []models.Move {
	panic("not implemented")
}

func TestMoveGeneratorInterface(
	test *testing.T,
) {
	gotType := reflect.TypeOf(
		models.MoveGenerator{},
	)
	wantType := reflect.
		TypeOf((*MoveGenerator)(nil)).
		Elem()
	if !gotType.Implements(wantType) {
		test.Fail()
	}
}

func TestMoveSearcherInterface(
	test *testing.T,
) {
	gotType := reflect.TypeOf(
		DefaultMoveSearcher{},
	)
	wantType := reflect.
		TypeOf((*MoveSearcher)(nil)).
		Elem()
	if !gotType.Implements(wantType) {
		test.Fail()
	}
}

func TestNewDefaultMoveSearcher(
	test *testing.T,
) {
	var terminator MockSearchTerminator
	var evaluator MockBoardEvaluator
	var generator MockMoveGenerator
	searcher := NewDefaultMoveSearcher(
		terminator,
		evaluator,
		generator,
	)

	if !reflect.DeepEqual(
		searcher.terminator,
		terminator,
	) {
		test.Fail()
	}
	if !reflect.DeepEqual(
		searcher.evaluator,
		evaluator,
	) {
		test.Fail()
	}
	if !reflect.DeepEqual(
		searcher.generator,
		generator,
	) {
		test.Fail()
	}

	// check a reference to itself
	if !reflect.DeepEqual(
		searcher.searcher,
		searcher,
	) {
		test.Fail()
	}
}
