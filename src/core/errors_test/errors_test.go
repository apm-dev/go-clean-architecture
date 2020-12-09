package errors_test

import (
	"errors"
	err "github.com/apm-dev/go-clean-architecture/core/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"testing"
)

func TestE(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		args args
		want *err.Error
	}{
		{
			"Simple",
			args{args: []interface{}{err.Op("blog.findByID"), err.KindNotFound, errors.New("not found error")}},
			&err.Error{
				Op:   "blog.findByID",
				Kind: err.KindNotFound,
				Err:  errors.New("not found error"),
			},
		},
		{
			"Nested",
			args{args: []interface{}{
				err.Op("blog.create"),
				err.KindUnauthorized,
				&err.Error{
					Op:   "account.getUser",
					Kind: err.KindNotFound,
					Err:  errors.New("user not found error"),
				},
			}},
			&err.Error{
				Op:   "blog.create",
				Kind: err.KindUnauthorized,
				Err: &err.Error{
					Op:   "account.getUser",
					Kind: err.KindNotFound,
					Err:  errors.New("user not found error"),
				},
			},
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, err.E(tt.args.args...), tt.name)
	}
}

func TestError_Error(t *testing.T) {
	type fields struct {
		Op   err.Op
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
				Kind: err.KindNotFound,
				Err:  errors.New("blog not found"),
			},
			"blog not found",
		},
		{
			"Error with nested Error",
			fields{
				Op:   "blog.findByID",
				Kind: err.KindNotFound,
				Err:  err.E(err.Op("account.getUser"), errors.New("unexpected error")),
			},
			"unexpected error",
		},
	}
	for _, tt := range tests {
		e := &err.Error{
			Op:   tt.fields.Op,
			Kind: tt.fields.Kind,
			Err:  tt.fields.Err,
		}
		assert.Equal(t, tt.want, e.Error(), tt.name)
	}
}

func TestOps(t *testing.T) {
	type args struct {
		e *err.Error
	}
	tests := []struct {
		name string
		args args
		want []err.Op
	}{
		{
			"Nested Errors",
			args{e: err.E(err.Op("blog.findByID"), err.E(err.Op("account.getUser")))},
			[]err.Op{"blog.findByID", "account.getUser"},
		},
		{
			"Error with nested error",
			args{e: err.E(err.Op("blog.findByID"), errors.New("unexpected error"))},
			[]err.Op{"blog.findByID"},
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, err.Ops(tt.args.e), tt.name)
	}
}
