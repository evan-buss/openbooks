package main

import (
	"context"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	openbooks "github.com/evan-buss/openbooks/proto/gen/proto/go/openbooks/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Search finds books that match the query
func (s *server) Search(ctx context.Context, req *openbooks.SearchRequest) (*openbooks.SearchResponse, error) {
	if req.GetQuery() == "" {
		return nil, status.Error(codes.InvalidArgument, "Query is required.")
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	response, err := s.searcher.Search(timeoutCtx, req.GetQuery())
	if err != nil {
		s.log.Errorw("search error", "err", err)
		return nil, status.Error(codes.Internal, "Error finding matches for search query.")
	}

	s.log.Infow("handled search query", "query", req.GetQuery(), "duration", response.Duration.AsDuration().Milliseconds())
	return response, nil
}

// StartIndex reads all documents from the source directory, parse them, and add them to the search index
func (s *server) StartIndex(ctx context.Context, req *openbooks.StartIndexRequest) (*openbooks.StartIndexResponse, error) {
	// s.log.Infow("started re-indexing")
	// return &openbooks.StartIndexResponse{}, nil

	duration, err := s.searcher.AddDocuments(ctx, s.indexer.Index(ctx))
	if err != nil {
		s.log.Errorw("indexing error", "err", err)
		return nil, status.Error(codes.Internal, "Error occurred during indexing.")
	}

	s.log.Infow("completed indexing", "duration", duration.Milliseconds())

	return &openbooks.StartIndexResponse{}, nil
}

func (s *server) LogIn(ctx context.Context, req *openbooks.LogInRequest) (*openbooks.LogInResponse, error) {
	if !s.store.Authenticate(req.GetUsername(), req.GetPassword()) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	tok, err := jwt.NewBuilder().
		Issuer("github.com/evan-buss/openbooks/cmd/service").
		Expiration(time.Now().Add(1*time.Hour)).
		Claim("username", req.GetUsername()).
		Build()

	if err != nil {
		s.log.Errorw("unable to create JWT")
	}

	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.HS256, []byte(s.config.Key)))
	if err != nil {
		s.log.Errorw("unable to sign JWT")
	}

	grpc.SetHeader(ctx, metadata.Pairs("token", string(signed)))
	return &openbooks.LogInResponse{}, nil
}
