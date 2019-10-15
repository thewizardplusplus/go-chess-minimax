package chessminimax

import (
	"errors"
	"math"
	"reflect"
	"testing"

	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	models "github.com/thewizardplusplus/go-chess-models"
)

func TestSearcherAdapterSearchMove(
	test *testing.T,
) {
	type fields struct {
		moveSearcher MoveSearcher
	}
	type args struct {
		storage models.PieceStorage
		color   models.Color
	}
	type data struct {
		fields   fields
		args     args
		wantMove models.Move
		wantErr  bool
	}

	for _, data := range []data{
		data{
			fields: fields{
				moveSearcher: MockMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds moves.Bounds,
					) (moves.ScoredMove, error) {
						_, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}
						if deep != 0 {
							test.Fail()
						}
						if !reflect.DeepEqual(
							bounds,
							moves.Bounds{
								Alpha: math.Inf(-1),
								Beta:  math.Inf(+1),
							},
						) {
							test.Fail()
						}

						move := moves.ScoredMove{
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
							Score: 2.3,
						}
						return move, nil
					},
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
			},
			wantMove: models.Move{
				Start: models.Position{
					File: 1,
					Rank: 2,
				},
				Finish: models.Position{
					File: 3,
					Rank: 4,
				},
			},
			wantErr: false,
		},
		data{
			fields: fields{
				moveSearcher: MockMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds moves.Bounds,
					) (moves.ScoredMove, error) {
						_, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}
						if deep != 0 {
							test.Fail()
						}
						if !reflect.DeepEqual(
							bounds,
							moves.Bounds{
								Alpha: math.Inf(-1),
								Beta:  math.Inf(+1),
							},
						) {
							test.Fail()
						}

						move := moves.ScoredMove{
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
							Score: 2.3,
						}
						return move, errors.New("dummy")
					},
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
			},
			wantMove: models.Move{
				Start: models.Position{
					File: 1,
					Rank: 2,
				},
				Finish: models.Position{
					File: 3,
					Rank: 4,
				},
			},
			wantErr: true,
		},
		data{
			fields: fields{
				moveSearcher: MockMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds moves.Bounds,
					) (moves.ScoredMove, error) {
						_, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}
						if deep != 0 {
							test.Fail()
						}
						if !reflect.DeepEqual(
							bounds,
							moves.Bounds{
								Alpha: math.Inf(-1),
								Beta:  math.Inf(+1),
							},
						) {
							test.Fail()
						}

						var move moves.ScoredMove
						return move, errors.New("dummy")
					},
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
			},
			wantMove: models.Move{},
			wantErr:  true,
		},
	} {
		adapter := SearcherAdapter{
			MoveSearcher: data.fields.
				moveSearcher,
		}
		gotMove, gotErr := adapter.SearchMove(
			data.args.storage,
			data.args.color,
		)

		if !reflect.DeepEqual(
			gotMove,
			data.wantMove,
		) {
			test.Fail()
		}

		hasErr := gotErr != nil
		if hasErr != data.wantErr {
			test.Fail()
		}
	}
}
