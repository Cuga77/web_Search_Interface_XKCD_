package db

import (
	"context"
	"encoding/json"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"yadro.com/course/search/core"
)

type DB struct {
	log  *slog.Logger
	conn *sqlx.DB
}

func New(log *slog.Logger, address string) (*DB, error) {
	db, err := sqlx.Connect("pgx", address)
	if err != nil {
		log.Error("connection problem", "address", address, "error", err)
		return nil, err
	}
	return &DB{
		log:  log,
		conn: db,
	}, nil
}

func (db *DB) Search(ctx context.Context, keywords []string, limit int) ([]core.Comic, int64, error) {
	query := `
		SELECT ID, URL_ADRESS 
		FROM comics 
		WHERE WORDS ?| $1::text[] 
		ORDER BY (
			SELECT COUNT(*) 
			FROM jsonb_array_elements_text(WORDS) AS w 
			WHERE w = ANY($1::text[])
		) DESC
		LIMIT $2
	`

	db.log.Debug("executing search query", "query", query, "keywords", keywords, "limit", limit)
	rows, err := db.conn.QueryContext(ctx, query, keywords, limit)
	if err != nil {
		db.log.Error("failed to search comics", "error", err)
		return nil, 0, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			db.log.Error("failed to close rows", "error", err)
		}
	}()

	var comics []core.Comic
	for rows.Next() {
		var c core.Comic
		if err := rows.Scan(&c.ID, &c.URL); err != nil {
			db.log.Error("failed to scan comic", "error", err)
			continue
		}
		comics = append(comics, c)
	}
	db.log.Debug("search results", "count", len(comics))

	return comics, int64(len(comics)), nil
}

func (db *DB) Scan(ctx context.Context) ([]core.Comic, error) {
	query := `SELECT ID, URL_ADRESS, WORDS FROM comics`
	rows, err := db.conn.QueryContext(ctx, query)
	if err != nil {
		db.log.Error("failed to scan comics", "error", err)
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			db.log.Error("failed to close rows", "error", err)
		}
	}()

	var comics []core.Comic
	for rows.Next() {
		var c core.Comic
		var wordsBytes []byte
		if err := rows.Scan(&c.ID, &c.URL, &wordsBytes); err != nil {
			db.log.Error("failed to scan comic for index", "error", err)
			continue
		}

		var words []string
		if err := json.Unmarshal(wordsBytes, &words); err != nil {
			db.log.Error("failed to unmarshal words", "id", c.ID, "error", err)
		}

		c.Keywords = words
		comics = append(comics, c)
	}
	return comics, nil
}
