//+build wireinject

package di

import (
	"github.com/apm-dev/go-clean-architecture/data/datasources/acl_service"
	"github.com/apm-dev/go-clean-architecture/data/datasources/postgres"
	"github.com/apm-dev/go-clean-architecture/data/repositories"
	"github.com/apm-dev/go-clean-architecture/domain/usecases"
	"github.com/apm-dev/go-clean-architecture/presentation/grpc/servers"
	"github.com/go-pg/pg/v10"
	"github.com/google/wire"
)

var (
	//	Third-Party
	postgresDB = wire.NewSet(providePostgresDB)

	//	DataSources
	postgresDS = wire.NewSet(postgres.NewPgDataSource, postgresDB)
	aclDS      = wire.NewSet(acl_service.NewACLDataSource)

	//	Repositories
	blogRepository = wire.NewSet(repositories.NewBlogRepository, postgresDS)
	aclRepository  = wire.NewSet(repositories.NewACLRepository, aclDS)

	//	UseCases
	createBlogUC  = wire.NewSet(usecases.NewCreateBlog, blogRepository)
	checkAccessUC = wire.NewSet(usecases.NewCheckAccess, aclRepository)

	//	Presentations
	grpcBlogServer = wire.NewSet(servers.NewBlogServer, createBlogUC)
)

func PostgresDB() *pg.DB {
	wire.Build(postgresDB)
	return nil
}

func PostgresDS() postgres.PgDataSource {
	wire.Build(postgresDS)
	return nil
}

func GrpcBlogServer() *servers.BlogServer {
	wire.Build(grpcBlogServer)
	return nil
}

func CheckAccess() usecases.CheckAccess {
	wire.Build(checkAccessUC)
	return nil
}
