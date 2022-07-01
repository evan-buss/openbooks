package main

import (
	"context"
	"time"

	openbooksv1 "github.com/evan-buss/openbooks/proto/gen/proto/go/openbooks/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Search finds books that match the query
func (s *server) Search(ctx context.Context, req *openbooksv1.SearchRequest) (*openbooksv1.SearchResponse, error) {
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
func (s *server) StartIndex(ctx context.Context, req *openbooksv1.StartIndexRequest) (*openbooksv1.StartIndexResponse, error) {
	s.log.Infow("started re-indexing")

	duration, err := s.searcher.AddDocuments(ctx, s.indexer.Index(ctx))
	if err != nil {
		s.log.Errorw("indexing error", "err", err)
		return nil, status.Error(codes.Internal, "Error occurred during indexing.")
	}

	s.log.Infow("completed indexing", "duration", duration.Milliseconds())

	return &openbooksv1.StartIndexResponse{}, nil
}
