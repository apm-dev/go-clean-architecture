package usecases

import (
	"errors"
	errs "github.com/apm-dev/go-clean-architecture/core/errors"
	"github.com/apm-dev/go-clean-architecture/domain/entities"
	"github.com/apm-dev/go-clean-architecture/domain/repositories"
)

type CheckAccess interface {
	Call(id entities.UserID, m entities.Method) (bool, *errs.Error)
}

func NewCheckAccess(r repositories.ACLRepository) CheckAccess {
	return &checkAccess{repo: r}
}

type checkAccess struct {
	repo repositories.ACLRepository
}

func (a *checkAccess) Call(id entities.UserID, m entities.Method) (bool, *errs.Error) {
	const op errs.Op = "usecase.CheckAccess"
	if _, ok := entities.AllMethods()[m]; !ok {
		return false, errs.E(op, errs.KindInternal, errors.New("method not exists"))
	}
	ok, err := a.repo.Call(id, m)
	if err != nil {
		return false, errs.E(op, err)
	}
	return ok, nil
}
