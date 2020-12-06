package postgres

import (
	"github.com/apm-dev/go-clean-architecture/core/errors"
	"github.com/apm-dev/go-clean-architecture/data/models"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type PgDataSource interface {
	CreateTables() *errors.Error
	InsertBlog(*models.BlogModel) (int64, *errors.Error)
}

func NewPgDataSource(db *pg.DB) PgDataSource {
	if db == nil {
		panic("postgres db instance should not be nil")
	}
	return &pgDataSource{db: db}
}

type pgDataSource struct {
	db *pg.DB
}

func (p *pgDataSource) InsertBlog(b *models.BlogModel) (int64, *errors.Error) {
	const op errors.Op = "pg_data_source.InsertBlog"
	b.ID = ""
	pb := new(PgBlogModel)
	err := pb.NewFromBlogModel(*b)
	if err != nil {
		return 0, errors.E(op, err)
	}
	_, dbErr := p.db.Model(pb).Insert()
	if dbErr != nil {
		return 0, errors.E(op, dbErr)
	}
	return pb.ID, nil
}

func (p *pgDataSource) CreateTables() *errors.Error {
	for _, model := range PgModels() {
		err := p.db.Model(model).CreateTable(&orm.CreateTableOptions{})
		if err != nil {
			return errors.E(errors.Op("pg_data_source.createTables"), errors.KindUnexpected, err)
		}
	}
	return nil
}
