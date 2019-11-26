package terminators

import (
	"reflect"
	"testing"
)

type MockSearchTerminator struct {
	isSearchTerminated func(deep int) bool
	searchProgress     func(deep int) float64
}

func (terminator MockSearchTerminator) IsSearchTerminated(deep int) bool {
	if terminator.isSearchTerminated == nil {
		panic("not implemented")
	}

	return terminator.isSearchTerminated(deep)
}

func (terminator MockSearchTerminator) SearchProgress(deep int) float64 {
	if terminator.searchProgress == nil {
		panic("not implemented")
	}

	return terminator.searchProgress(deep)
}

func TestNewGroupTerminator(test *testing.T) {
	type args struct {
		terminators []SearchTerminator
	}
	type data struct {
		args args
	}

	for _, data := range []data{
		{
			args: args{nil},
		},
		{
			args: args{
				terminators: []SearchTerminator{
					MockSearchTerminator{},
					MockSearchTerminator{},
				},
			},
		},
	} {
		group := NewGroupTerminator(data.args.terminators...)

		if !reflect.DeepEqual(group.terminators, data.args.terminators) {
			test.Fail()
		}
	}
}

func TestGroupTerminatorIsSearchTerminated(test *testing.T) {
	type fields struct {
		terminators []SearchTerminator
	}
	type args struct {
		deep int
	}
	type data struct {
		fields fields
		args   args
		want   bool
	}

	for _, data := range []data{
		{
			fields: fields{nil},
			args:   args{5},
			want:   false,
		},
		{
			fields: fields{
				terminators: []SearchTerminator{
					MockSearchTerminator{
						isSearchTerminated: func(deep int) bool {
							if deep != 5 {
								test.Fail()
							}

							return false
						},
					},
					MockSearchTerminator{
						isSearchTerminated: func(deep int) bool {
							if deep != 5 {
								test.Fail()
							}

							return false
						},
					},
				},
			},
			args: args{5},
			want: false,
		},
		{
			fields: fields{
				terminators: []SearchTerminator{
					MockSearchTerminator{
						isSearchTerminated: func(deep int) bool {
							if deep != 5 {
								test.Fail()
							}

							return true
						},
					},
					MockSearchTerminator{},
				},
			},
			args: args{5},
			want: true,
		},
	} {
		group := GroupTerminator{
			terminators: data.fields.terminators,
		}
		got := group.IsSearchTerminated(data.args.deep)

		if got != data.want {
			test.Fail()
		}
	}
}

func TestGroupTerminatorSearchProgress(test *testing.T) {
	type fields struct {
		terminators []SearchTerminator
	}
	type args struct {
		deep int
	}
	type data struct {
		fields fields
		args   args
		want   float64
	}

	for _, data := range []data{
		{
			fields: fields{nil},
			args:   args{5},
			want:   0,
		},
		{
			fields: fields{
				terminators: []SearchTerminator{
					MockSearchTerminator{
						isSearchTerminated: func(deep int) bool {
							if deep != 5 {
								test.Fail()
							}

							return false
						},
						searchProgress: func(deep int) float64 {
							if deep != 5 {
								test.Fail()
							}

							return 0.25
						},
					},
					MockSearchTerminator{
						isSearchTerminated: func(deep int) bool {
							if deep != 5 {
								test.Fail()
							}

							return false
						},
						searchProgress: func(deep int) float64 {
							if deep != 5 {
								test.Fail()
							}

							return 0.75
						},
					},
					MockSearchTerminator{
						isSearchTerminated: func(deep int) bool {
							if deep != 5 {
								test.Fail()
							}

							return false
						},
						searchProgress: func(deep int) float64 {
							if deep != 5 {
								test.Fail()
							}

							return 0.5
						},
					},
				},
			},
			args: args{5},
			want: 0.75,
		},
		{
			fields: fields{
				terminators: []SearchTerminator{
					MockSearchTerminator{
						isSearchTerminated: func(deep int) bool {
							if deep != 5 {
								test.Fail()
							}

							return true
						},
						searchProgress: func(deep int) float64 {
							if deep != 5 {
								test.Fail()
							}

							return 0.25
						},
					},
					MockSearchTerminator{
						searchProgress: func(deep int) float64 {
							if deep != 5 {
								test.Fail()
							}

							return 0.75
						},
					},
					MockSearchTerminator{
						searchProgress: func(deep int) float64 {
							if deep != 5 {
								test.Fail()
							}

							return 0.5
						},
					},
				},
			},
			args: args{5},
			want: 1,
		},
	} {
		group := GroupTerminator{
			terminators: data.fields.terminators,
		}
		got := group.SearchProgress(data.args.deep)

		if got != data.want {
			test.Fail()
		}
	}
}
