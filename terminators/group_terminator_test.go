package terminators

import (
	"reflect"
	"testing"
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

func TestGroupTerminatorInterface(
	test *testing.T,
) {
	gotType := reflect.TypeOf(
		GroupTerminator{},
	)
	wantType := reflect.
		TypeOf((*SearchTerminator)(nil)).
		Elem()
	if !gotType.Implements(wantType) {
		test.Fail()
	}
}

func TestNewGroupTerminator(
	test *testing.T,
) {
	type args struct {
		terminators []SearchTerminator
	}
	type data struct {
		args args
	}

	for _, data := range []data{
		data{
			args: args{nil},
		},
		data{
			args: args{
				terminators: []SearchTerminator{
					MockSearchTerminator{},
					MockSearchTerminator{},
				},
			},
		},
	} {
		group := NewGroupTerminator(
			data.args.terminators...,
		)

		if !reflect.DeepEqual(
			group.terminators,
			data.args.terminators,
		) {
			test.Fail()
		}
	}
}

func TestGroupTerminatorIsSearchTerminate(
	test *testing.T,
) {
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
		data{
			fields: fields{nil},
			args:   args{5},
			want:   false,
		},
		data{
			fields: fields{
				terminators: []SearchTerminator{
					MockSearchTerminator{
						isSearchTerminate: func(
							deep int,
						) bool {
							if deep != 5 {
								test.Fail()
							}

							return false
						},
					},
					MockSearchTerminator{
						isSearchTerminate: func(
							deep int,
						) bool {
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
		data{
			fields: fields{
				terminators: []SearchTerminator{
					MockSearchTerminator{
						isSearchTerminate: func(
							deep int,
						) bool {
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
		got := group.IsSearchTerminate(
			data.args.deep,
		)

		if got != data.want {
			test.Fail()
		}
	}
}
