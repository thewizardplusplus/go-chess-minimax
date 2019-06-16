package chessminimax

import (
	"reflect"
	"testing"
)

func TestDeepTerminatorInterface(
	test *testing.T,
) {
	gotType := reflect.TypeOf(DeepTerminator{})
	wantType := reflect.
		TypeOf((*SearchTerminator)(nil)).
		Elem()
	if !gotType.Implements(wantType) {
		test.Fail()
	}
}

func TestNewDeepTerminator(test *testing.T) {
	terminator := NewDeepTerminator(5)

	if terminator.maximalDeep != 5 {
		test.Fail()
	}
}

func TestDeepTerminatorIsSearchTerminate(
	test *testing.T,
) {
	type fields struct {
		maximalDeep int
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
			fields: fields{5},
			args:   args{4},
			want:   false,
		},
		data{
			fields: fields{5},
			args:   args{5},
			want:   true,
		},
		data{
			fields: fields{5},
			args:   args{6},
			want:   true,
		},
	} {
		terminator := DeepTerminator{
			maximalDeep: data.fields.maximalDeep,
		}
		got := terminator.IsSearchTerminate(
			data.args.deep,
		)

		if got != data.want {
			test.Fail()
		}
	}
}
