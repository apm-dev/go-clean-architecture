package postgres

import (
	errs "github.com/apm-dev/go-clean-architecture/core/errors"
	"github.com/apm-dev/go-clean-architecture/core/util"
	"github.com/apm-dev/go-clean-architecture/data/models"
	"time"
)

func PgModels() []interface{} {
	return []interface{}{
		(*PgBlogModel)(nil),
	}
}

type PgBlogModel struct {
	tableName struct{} `pg:"blogs"`
	ID        int64
	Title     string `pg:"type:varchar(128)"`
	Content   string
	AuthorID  int64
	CreatedAt time.Time `pg:"default:now()"`
	UpdatedAt time.Time `pg:"default:now()"`
}

func (pm *PgBlogModel) ToBlogModel() *models.BlogModel {
	//	TODO: implement ToBlogModel
	return nil
}

func (pm *PgBlogModel) NewFromBlogModel(b models.BlogModel) *errs.Error {
	const op errs.Op = "pg_models.newFromBlogModel"
	id, err := util.InputConverter.StringToUnsignedInt(b.ID)
	if err != nil {
		return errs.E(op, err)
	}
	authorId, err := util.InputConverter.StringToUnsignedInt(b.AuthorID)
	if err != nil {
		return errs.E(op, err)
	}
	pm.ID = int64(id)
	pm.Title = b.Title
	pm.Content = b.Content
	pm.AuthorID = int64(authorId)
	pm.CreatedAt = util.Time.FromUnix(b.CreatedAt)
	pm.UpdatedAt = util.Time.FromUnix(b.UpdatedAt)
	return nil
}
