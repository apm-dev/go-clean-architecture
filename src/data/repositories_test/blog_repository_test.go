package repositories_test

import (
	"errors"
	errs "github.com/apm-dev/go-clean-architecture/core/errors"
	"github.com/apm-dev/go-clean-architecture/core/util"
	"github.com/apm-dev/go-clean-architecture/data/datasources/postgres"
	"github.com/apm-dev/go-clean-architecture/data/mocks"
	"github.com/apm-dev/go-clean-architecture/data/models"
	"github.com/apm-dev/go-clean-architecture/data/repositories"
	"github.com/apm-dev/go-clean-architecture/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestNewBlogRepository(t *testing.T) {
	pgdsMock := new(mocks.PgDataSource)
	assert.NotNil(t, repositories.NewBlogRepository(pgdsMock), "NewBlogRepository() should not return nil")
}

func TestCreate(t *testing.T) {
	pgdsMock := new(mocks.PgDataSource)
	pgdsBadMock := new(mocks.PgDataSource)
	blogModelType := reflect.TypeOf(&models.BlogModel{}).Name()

	pgdsMock.On("InsertBlog", mock.AnythingOfType(blogModelType)).Return(int64(1), nil)

	pgdsErr := errs.E(
		errs.Op("pg_data_source.InsertBlog"),
		errs.KindUnexpected,
		errors.New("there is no blogs table in database"),
	)
	pgdsBadMock.On("InsertBlog", mock.AnythingOfType(blogModelType)).Return(int64(0), pgdsErr)

	b1in := entities.Blog{
		AuthorID:  "15",
		Title:     "new go version released",
		Content:   "<html><h1>This is awesome</h1></html>",
		CreatedAt: util.Time.NowUnix(),
		UpdatedAt: util.Time.NowUnix(),
	}

	b1want := &b1in
	b1want.ID = "1"

	type fields struct {
		pgds postgres.PgDataSource
	}

	tests := []struct {
		msg    string
		fields fields
		in     entities.Blog
		want   *entities.Blog
		err    *errs.Error
	}{
		{
			msg:    "should return blog pointer with filled id when everything is correct",
			fields: fields{pgds: pgdsMock},
			in:     b1in,
			want:   b1want,
			err:    nil,
		},
		{
			msg:    "should return wrapped error when postgres data source returns an error",
			fields: fields{pgds: pgdsBadMock},
			in:     b1in,
			want:   nil,
			err:    errs.E(errs.Op("blog_repository.Create"), pgdsErr),
		},
	}

	for _, tt := range tests {
		blog, err := repositories.NewBlogRepository(tt.fields.pgds).Create(tt.in)
		assert.Equal(t, tt.want, blog, tt.msg)
		assert.Equal(t, tt.err, err, tt.msg)
	}
	pgdsMock.AssertExpectations(t)
	pgdsBadMock.AssertExpectations(t)
}
