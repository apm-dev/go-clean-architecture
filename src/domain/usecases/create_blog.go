package usecases

import (
	"github.com/apm-dev/go-clean-architecture/core/errors"
	"github.com/apm-dev/go-clean-architecture/domain/entities"
	"github.com/apm-dev/go-clean-architecture/domain/repositories"
)

type CreateBlog interface {
	Call(blog entities.Blog) (*entities.Blog, *errors.Error)
}

func NewCreateBlog(r repositories.BlogRepository) CreateBlog {
	return &createBlog{
		repo: r,
	}
}

type createBlog struct {
	repo repositories.BlogRepository
}

func (s *createBlog) Call(b entities.Blog) (*entities.Blog, *errors.Error) {
	const op errors.Op = "usecase.createBlog"
	blog, err := s.repo.Create(b)
	if err != nil {
		err = errors.E(op, err)
	}
	return blog, err
}
