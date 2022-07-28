package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"strings"
	"time"
)

const metadataAuthKey = "token"

type AuthInterceptor struct {
	authKey          string
	protectedMethods []string
}

func NewAuthInterceptor(authKey string, protectedMethods []string) *AuthInterceptor {
	return &AuthInterceptor{
		authKey,
		protectedMethods,
	}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
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

		authValues := md[metadataAuthKey]
		if len(authValues) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "credentials not provided")
		}

		token, err := jwt.Parse([]byte(authValues[0]), jwt.WithKey(jwa.HS256, []byte(interceptor.authKey)))
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "credentials not valid")
		}

		log.Println(token)

		return handler(ctx, req)
	}
}

func ForwardTokenToCookie(ctx context.Context, writer http.ResponseWriter, message proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	token := md.HeaderMD.Get(metadataAuthKey)
	if len(token) == 0 {
		return nil
	}

	cookie := http.Cookie{
		Name:     "OpenBooksService",
		Value:    token[0],
		Secure:   false,
		HttpOnly: false,
		Expires:  time.Now().Add(time.Hour * 1),
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(writer, &cookie)

	return nil
}

func SetMetadataTokenFromCookie(ctx context.Context, request *http.Request) metadata.MD {
	cookie, err := request.Cookie("OpenBooksService")
	if err != nil {
		return nil
	}

	return metadata.Pairs("token", cookie.Value)
}
