package repositories

import (
	"github.com/apm-dev/go-clean-architecture/core/errors"
	"github.com/apm-dev/go-clean-architecture/domain/entities"
)

type ACLRepository interface {
	Call(id entities.UserID, m entities.Method) (bool, *errors.Error)
}
