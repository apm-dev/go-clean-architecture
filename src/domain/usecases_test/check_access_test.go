package usecases

import (
	"errors"
	errs "github.com/apm-dev/go-clean-architecture/core/errors"
	"github.com/apm-dev/go-clean-architecture/domain/entities"
	"github.com/apm-dev/go-clean-architecture/domain/mocks"
	"github.com/apm-dev/go-clean-architecture/domain/usecases"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewACL(t *testing.T) {
	repo := new(mocks.ACLRepository)
	assert.NotNil(t, usecases.NewCheckAccess(repo), "acl instance should not be nil")
}

func TestHasAccess(t *testing.T) {

	type args struct {
		id     entities.UserID
		method entities.Method
	}

	tests := []struct {
		msg  string
		args args
		want bool
		err  *errs.Error
	}{
		{
			msg: "should return true when repository returns true",
			args: args{
				id:     entities.UserID(1),
				method: entities.CreateBlog,
			},
			want: true,
			err:  nil,
		},
		{
			msg: "should return false when repository returns false",
			args: args{
				id:     entities.UserID(2),
				method: entities.CreateBlog,
			},
			want: false,
			err:  nil,
		},
		{
			msg: "should return error when method not exists",
			args: args{
				id:     entities.UserID(1),
				method: entities.Method("wrong.method"),
			},
			want: false,
			err: errs.E(
				errs.Op("usecase.CheckAccess"),
				errs.KindInternal,
				errors.New("method not exists"),
			),
		},
	}

	repo := new(mocks.ACLRepository)
	repo.On("Call", entities.UserID(1), entities.CreateBlog).Return(true, nil)
	repo.On("Call", entities.UserID(2), entities.CreateBlog).Return(false, nil)

	acl := usecases.NewCheckAccess(repo)

	for _, tt := range tests {
		ok, err := acl.Call(tt.args.id, tt.args.method)
		assert.EqualValues(t, tt.want, ok, tt.msg)
		assert.EqualValues(t, tt.err, err, tt.msg)
	}
	repo.AssertExpectations(t)
}
