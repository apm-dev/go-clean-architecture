package util

import (
	"errors"
	"fmt"
	errs "github.com/apm-dev/go-clean-architecture/core/errors"
	"strconv"
	"strings"
)

var InputConverter InputConverterInterface = &inputConverter{}

type InputConverterInterface interface {
	StringToUnsignedInt(i string) (int, *errs.Error)
}

type inputConverter struct {
}

func (*inputConverter) StringToUnsignedInt(i string) (int, *errs.Error) {
	i = strings.Split(i, ".")[0]
	if i == "" {
		return 0, nil
	}
	const op errs.Op = "inputConverter.stringToUnsignedInt"
	number, err := strconv.Atoi(i)
	if err != nil {
		return 0, errs.E(op, errs.KindUnexpected, err)
	}
	if number < 0 {
		return 0, errs.E(
			op, errs.KindUnexpected,
			errors.New(fmt.Sprintf("negative number: %d", number)),
		)
	}
	return number, nil
}
