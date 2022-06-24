package main

import (
	"context"

	openbooksv1 "github.com/evan-buss/openbooks/proto/gen/proto/go/openbooks/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Search finds books that match the query
func (s *server) Search(ctx context.Context, req *openbooksv1.SearchRequest) (*openbooksv1.SearchResponse, error) {
	if req.GetQuery() == "" {
		return nil, status.Error(codes.InvalidArgument, "Query is required.")
	}

	response := s.searcher.Search(req.GetQuery())

	return &response, nil
}
