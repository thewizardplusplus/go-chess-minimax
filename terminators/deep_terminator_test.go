package terminators

import (
	"testing"
)

func TestNewDeepTerminator(test *testing.T) {
	terminator := NewDeepTerminator(5)

	if terminator.maximalDeep != 5 {
		test.Fail()
	}
}

func TestDeepTerminatorIsSearchTerminated(test *testing.T) {
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
		{
			fields: fields{5},
			args:   args{4},
			want:   false,
		},
		{
			fields: fields{5},
			args:   args{5},
			want:   true,
		},
		{
			fields: fields{5},
			args:   args{6},
			want:   true,
		},
	} {
		terminator := DeepTerminator{
			maximalDeep: data.fields.maximalDeep,
		}
		got := terminator.IsSearchTerminated(data.args.deep)

		if got != data.want {
			test.Fail()
		}
	}
}

func TestDeepTerminatorSearchProgress(test *testing.T) {
	type fields struct {
		maximalDeep int
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
			fields: fields{100},
			args:   args{0},
			want:   0,
		},
		{
			fields: fields{100},
			args:   args{75},
			want:   0.75,
		},
		{
			fields: fields{100},
			args:   args{100},
			want:   1,
		},
		{
			fields: fields{100},
			args:   args{110},
			want:   1,
		},
	} {
		terminator := DeepTerminator{
			maximalDeep: data.fields.maximalDeep,
		}
		got := terminator.SearchProgress(data.args.deep)

		if got != data.want {
			test.Fail()
		}
	}
}
