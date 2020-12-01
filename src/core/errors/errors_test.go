package errors

import (
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"reflect"
	"testing"
)

func TestE(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			"Simple",
			args{args: []interface{}{Op("blog.findByID"), KindNotFound, errors.New("not found error")}},
			&Error{
				Op:   "blog.findByID",
				Kind: KindNotFound,
				Err:  errors.New("not found error"),
			},
		},
		{
			"Nested",
			args{args: []interface{}{
				Op("blog.create"),
				KindUnauthorized,
				&Error{
					Op:   "account.getUser",
					Kind: KindNotFound,
					Err:  errors.New("user not found error"),
				},
			}},
			&Error{
				Op:   "blog.create",
				Kind: KindUnauthorized,
				Err: &Error{
					Op:   "account.getUser",
					Kind: KindNotFound,
					Err:  errors.New("user not found error"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := E(tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("E() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	type fields struct {
		Op   Op
		Kind codes.Code
		Err  error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"Error with nested error",
			fields{
				Op:   "blog.findByID",
				Kind: KindNotFound,
				Err:  errors.New("blog not found"),
			},
			fmt.Sprintf("K:5  Op:blog.findByID  Err:blog not found"),
		},
		{
			"Error with nested Error",
			fields{
				Op:   "blog.findByID",
				Kind: KindNotFound,
				Err:  E(Op("account.getUser"), errors.New("unexpected error")),
			},
			"K:5  Op:blog.findByID  Err:\n\tK:0  Op:account.getUser  Err:unexpected error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Op:   tt.fields.Op,
				Kind: tt.fields.Kind,
				Err:  tt.fields.Err,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOps(t *testing.T) {
	type args struct {
		e *Error
	}
	tests := []struct {
		name string
		args args
		want []Op
	}{
		{
			"Nested Errors",
			args{e: E(Op("blog.findByID"), E(Op("account.getUser")))},
			[]Op{"blog.findByID", "account.getUser"},
		},
		{
			"Error with nested error",
			args{e: E(Op("blog.findByID"), errors.New("unexpected error"))},
			[]Op{"blog.findByID"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ops(tt.args.e); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ops() = %v, want %v", got, tt.want)
			}
		})
	}
}
