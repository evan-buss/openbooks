package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

type AuthZInterceptor struct {
	password         string
	protectedMethods []string
}

func NewAuthZInterceptor(password string, protectedMethods []string) *AuthZInterceptor {
	return &AuthZInterceptor{password, protectedMethods}
}

func (interceptor *AuthZInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		// If the method doesn't need authentication, go next
		for _, method := range interceptor.protectedMethods {
			if method != info.FullMethod && !strings.HasPrefix(info.FullMethod, method) {
				return handler(ctx, req)
			}
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata not provided")
		}

		authValues := md["authorization"]
		if len(authValues) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "credentials not provided")
		}

		if interceptor.password != authValues[0] {
			return nil, status.Error(codes.Unauthenticated, "credentials not valid")
		}

		return handler(ctx, req)
	}
}
