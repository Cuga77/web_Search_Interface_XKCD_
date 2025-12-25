package core

import (
	"context"
)

type DB interface {
	Search(ctx context.Context, keywords []string, limit int) ([]Comic, int64, error)
	Scan(ctx context.Context) ([]Comic, error)
}

type Words interface {
	Norm(ctx context.Context, phrase string) ([]string, error)
}
