package chessminimax

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-chess-models"
)

type MockSearchTerminator struct {
	isSearchTerminate func(deep int) bool
}

func (
	terminator MockSearchTerminator,
) IsSearchTerminate(deep int) bool {
	if terminator.isSearchTerminate == nil {
		panic("not implemented")
	}

	return terminator.isSearchTerminate(deep)
}

type MockBoardEvaluator struct {
	evaluateBoard func(
		storage models.PieceStorage,
		color models.Color,
	) float64
}

func (
	evaluator MockBoardEvaluator,
) EvaluateBoard(
	storage models.PieceStorage,
	color models.Color,
) float64 {
	if evaluator.evaluateBoard == nil {
		panic("not implemented")
	}

	return evaluator.evaluateBoard(
		storage,
		color,
	)
}

type MockMoveGenerator struct {
	movesForColor func(
		storage models.PieceStorage,
		color models.Color,
	) []models.Move
}

func (
	generator MockMoveGenerator,
) MovesForColor(
	storage models.PieceStorage,
	color models.Color,
) []models.Move {
	if generator.movesForColor == nil {
		panic("not implemented")
	}

	return generator.movesForColor(
		storage,
		color,
	)
}

type MockMoveSearcher struct {
	searchMove func(
		storage models.PieceStorage,
		color models.Color,
		deep int,
	) (ScoredMove, error)
}

func (searcher MockMoveSearcher) SearchMove(
	storage models.PieceStorage,
	color models.Color,
	deep int,
) (ScoredMove, error) {
	if searcher.searchMove == nil {
		panic("not implemented")
	}

	return searcher.searchMove(
		storage,
		color,
		deep,
	)
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

func TestDefaultMoveSearcherSearchMove(
	test *testing.T,
) {
	type fields struct {
		terminator SearchTerminator
		evaluator  BoardEvaluator
		generator  MoveGenerator
		searcher   MoveSearcher
	}
	type args struct {
		storage models.PieceStorage
		color   models.Color
		deep    int
	}
	type data struct {
		fields   fields
		args     args
		wantMove ScoredMove
		wantErr  error
	}

	for _, data := range []data{} {
		searcher := DefaultMoveSearcher{
			terminator: data.fields.terminator,
			evaluator:  data.fields.evaluator,
			generator:  data.fields.generator,
			searcher:   data.fields.searcher,
		}
		gotMove, gotErr := searcher.SearchMove(
			data.args.storage,
			data.args.color,
			data.args.deep,
		)

		if !reflect.DeepEqual(
			gotMove,
			data.wantMove,
		) {
			test.Fail()
		}
		if !reflect.DeepEqual(
			gotErr,
			data.wantErr,
		) {
			test.Fail()
		}
	}
}
