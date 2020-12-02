package repositories

import (
	"github.com/apm-dev/go-clean-architecture/core/errors"
	"github.com/apm-dev/go-clean-architecture/domain/entities"
)

type BlogRepository interface {
	Create(blog entities.Blog) (*entities.Blog, *errors.Error)
}
