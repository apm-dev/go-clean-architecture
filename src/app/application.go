package app

import (
	"fmt"
	"github.com/apm-dev/go-clean-architecture/app/di"
	"github.com/apm-dev/go-clean-architecture/core/logger"
	"github.com/apm-dev/go-clean-architecture/presentation/grpc"
)

func StartApplication() {
	//	Migrate DBs tables
	migratePostgres()

	//	Start grpc server
	startGRPC()

}

func StopApplication() {
	//	Stop grpc server
	grpc.GetInstance().Stop()
	fmt.Println("gRPC server stopped")
}

func startGRPC() {
	err := grpc.GetInstance().Start()
	if err != nil {
		logger.SysError(err)
		panic("failed start grpc server")
	}
	fmt.Println("gRPC server started")
}

func migratePostgres() {
	var tables []string
	//	get existing tables
	_, pgmErr := di.PostgresDB().Query(&tables, `
		SELECT tablename
		FROM pg_catalog.pg_tables
		WHERE schemaname = 'public'
	`)
	if pgmErr != nil {
		logger.SysError(pgmErr)
		panic("could not get postgres tables")
	}
	//	if there was no tables migration will run
	if len(tables) == 0 {
		pgdsErr := di.PostgresDS().CreateTables()
		if pgdsErr != nil {
			logger.SysError(pgmErr)
			panic("could not create postgres tables")
		}
	}
}
