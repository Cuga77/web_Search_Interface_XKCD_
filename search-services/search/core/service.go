package core

import (
	"context"
	"fmt"
	"log/slog"
)

type Service struct {
	log   *slog.Logger
	db    DB
	words Words
	index *Index
}

func NewService(log *slog.Logger, db DB, words Words) *Service {
	return &Service{
		log:   log,
		db:    db,
		words: words,
		index: NewIndex(),
	}
}

func (s *Service) Search(ctx context.Context, phrase string, limit int) (SearchResult, error) {
	s.log.Debug("normalizing phrase", "phrase", phrase)
	keywords, err := s.words.Norm(ctx, phrase)
	if err != nil {
		return SearchResult{}, fmt.Errorf("failed to normalize phrase: %w", err)
	}

	if len(keywords) == 0 {
		return SearchResult{}, nil
	}

	s.log.Debug("searching comics", "keywords", keywords, "limit", limit)
	comics, total, err := s.db.Search(ctx, keywords, limit)
	if err != nil {
		return SearchResult{}, fmt.Errorf("failed to search comics: %w", err)
	}

	return SearchResult{
		Comics: comics,
		Total:  total,
	}, nil
}

func (s *Service) BuildIndex(ctx context.Context) error {
	s.log.Info("building index")
	comics, err := s.db.Scan(ctx)
	if err != nil {
		return fmt.Errorf("failed to scan comics: %w", err)
	}
	s.log.Info("scanned comics for index", "count", len(comics))

	s.index.Add(comics)
	s.log.Info("index built")
	return nil
}

func (s *Service) ISearch(ctx context.Context, phrase string, limit int) (SearchResult, error) {
	s.log.Debug("isearch: normalizing phrase", "phrase", phrase)
	keywords, err := s.words.Norm(ctx, phrase)
	if err != nil {
		return SearchResult{}, fmt.Errorf("failed to normalize phrase: %w", err)
	}

	if len(keywords) == 0 {
		return SearchResult{}, nil
	}

	s.log.Debug("isearch: searching index", "keywords", keywords)
	foundComics := s.index.Search(keywords)

	s.log.Debug("isearch: found comics", "count", len(foundComics))

	total := int64(len(foundComics))

	if limit > 0 && len(foundComics) > limit {
		foundComics = foundComics[:limit]
	}

	return SearchResult{
		Comics: foundComics,
		Total:  total,
	}, nil
}
