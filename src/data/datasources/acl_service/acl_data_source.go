package acl_service

import "github.com/apm-dev/go-clean-architecture/core/errors"

type ACLDataSource interface {
	HasAccess(id int64, method string) (bool, *errors.Error)
}

func NewACLDataSource() ACLDataSource {
	return &aclDataSource{}
}

type aclDataSource struct {
	//	dependencies
}

func (a *aclDataSource) HasAccess(id int64, method string) (bool, *errors.Error) {
	//	send grpc request to acl service and check access
	if id%2 == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
