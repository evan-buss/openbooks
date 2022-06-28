package main

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"
	"net/http"

	"github.com/evan-buss/openbooks/cmd/service/util"
	openbooksv1 "github.com/evan-buss/openbooks/proto/gen/proto/go/openbooks/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var protectedServices = []string{"/openbooks.v1.AdminService/"}

func (s *server) StartGrpc() *grpc.Server {

	authzInterceptor := NewAuthZInterceptor("evan", protectedServices)

	grpcServer := grpc.NewServer(grpc_middleware.WithUnaryServerChain(grpcLogger(s.log), authzInterceptor.Unary()))
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

func grpcLogger(logger *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		logger.Infow("grpc request", "method", info.FullMethod)

		resp, err = handler(ctx, req)

		if err != nil {
			logger.Warnw("grpc request error", "err", err)
		}

		return
	}
}

func requireAuthorization(serviceName string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		return handler(ctx, req)
	}
}

func rateLimit() grpc.UnaryServerInterceptor {
	rl := util.NewRateLimiter()

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		p, ok := peer.FromContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unavailable, "Unable to find IP address.")
		}

		rl.Access(p.Addr.String())

		return handler(ctx, req)
	}
}
