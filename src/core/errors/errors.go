package errors

import (
	"fmt"
	"google.golang.org/grpc/codes"
)

//	Error kinds based on grpc codes
const (
	KindNotFound        = codes.NotFound
	KindInvalidArgument = codes.InvalidArgument
	KindUnauthorized    = codes.Unauthenticated
	KindInternal        = codes.Internal
	KindUnexpected      = codes.Unknown
)

type Op string

type Error struct {
	Op   Op
	Kind codes.Code
	Err  error
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

func (e *Error) Error() string {
	if subErr, ok := e.Err.(*Error); ok {
		return fmt.Sprintf("K:%d  Op:%s  Err:\n\t%s", e.Kind, e.Op, subErr.Error())
	}
	return fmt.Sprintf("K:%d  Op:%s  Err:%s", e.Kind, e.Op, e.Err)
}
