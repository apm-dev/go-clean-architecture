package postgres

import (
	"github.com/apm-dev/go-clean-architecture/core/errors"
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
	return nil
}

func (pm *PgBlogModel) NewFromBlogModel(b models.BlogModel) *errors.Error {
	const op errors.Op = "pg_models.newFromBlogModel"
	id, err := util.InputConverter.StringToUnsignedInt(b.ID)
	if err != nil {
		return errors.E(op, err)
	}
	authorId, err := util.InputConverter.StringToUnsignedInt(b.AuthorID)
	if err != nil {
		return errors.E(op, err)
	}
	pm.ID = int64(id)
	pm.Title = b.Title
	pm.Content = b.Content
	pm.AuthorID = int64(authorId)
	pm.CreatedAt = time.Unix(b.CreatedAt, 0)
	pm.UpdatedAt = time.Unix(b.UpdateAt, 0)
	return nil
}
