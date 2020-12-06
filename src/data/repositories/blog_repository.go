package repositories

import (
	"github.com/apm-dev/go-clean-architecture/core/errors"
	"github.com/apm-dev/go-clean-architecture/data/datasources/postgres"
	"github.com/apm-dev/go-clean-architecture/data/models"
	"github.com/apm-dev/go-clean-architecture/domain/entities"
	"github.com/apm-dev/go-clean-architecture/domain/repositories"
	"strconv"
)

func NewBlogRepository(pgds postgres.PgDataSource) repositories.BlogRepository {
	return &blogRepository{pgds: pgds}
}

type blogRepository struct {
	pgds postgres.PgDataSource
}

func (b *blogRepository) Create(blog entities.Blog) (*entities.Blog, *errors.Error) {
	const op errors.Op = "blog_repository.Create"
	id, err := b.pgds.InsertBlog(&models.BlogModel{
		Title:     blog.Title,
		Content:   blog.Content,
		AuthorID:  blog.AuthorID,
		CreatedAt: blog.CreatedAt,
		UpdateAt:  blog.UpdatedAt,
	})
	if err != nil {
		return nil, errors.E(op, err)
	}
	blog.ID = strconv.Itoa(int(id))
	return &blog, nil
}
