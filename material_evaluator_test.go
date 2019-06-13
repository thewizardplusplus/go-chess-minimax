package chessminimax

import (
	"testing"

	models "github.com/thewizardplusplus/go-chess-models"
)

func TestMaterialEvaluatorEvaluateBoard(
	test *testing.T,
) {
	type args struct {
		storage models.PieceStorage
		color   models.Color
	}
	type data struct {
		args args
		want float64
	}

	for _, data := range []data{} {
		var evaluator MaterialEvaluator
		got := evaluator.EvaluateBoard(
			data.args.storage,
			data.args.color,
		)

		if got != data.want {
			test.Fail()
		}
	}
}

func TestColorSign(test *testing.T) {
	type args struct {
		piece models.Piece
		color models.Color
	}
	type data struct {
		args args
		want float64
	}

	for _, data := range []data{} {
		got := colorSign(
			data.args.piece,
			data.args.color,
		)

		if got != data.want {
			test.Fail()
		}
	}
}
