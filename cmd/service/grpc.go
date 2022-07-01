package main

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"
	"net/http"

	openbooksv1 "github.com/evan-buss/openbooks/proto/gen/proto/go/openbooks/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var protectedServices = []string{"/openbooks.AdminService/"}

func (s *server) StartGrpc() *grpc.Server {

	authInterceptor := NewAuthInterceptor("evan", "authorization", protectedServices)
	rateLimiter := NewRateLimiter()

	interceptors := grpc_middleware.WithUnaryServerChain(
		grpcLoggerUnary(s.log),
		rateLimiter.Unary(),
		authInterceptor.Unary(),
	)

	grpcServer := grpc.NewServer(interceptors)

	openbooksv1.RegisterOpenBooksServiceServer(grpcServer, s)
	openbooksv1.RegisterAdminServiceServer(grpcServer, s)

	if s.config.debug {
		reflection.Register(grpcServer)
	}

	return grpcServer
}

func (s *server) StartGrpcGateway() http.Handler {
	grpcGatewayMux := runtime.NewServeMux()
	credentials := insecure.NewCredentials()
	err := openbooksv1.RegisterOpenBooksServiceHandlerFromEndpoint(
		context.Background(),
		grpcGatewayMux,
		"localhost:"+s.config.port,
		[]grpc.DialOption{grpc.WithTransportCredentials(credentials)},
	)
	if err != nil {
		panic(err)
	}

	err = openbooksv1.RegisterAdminServiceHandlerFromEndpoint(
		context.Background(),
		grpcGatewayMux,
		"localhost:"+s.config.port,
		[]grpc.DialOption{grpc.WithTransportCredentials(credentials)},
	)

	if err != nil {
		panic(err)
	}

	return grpcGatewayMux
}

func grpcLoggerUnary(logger *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		logger.Infow("grpc request", "method", info.FullMethod)

		resp, err = handler(ctx, req)

		if err != nil {
			logger.Warnw("grpc request error", "err", err)
		}

		return
	}
}
