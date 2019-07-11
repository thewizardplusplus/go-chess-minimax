package models

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-chess-models"
)

func TestNewScoredMove(test *testing.T) {
	got := NewScoredMove()

	want := ScoredMove{Score: initialScore}
	if !reflect.DeepEqual(got, want) {
		test.Fail()
	}
}

func TestScoredMoveIsUpdated(
	test *testing.T,
) {
	type fields struct {
		score float64
	}
	type data struct {
		fields fields
		want   bool
	}

	for _, data := range []data{
		data{
			fields: fields{initialScore},
			want:   false,
		},
		data{
			fields: fields{2.3},
			want:   true,
		},
	} {
		move := ScoredMove{
			Score: data.fields.score,
		}
		got := move.IsUpdated()

		if got != data.want {
			test.Fail()
		}
	}
}

func TestScoredMoveUpdate(test *testing.T) {
	type fields struct {
		move  models.Move
		score float64
	}
	type args struct {
		scoredMove ScoredMove
		rawMove    models.Move
	}
	type data struct {
		fields   fields
		args     args
		wantMove ScoredMove
	}

	for _, data := range []data{
		data{
			fields: fields{
				move: models.Move{
					Start: models.Position{
						File: 1,
						Rank: 2,
					},
					Finish: models.Position{
						File: 3,
						Rank: 4,
					},
				},
				score: 4.2,
			},
			args: args{
				scoredMove: ScoredMove{Score: 2.3},
				rawMove: models.Move{
					Start: models.Position{
						File: 5,
						Rank: 6,
					},
					Finish: models.Position{
						File: 7,
						Rank: 8,
					},
				},
			},
			wantMove: ScoredMove{
				Move: models.Move{
					Start: models.Position{
						File: 1,
						Rank: 2,
					},
					Finish: models.Position{
						File: 3,
						Rank: 4,
					},
				},
				Score: 4.2,
			},
		},
		data{
			fields: fields{
				move: models.Move{
					Start: models.Position{
						File: 1,
						Rank: 2,
					},
					Finish: models.Position{
						File: 3,
						Rank: 4,
					},
				},
				score: -4.2,
			},
			args: args{
				scoredMove: ScoredMove{Score: 2.3},
				rawMove: models.Move{
					Start: models.Position{
						File: 5,
						Rank: 6,
					},
					Finish: models.Position{
						File: 7,
						Rank: 8,
					},
				},
			},
			wantMove: ScoredMove{
				Move: models.Move{
					Start: models.Position{
						File: 5,
						Rank: 6,
					},
					Finish: models.Position{
						File: 7,
						Rank: 8,
					},
				},
				Score: -2.3,
			},
		},
	} {
		move := ScoredMove{
			Move:  data.fields.move,
			Score: data.fields.score,
		}
		move.Update(
			data.args.scoredMove,
			data.args.rawMove,
		)

		if !reflect.DeepEqual(
			move,
			data.wantMove,
		) {
			test.Fail()
		}
	}
}
