package repositories

import (
	errs "github.com/apm-dev/go-clean-architecture/core/errors"
	"github.com/apm-dev/go-clean-architecture/data/datasources/acl_service"
	"github.com/apm-dev/go-clean-architecture/domain/entities"
	"github.com/apm-dev/go-clean-architecture/domain/repositories"
)

func NewACLRepository(ds acl_service.ACLDataSource) repositories.ACLRepository {
	return &aclRepository{ds: ds}
}

type aclRepository struct {
	ds acl_service.ACLDataSource
}

func (a *aclRepository) Call(id entities.UserID, m entities.Method) (bool, *errs.Error) {
	const op errs.Op = "repository.acl.Call"
	ok, err := a.ds.HasAccess(int64(id), string(m))
	if err != nil {
		return false, errs.E(op, err)
	}
	return ok, nil
}
