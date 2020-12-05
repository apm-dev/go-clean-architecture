package usecases_test

import (
	"errors"
	err "github.com/apm-dev/go-clean-architecture/core/errors"
	"github.com/apm-dev/go-clean-architecture/domain/entities"
	"github.com/apm-dev/go-clean-architecture/domain/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockBlogRepository struct {
	mock.Mock
}

func (m *MockBlogRepository) Create(blog entities.Blog) (*entities.Blog, *err.Error) {
	args := m.Called(blog)
	r0, ok := args.Get(0).(*entities.Blog)
	if !ok {
		r0 = nil
	}
	r1, ok := args.Get(1).(*err.Error)
	if !ok {
		r1 = nil
	}
	return r0, r1
}

func TestNewCreateBlog(t *testing.T) {
	mockRepo := new(MockBlogRepository)
	cb, ok := usecases.NewCreateBlog(mockRepo).(usecases.CreateBlog)
	assert.True(t, ok, "instance should be of type CreateBlog")
	assert.NotNil(t, cb, "instance should not be nil")
}

func Test_CreateBlog_Call(t *testing.T) {
	rp0in := entities.Blog{
		AuthorID:  "1",
		Title:     "New Blog",
		Content:   "<html>some content</html>",
		CreatedAt: time.Now().UTC().Unix(),
		UpdatedAt: time.Now().UTC().Unix(),
	}
	rp0out := rp0in
	rp0out.ID = "1"
	var rp0err *err.Error = nil

	rp1in := entities.Blog{
		AuthorID:  "2",
		Title:     "Another Blog",
		Content:   "<html>some other content</html>",
		CreatedAt: time.Now().UTC().Unix(),
		UpdatedAt: time.Now().UTC().Unix(),
	}
	var rp1out *entities.Blog = nil
	rp1err := err.E(err.Op("repository.createBlog"), errors.New("author not found"))

	mockRepo := new(MockBlogRepository)
	mockRepo.On("Create", rp0in).Return(&rp0out, rp0err)
	mockRepo.On("Create", rp1in).Return(rp1out, rp1err)

	tests := []struct {
		name    string
		input   entities.Blog
		want    *entities.Blog
		wantErr *err.Error
	}{
		{
			name:    "successful blog creation",
			input:   rp0in,
			want:    &rp0out,
			wantErr: nil,
		},
		{
			name:    "failed blog creation",
			input:   rp1in,
			want:    nil,
			wantErr: err.E(err.Op("usecase.createBlog"), rp1err),
		},
	}

	cb := usecases.NewCreateBlog(mockRepo)
	for _, tt := range tests {
		got, gotErr := cb.Call(tt.input)
		assert.Equal(t, tt.want, got, tt.name)
		assert.Equal(t, tt.wantErr, gotErr, tt.name)
	}
	mockRepo.AssertExpectations(t)
}
