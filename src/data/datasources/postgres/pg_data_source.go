package postgres

import (
	"context"
	"fmt"
	errs "github.com/apm-dev/go-clean-architecture/core/errors"
	"github.com/apm-dev/go-clean-architecture/data/models"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type PgDataSource interface {
	CreateTables() *errs.Error
	InsertBlog(*models.BlogModel) (int64, *errs.Error)
}

func NewPgDataSource(db *pg.DB) PgDataSource {
	if db == nil {
		panic("postgres db instance should not be nil")
	}
	if err := db.Ping(context.TODO()); err != nil {
		panic(fmt.Sprintf("can not connect to postgres: %s", err))
	}
	return &pgDataSource{db: db}
}

type pgDataSource struct {
	db *pg.DB
}

func (p *pgDataSource) InsertBlog(b *models.BlogModel) (int64, *errs.Error) {
	const op errs.Op = "data_sources.postgres.InsertBlog"
	//	Skip blog.ID if assigned
	b.ID = ""
	//	Cast general blog model to postgres specific blog model
	pb := new(PgBlogModel)
	err := pb.NewFromBlogModel(*b)
	if err != nil {
		return 0, errs.E(op, err)
	}
	//	Insert pg blog model to database
	_, dbErr := p.db.Model(pb).Insert()
	if dbErr != nil {
		return 0, errs.E(op, dbErr, errs.LevelError, errs.KindInternal)
	}
	return pb.ID, nil
}

func (p *pgDataSource) CreateTables() *errs.Error {
	const op errs.Op = "data_sources.postgres.createTables"
	for _, model := range PgModels() {
		err := p.db.Model(model).CreateTable(&orm.CreateTableOptions{})
		if err != nil {
			return errs.E(op, errs.KindInternal, err)
		}
	}
	return nil
}
