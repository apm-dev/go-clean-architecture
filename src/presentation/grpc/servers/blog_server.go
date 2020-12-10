package servers

import (
	"context"
	"github.com/apm-dev/go-clean-architecture/core/errors"
	"github.com/apm-dev/go-clean-architecture/core/logger"
	"github.com/apm-dev/go-clean-architecture/domain/entities"
	"github.com/apm-dev/go-clean-architecture/domain/usecases"
	"github.com/apm-dev/go-clean-architecture/presentation/grpc/pbs"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/status"
)

func NewBlogServer(cb usecases.CreateBlog) *BlogServer {
	return &BlogServer{
		CreateBlog: cb,
	}
}

type BlogServer struct {
	CreateBlog usecases.CreateBlog
}

func protoBlogFromBlog(b *entities.Blog) *pbs.Blog {
	return &pbs.Blog{
		Id:        b.ID,
		AuthorId:  b.AuthorID,
		Title:     b.Title,
		Content:   b.Content,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

func (b *BlogServer) Create(ctx context.Context, req *pbs.CreateBlogRequest) (*pbs.CreateBlogResponse, error) {
	blog, err := b.CreateBlog.Call(entities.Blog{
		AuthorID: metautils.ExtractIncoming(ctx).Get("user_id"),
		Title:    req.GetTitle(),
		Content:  req.GetContent(),
	})
	if err != nil {
		logger.SysError(err)
		return nil, status.Error(errors.Kind(err), "there was a problem")
	}
	return &pbs.CreateBlogResponse{Blog: protoBlogFromBlog(blog)}, nil
}
