package chessminimax

import (
	"reflect"
	"testing"
)

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
