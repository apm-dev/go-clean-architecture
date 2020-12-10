package errors

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
)

//	Error kinds based on grpc codes
const (
	KindNotFound        = codes.NotFound
	KindInvalidArgument = codes.InvalidArgument
	KindUnauthenticated = codes.Unauthenticated
	KindUnauthorized    = codes.PermissionDenied
	KindInternal        = codes.Internal
	KindUnexpected      = codes.Unknown

	LevelInfo  = logrus.InfoLevel
	LevelDebug = logrus.DebugLevel
	LevelWarn  = logrus.WarnLevel
	LevelError = logrus.ErrorLevel
	LevelPanic = logrus.PanicLevel
)

type Op string

type Error struct {
	Op       Op
	Kind     codes.Code
	Err      error
	Severity logrus.Level
}

func E(args ...interface{}) *Error {
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case error:
			e.Err = arg
		case codes.Code:
			e.Kind = arg
		case logrus.Level:
			e.Severity = arg
		default:
			panic(fmt.Sprintf("bad call to E: %v", arg))
		}
	}
	return e
}

func Ops(e *Error) []Op {
	res := []Op{e.Op}
	subErr, ok := e.Err.(*Error)
	if !ok {
		return res
	}
	res = append(res, Ops(subErr)...)
	return res
}

func Kind(err error) codes.Code {
	e, ok := err.(*Error)
	if !ok {
		return KindUnexpected
	}
	if e.Kind != 0 {
		return e.Kind
	}
	return Kind(e.Err)
}

func Level(err error) logrus.Level {
	e, ok := err.(*Error)
	if !ok {
		return logrus.ErrorLevel
	}
	if e.Severity != 0 {
		return e.Severity
	}
	return Level(e.Err)
}

func (e *Error) Error() string {
	if subErr, ok := e.Err.(*Error); ok {
		return subErr.Error()
	}
	return e.Err.Error()
}
