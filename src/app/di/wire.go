//+build wireinject

package di

import (
	"github.com/apm-dev/go-clean-architecture/data/datasources/acl_service"
	"github.com/apm-dev/go-clean-architecture/data/repositories"
	"github.com/apm-dev/go-clean-architecture/domain/usecases"
	"github.com/google/wire"
)

var (
	//	Third-Party

	//	DataSources
	aclDS = wire.NewSet(acl_service.NewACLDataSource)

	//	Repositories
	aclRepository = wire.NewSet(repositories.NewACLRepository, aclDS)

	//	UseCases
	checkAccessUC = wire.NewSet(usecases.NewCheckAccess, aclRepository)

	//	Presentations
)

func CheckAccess() usecases.CheckAccess {
	wire.Build(checkAccessUC)
	return nil
}
