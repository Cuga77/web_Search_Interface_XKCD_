package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"yadro.com/course/proto/search"
	"yadro.com/course/search/core"
)

type Server struct {
	search.UnimplementedSearchServer
	service *core.Service
}

func NewServer(service *core.Service) *Server {
	return &Server{service: service}
}

func (s *Server) Search(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error) {
	res, err := s.service.Search(ctx, req.Phrase, int(req.Limit))
	if err != nil {
		return nil, err
	}

	var comics []*search.Comic
	for _, c := range res.Comics {
		comics = append(comics, &search.Comic{
			Id:  c.ID,
			Url: c.URL,
		})
	}

	return &search.SearchResponse{
		Comics: comics,
		Total:  res.Total,
	}, nil
}

func (s *Server) Ping(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *Server) ISearch(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error) {
	res, err := s.service.ISearch(ctx, req.Phrase, int(req.Limit))
	if err != nil {
		return nil, err
	}

	var comics []*search.Comic
	for _, c := range res.Comics {
		comics = append(comics, &search.Comic{
			Id:  c.ID,
			Url: c.URL,
		})
	}

	return &search.SearchResponse{
		Comics: comics,
		Total:  res.Total,
	}, nil
}
