package postgres_test

import (
	"context"
	"github.com/apm-dev/go-clean-architecture/data/datasources/postgres"
	"github.com/apm-dev/go-clean-architecture/data/models"
	"github.com/go-pg/pg/v10"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/suite"
	"log"
	"strconv"
	"testing"
	"time"
)

type PgTestSuite struct {
	suite.Suite
	db *pg.DB
}

var pool *dockertest.Pool
var resource *dockertest.Resource

func (s *PgTestSuite) SetupTest() {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	// pulls an image, creates a container based on it and runs it
	err = nil
	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.3",
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_DB=pg_test",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.Cgroup = "fsync=off synchronous_commit=off archive_mode=off wal_level=minimal shared_buffers=512MB"
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		s.db = pg.Connect(&pg.Options{
			Addr:     ":" + resource.GetPort("5432/tcp"),
			User:     "postgres",
			Password: "secret",
			Database: "pg_test",
		})
		return s.db.Ping(context.Background())
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
}

func (s *PgTestSuite) TearDownTest() {
	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	//os.Exit(0)
}

func (s *PgTestSuite) TestNewPgDataSource() {
	s.NotNil(postgres.NewPgDataSource(s.db), "new pg data source should not be nil")
}

func (s *PgTestSuite) TestCreateTables() {
	var err error
	err = postgres.NewPgDataSource(s.db).CreateTables()
	s.Nilf(err, "create table should not return error: %s", err)

	var tables []string
	_, err = s.db.Query(&tables, `
		SELECT tablename
		FROM pg_catalog.pg_tables
		WHERE schemaname = 'public'
	`)
	panicIf(err)
	mLen := len(postgres.PgModels())
	tLen := len(tables)
	s.EqualValuesf(mLen, tLen, "have %d pg_models but only %d tables", mLen, tLen)
}

func (s *PgTestSuite) TestCreateBlog() {
	//	Creating tables
	pgds := postgres.NewPgDataSource(s.db)
	err := pgds.CreateTables()
	if err != nil {
		log.Fatalf("failed when creating tables: %v", err)
	}

	//	Create a blog
	blog := models.BlogModel{
		Title:     "my title",
		Content:   "<html>my html content</html>",
		AuthorID:  "12",
		CreatedAt: time.Now().UTC().Unix(),
		UpdateAt:  time.Now().UTC().Unix(),
	}
	id, err := pgds.InsertBlog(&blog)
	if err != nil {
		log.Fatalf("insert blog error: %s", err)
	}
	s.NotEqualValues(0, id, "created blog id should not be 0")

	//	Find created blog in DB
	pgBlog := new(postgres.PgBlogModel)
	dbErr := s.db.Model(pgBlog).Where("id = ?", id).Select()
	if dbErr != nil {
		log.Fatalf("find blog: %d error: %s", id, dbErr)
	}
	s.Equal(blog.Title, pgBlog.Title, "blog title is not equal to db row title")
	s.EqualValues(blog.CreatedAt, pgBlog.CreatedAt.Unix(), "blog createdAt is not equal to db row createdAt")

	//	Create blog with id, it should skip it's id and generate new one
	fakeId := 14323
	blog.ID = strconv.Itoa(fakeId)
	id2, err := pgds.InsertBlog(&blog)
	if err != nil {
		log.Fatalf("insert blog error: %v", err)
	}
	s.NotEqualValues(id2, fakeId, "id of blog that passed to insertBlog should skip and don't persist in DB")
}

func TestPgTestSuite(t *testing.T) {
	suite.Run(t, new(PgTestSuite))
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
