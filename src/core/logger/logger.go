package logger

import (
	"github.com/apm-dev/go-clean-architecture/core/configs"
	"github.com/apm-dev/go-clean-architecture/core/errors"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//	ToDo: implement dynamic formatter based on config and env
	log.SetFormatter(&log.JSONFormatter{PrettyPrint: true})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	//	ToDo: implement dynamic output based on config and env
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	if configs.IsDebugMode() {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
}

func SysError(err error) {
	sysError, ok := err.(*errors.Error)
	if !ok {
		log.Error(err)
		return
	}

	entry := log.WithFields(log.Fields{
		"operations": errors.Ops(sysError),
		"kind":       errors.Kind(sysError),
	})

	switch errors.Level(err) {
	case log.WarnLevel:
		entry.Warnf("%s: %v\n", sysError.Op, err)
	case log.InfoLevel:
		entry.Infof("%s: %v\n", sysError.Op, err)
	case log.DebugLevel:
		entry.Debugf("%s: %v\n", sysError.Op, err)
	default:
		entry.Errorf("%s: %v\n", sysError.Op, err)
	}
}
