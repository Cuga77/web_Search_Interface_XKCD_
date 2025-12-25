package core

import (
	"context"
)

type Normalizer interface {
	Norm(context.Context, string) ([]string, error)
}

type Pinger interface {
	Ping(context.Context) error
}

type Updater interface {
	Update(context.Context) error
	Stats(context.Context) (UpdateStats, error)
	Status(context.Context) (UpdateStatus, error)
	Drop(context.Context) error
}

type Searcher interface {
	Search(ctx context.Context, phrase string, limit int) (SearchResult, error)
	ISearch(ctx context.Context, phrase string, limit int) (SearchResult, error)
}

type DBStats struct {
	ComicsFetched int
	WordsTotal    int
	WordsUnique   int
}

type ComicsStorage interface {
	SaveComics(context.Context, Comics) error
	GetDBStats(context.Context) (DBStats, error)
	GetSavedComicIDs(context.Context) ([]int, error)
	DropAll(context.Context) error
}
