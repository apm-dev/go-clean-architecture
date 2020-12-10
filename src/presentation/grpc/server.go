package grpc

import (
	"errors"
	"fmt"
	"github.com/apm-dev/go-clean-architecture/app/di"
	"github.com/apm-dev/go-clean-architecture/core/configs"
	errs "github.com/apm-dev/go-clean-architecture/core/errors"
	"github.com/apm-dev/go-clean-architecture/core/logger"
	"github.com/apm-dev/go-clean-architecture/presentation/grpc/interceptors"
	"github.com/apm-dev/go-clean-architecture/presentation/grpc/pbs"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type GRPC interface {
	Start() *errs.Error
	Stop()
}

var (
	instance   GRPC = &server{}
	grpcServer *grpc.Server
)

func GetInstance() GRPC {
	return instance
}

type server struct {
}

func (s *server) Start() *errs.Error {
	const op errs.Op = "grpc.server.Serve"
	fmt.Println("gRPC server starting...")

	conf, err := s.loadConfigs()
	if err != nil {
		return errs.E(op, err)
	}

	lis, err := s.makeListener(conf.ConnectionType, conf.Address)
	if err != nil {
		return errs.E(op, err)
	}

	grpcServer = s.makeServer()

	err = s.registerServices(grpcServer)
	if err != nil {
		return errs.E(op, err)
	}

	if configs.IsDebugMode() {
		fmt.Println("reflection registered")
		reflection.Register(grpcServer)
	}

	s.serve(grpcServer, lis)
	return nil
}

func (*server) Stop() {
	if grpcServer != nil {
		fmt.Println("Stopping gRPC server...")
		grpcServer.Stop()
	}
}

func (s *server) makeListener(ctype, addrs string) (net.Listener, *errs.Error) {
	const op errs.Op = "grpc.server.makeListener"
	lis, err := net.Listen(ctype, addrs)
	if err != nil {
		return nil, errs.E(op, errs.KindInternal, err)
	}
	return lis, nil
}

type grpcConf struct {
	ConnectionType string
	Address        string
}

func (s *server) loadConfigs() (*grpcConf, *errs.Error) {
	const op errs.Op = "grpc.server.loadConfigs"
	ctype := configs.Env("GRPC_CONNECTION_TYPE")
	addrs := configs.Env("GRPC_ADDRESS")
	if ctype == "" || addrs == "" {
		return nil, errs.E(op, errs.KindInternal,
			errors.New("can not find GRPC_CONNECTION_TYPE or GRPC_ADDRESS in ."+
				"env file"))
	}
	return &grpcConf{ConnectionType: ctype, Address: addrs}, nil
}

func (s *server) makeServer() *grpc.Server {
	return grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_validator.UnaryServerInterceptor(),
				interceptors.AuthFunc,
			),
		),
	)
}

func (s *server) registerServices(server *grpc.Server) *errs.Error {
	//const op errs.Op = "grpc.server.registerServices"
	pbs.RegisterBlogServiceServer(server, di.GrpcBlogServer())
	return nil
}

func (s *server) serve(server *grpc.Server, lis net.Listener) {
	const op errs.Op = "grpc.server.serve"
	go func() {
		if err := server.Serve(lis); err != nil {
			logger.SysError(errs.E(op, errs.KindInternal, err))
			panic("failed to serve grpc server")
		}
	}()
}
