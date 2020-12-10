package interceptors

import (
	"context"
	"fmt"
	"github.com/apm-dev/go-clean-architecture/app/di"
	errs "github.com/apm-dev/go-clean-architecture/core/errors"
	"github.com/apm-dev/go-clean-architecture/core/logger"
	"github.com/apm-dev/go-clean-architecture/core/util"
	"github.com/apm-dev/go-clean-architecture/domain/entities"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Method string

func toMethod(method string) (entities.Method, error) {
	const blogServicePath = "/blog.BlogService/"

	m, ok := map[string]entities.Method{
		blogServicePath + "Create": entities.CreateBlog,
	}[method]
	if !ok {
		return "", fmt.Errorf("method %s not implemented", method)
	}
	return m, nil
}

func AuthFunc(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	const op errs.Op = "grpc.interceptor.acl"
	//	Get user_id from context(request header)
	strID := metautils.ExtractIncoming(ctx).Get("user_id")
	id, icErr := util.InputConverter.StringToUnsignedInt(strID)
	if icErr != nil {
		logger.SysError(errs.E(op, icErr, errs.LevelInfo))
		return nil, status.Error(codes.InvalidArgument, "user_id should be unsigned integer")
	}

	//	Convert grpc method name to our Method type
	m, mErr := toMethod(info.FullMethod)
	if mErr != nil {
		logger.SysError(errs.E(op, mErr, errs.LevelError))
		return nil, status.Error(codes.Unimplemented, "method not implemented")
	}

	//	Check access
	ok, err := di.CheckAccess().Call(entities.UserID(id), m)
	if err != nil {
		logger.SysError(errs.E(op, err))
		return nil, status.Error(errs.Kind(err), "there was a problem")
	}

	//	return PermissionDenied error if user has no access to method
	if !ok {
		logger.SysError(errs.E(
			op,
			errs.KindUnauthorized,
			fmt.Errorf("unauthorized access to '%s':'%s' with user_id: %d", info.FullMethod, m, id),
			errs.LevelWarn,
		))
		return nil, status.Errorf(codes.PermissionDenied, "user: %d don't have access to method: %s", id,
			info.FullMethod)
	}
	//	dispatch request if access was ok
	return handler(ctx, req)
}
