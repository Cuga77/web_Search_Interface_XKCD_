package search

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"yadro.com/course/api/core"
	searchpb "yadro.com/course/proto/search"
)

type Client struct {
	log    *slog.Logger
	client searchpb.SearchClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client: %w", err)
	}
	return &Client{
		client: searchpb.NewSearchClient(conn),
		log:    log,
	}, nil
}

func (c *Client) Search(ctx context.Context, phrase string, limit int) (core.SearchResult, error) {
	req := &searchpb.SearchRequest{
		Phrase: phrase,
		Limit:  int32(limit),
	}
	resp, err := c.client.Search(ctx, req)
	if err != nil {
		c.log.Error("gRPC Search call failed", "error", err)
		return core.SearchResult{}, err
	}

	var comics []core.Comic
	for _, c := range resp.Comics {
		comics = append(comics, core.Comic{
			ID:  c.Id,
			URL: c.Url,
		})
	}
	return core.SearchResult{
		Comics: comics,
		Total:  resp.Total,
	}, nil
}

func (c *Client) ISearch(ctx context.Context, phrase string, limit int) (core.SearchResult, error) {
	req := &searchpb.SearchRequest{
		Phrase: phrase,
		Limit:  int32(limit),
	}
	resp, err := c.client.ISearch(ctx, req)
	if err != nil {
		c.log.Error("gRPC ISearch call failed", "error", err)
		return core.SearchResult{}, err
	}

	var comics []core.Comic
	for _, c := range resp.Comics {
		comics = append(comics, core.Comic{
			ID:  c.Id,
			URL: c.Url,
		})
	}
	return core.SearchResult{
		Comics: comics,
		Total:  resp.Total,
	}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.client.Ping(ctx, &emptypb.Empty{})
	return err
}
