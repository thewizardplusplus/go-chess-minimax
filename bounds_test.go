package chessminimax

import (
	"math"
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-chess-models"
)

func TestNewBounds(test *testing.T) {
	got := NewBounds()

	want := Bounds{math.Inf(-1), math.Inf(+1)}
	if !reflect.DeepEqual(got, want) {
		test.Fail()
	}
}

func TestBoundsNext(test *testing.T) {
	got := Bounds{2.3, 4.2}.next()

	want := Bounds{-4.2, -2.3}
	if !reflect.DeepEqual(got, want) {
		test.Fail()
	}
}

func TestBoundsUpdate(test *testing.T) {
	type fields struct {
		alpha float64
		beta  float64
	}
	type args struct {
		scoredMove ScoredMove
		rawMove    models.Move
	}
	type data struct {
		fields     fields
		args       args
		wantBounds Bounds
		wantMove   ScoredMove
		wantOk     bool
	}

	for _, data := range []data{
		data{
			fields: fields{-2.3, 4.2},
			args: args{
				scoredMove: ScoredMove{
					Move: models.Move{
						Start:  models.Position{1, 2},
						Finish: models.Position{3, 4},
					},
					Score: 5,
				},
				rawMove: models.Move{
					Start:  models.Position{5, 6},
					Finish: models.Position{7, 8},
				},
			},
			wantBounds: Bounds{-2.3, 4.2},
			wantMove: ScoredMove{
				Move: models.Move{
					Start:  models.Position{1, 2},
					Finish: models.Position{3, 4},
				},
				Score: 5,
			},
			wantOk: true,
		},
		data{
			fields: fields{-2.3, 4.2},
			args: args{
				scoredMove: ScoredMove{
					Move: models.Move{
						Start:  models.Position{1, 2},
						Finish: models.Position{3, 4},
					},
					Score: 1.2,
				},
				rawMove: models.Move{
					Start:  models.Position{5, 6},
					Finish: models.Position{7, 8},
				},
			},
			wantBounds: Bounds{-1.2, 4.2},
			wantMove: ScoredMove{
				Move: models.Move{
					Start:  models.Position{1, 2},
					Finish: models.Position{3, 4},
				},
				Score: 1.2,
			},
			wantOk: true,
		},
		data{
			fields: fields{-2.3, 4.2},
			args: args{
				scoredMove: ScoredMove{
					Move: models.Move{
						Start:  models.Position{1, 2},
						Finish: models.Position{3, 4},
					},
					Score: -4.2,
				},
				rawMove: models.Move{
					Start:  models.Position{5, 6},
					Finish: models.Position{7, 8},
				},
			},
			wantBounds: Bounds{4.2, 4.2},
			wantMove: ScoredMove{
				Move: models.Move{
					Start:  models.Position{5, 6},
					Finish: models.Position{7, 8},
				},
				Score: 4.2,
			},
			wantOk: false,
		},
		data{
			fields: fields{-2.3, 4.2},
			args: args{
				scoredMove: ScoredMove{
					Move: models.Move{
						Start:  models.Position{1, 2},
						Finish: models.Position{3, 4},
					},
					Score: -5,
				},
				rawMove: models.Move{
					Start:  models.Position{5, 6},
					Finish: models.Position{7, 8},
				},
			},
			wantBounds: Bounds{5, 4.2},
			wantMove: ScoredMove{
				Move: models.Move{
					Start:  models.Position{5, 6},
					Finish: models.Position{7, 8},
				},
				Score: 5,
			},
			wantOk: false,
		},
	} {
		bounds := Bounds{
			Alpha: data.fields.alpha,
			Beta:  data.fields.beta,
		}
		gotMove, gotOk := bounds.update(
			data.args.scoredMove,
			data.args.rawMove,
		)

		if !reflect.DeepEqual(
			bounds,
			data.wantBounds,
		) {
			test.Fail()
		}
		if !reflect.DeepEqual(
			gotMove,
			data.wantMove,
		) {
			test.Fail()
		}
		if gotOk != data.wantOk {
			test.Fail()
		}
	}
}
