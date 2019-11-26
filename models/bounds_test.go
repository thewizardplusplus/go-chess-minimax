package models

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
	got := Bounds{2.3, 4.2}.Next()

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
		scoredMove  ScoredMove
		rawMove     models.Move
		moveQuality float64
	}
	type data struct {
		fields     fields
		args       args
		wantBounds Bounds
		wantMove   ScoredMove
		wantOk     bool
	}

	for _, data := range []data{
		{
			fields: fields{-2.3, 4.2},
			args: args{
				scoredMove: ScoredMove{
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
					Score:   5,
					Quality: 0.25,
				},
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
				moveQuality: 0.75,
			},
			wantBounds: Bounds{-2.3, 4.2},
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
				Score:   5,
				Quality: 0.25,
			},
			wantOk: true,
		},
		{
			fields: fields{-2.3, 4.2},
			args: args{
				scoredMove: ScoredMove{
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
					Score:   1.2,
					Quality: 0.25,
				},
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
				moveQuality: 0.75,
			},
			wantBounds: Bounds{-1.2, 4.2},
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
				Score:   1.2,
				Quality: 0.25,
			},
			wantOk: true,
		},
		{
			fields: fields{-2.3, 4.2},
			args: args{
				scoredMove: ScoredMove{
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
					Score:   -4.2,
					Quality: 0.25,
				},
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
				moveQuality: 0.75,
			},
			wantBounds: Bounds{4.2, 4.2},
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
				Score:   4.2,
				Quality: 0.75,
			},
			wantOk: false,
		},
		{
			fields: fields{-2.3, 4.2},
			args: args{
				scoredMove: ScoredMove{
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
					Score:   -5,
					Quality: 0.25,
				},
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
				moveQuality: 0.75,
			},
			wantBounds: Bounds{5, 4.2},
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
				Score:   5,
				Quality: 0.75,
			},
			wantOk: false,
		},
	} {
		bounds := Bounds{
			Alpha: data.fields.alpha,
			Beta:  data.fields.beta,
		}
		gotMove, gotOk :=
			bounds.Update(data.args.scoredMove, data.args.rawMove, data.args.moveQuality)

		if !reflect.DeepEqual(bounds, data.wantBounds) {
			test.Fail()
		}
		if !reflect.DeepEqual(gotMove, data.wantMove) {
			test.Fail()
		}
		if gotOk != data.wantOk {
			test.Fail()
		}
	}
}
