package core

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type Service struct {
	log         *slog.Logger
	db          DB
	xkcd        XKCD
	words       Words
	eb          EventBus
	concurrency int

	mu     sync.Mutex
	status ServiceStatus
}

func NewService(
	log *slog.Logger, db DB, xkcd XKCD, words Words, eb EventBus, concurrency int,
) (*Service, error) {
	if concurrency < 1 {
		return nil, fmt.Errorf("wrong concurrency specified: %d", concurrency)
	}
	return &Service{
		log:         log,
		db:          db,
		xkcd:        xkcd,
		words:       words,
		eb:          eb,
		concurrency: concurrency,
		status:      StatusIdle,
	}, nil
}

func (s *Service) Update(ctx context.Context) (err error) {
	s.mu.Lock()
	if s.status == StatusRunning {
		s.mu.Unlock()
		return ErrUpdateInProgress
	}

	s.status = StatusRunning
	s.mu.Unlock()

	s.log.Info("starting comics update")
	defer func() {
		s.mu.Lock()
		s.status = StatusIdle
		s.mu.Unlock()
		s.log.Info("update finished")
	}()

	maxID, err := s.xkcd.LastID(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch latest comic: %w", err)
	}

	savedIDs, err := s.db.IDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get saved IDs: %w", err)
	}

	savedIDsMap := make(map[int]struct{}, len(savedIDs))
	for _, id := range savedIDs {
		savedIDsMap[id] = struct{}{}
	}

	var comicsToFetch []int
	for id := 1; id <= maxID; id++ {
		if id == 404 {
			continue
		}
		if _, exists := savedIDsMap[id]; !exists {
			comicsToFetch = append(comicsToFetch, id)
		}
	}

	if len(comicsToFetch) == 0 {
		s.log.Info("no new comics to fetch, database is up to date")
		return nil
	}
	s.log.Info("comics to fetch", "count", len(comicsToFetch))

	var wg sync.WaitGroup
	jobs := make(chan int, s.concurrency)

	for i := 0; i < s.concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range jobs {
				var comicData XKCDInfo
				var err error
				for attempt := 0; attempt < 10; attempt++ {
					comicData, err = s.xkcd.Get(ctx, id)
					if err == nil {
						break
					}
					s.log.Warn("failed to fetch comic, retrying", "id", id, "attempt", attempt+1, "error", err)
					select {
					case <-ctx.Done():
						return
					case <-time.After(1 * time.Second):
					}
				}
				if err != nil {
					s.log.Error("failed to fetch comic after retries", "id", id, "error", err)
					continue
				}

				textToNorm := comicData.Alt + " " + comicData.Title + " " + comicData.Transcript

				var keywords []string

				for attempt := 0; attempt < 10; attempt++ {
					keywords, err = s.words.Norm(ctx, textToNorm)
					if err == nil {
						break
					}
					s.log.Warn("failed to normalize words, retrying", "id", id, "attempt", attempt+1, "error", err)
					select {
					case <-ctx.Done():
						return
					case <-time.After(1 * time.Second):
					}
				}
				if err != nil {
					s.log.Error("failed to normalize words after retries", "id", id, "error", err)
					keywords = []string{}
				}

				comicToSave := Comics{
					ID:         comicData.ID,
					URL:        comicData.URL,
					Words:      keywords,
					Title:      comicData.Title,
					Alt:        comicData.Alt,
					Transcript: comicData.Transcript,
					SafeTitle:  comicData.SafeTitle,
				}

				err = s.db.Add(ctx, comicToSave)
				if err != nil {
					s.log.Warn("failed to save comic", "id", id, "error", err)
					continue
				}
				s.log.Debug("successfully saved comic", "id", id)
			}
		}()
	}

	for _, id := range comicsToFetch {
		jobs <- id
	}
	close(jobs)
	wg.Wait()

	if err := s.eb.PublishUpdate(); err != nil {
		s.log.Error("failed to publish update event", "error", err)
	}

	return nil
}

func (s *Service) Stats(ctx context.Context) (ServiceStats, error) {
	s.log.Debug("getting stats")

	dbStats, err := s.db.Stats(ctx)
	if err != nil {
		s.log.Error("failed to get db stats", "error", err)
		return ServiceStats{}, err
	}

	maxID, err := s.xkcd.LastID(ctx)
	if err != nil {
		s.log.Error("failed to fetch latest comic from xkcd", "error", err)
		return ServiceStats{}, err
	}

	if maxID >= 404 {
		maxID--
	}

	stats := ServiceStats{
		DBStats:     dbStats,
		ComicsTotal: maxID,
	}

	return stats, nil
}

func (s *Service) Status(ctx context.Context) ServiceStatus {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.status
}

func (s *Service) Drop(ctx context.Context) error {
	s.log.Debug("dropping all comics")
	err := s.db.Drop(ctx)
	if err != nil {
		s.log.Error("failed to drop comics", "error", err)
		return err
	}
	if err := s.eb.PublishUpdate(); err != nil {
		s.log.Error("failed to publish update event", "error", err)
	}
	return nil
}
