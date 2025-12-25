package db

import (
	"context"
	"encoding/json"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"yadro.com/course/update/core"
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

func (db *DB) Add(ctx context.Context, comics core.Comics) error {
	sqlStmt := `INSERT INTO comics (ID, URL_ADRESS, WORDS, TITLE, ALT, TRANSCRIPT, SAFE_TITLE) VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (ID) DO UPDATE SET
		URL_ADRESS = EXCLUDED.URL_ADRESS,
		WORDS = EXCLUDED.WORDS,
		TITLE = EXCLUDED.TITLE,
		ALT = EXCLUDED.ALT,
		TRANSCRIPT = EXCLUDED.TRANSCRIPT,
		SAFE_TITLE = EXCLUDED.SAFE_TITLE;`
	if comics.Words == nil {
		comics.Words = []string{}
	}
	wordsJSON, err := json.Marshal(comics.Words)
	if err != nil {
		db.log.Error("failed to marshal words to JSON", "error", err, "comic_id", comics.ID)
		return err
	}

	_, err = db.conn.ExecContext(ctx, sqlStmt, comics.ID, comics.URL, wordsJSON, comics.Title, comics.Alt, comics.Transcript, comics.SafeTitle)
	if err != nil {
		db.log.Error("failed to insert comic", "error", err, "comic_id", comics.ID)
		return err
	}

	return nil
}

func (db *DB) Stats(ctx context.Context) (core.DBStats, error) {
	const sqlStmt = `SELECT
    COUNT(*) AS comics_fetched,
    COALESCE(SUM(jsonb_array_length(WORDS)), 0) AS words_total,
    (SELECT COUNT(DISTINCT word) FROM comics, jsonb_array_elements_text(comics.words) AS word) AS words_unique
FROM
    comics;`
	var stats core.DBStats
	err := db.conn.GetContext(ctx, &stats, sqlStmt)
	if err != nil {
		db.log.Error("failed to get DB stats", "error", err)
		return core.DBStats{}, err
	}

	return stats, nil
}

func (db *DB) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	const query = "SELECT ID FROM comics"
	err := db.conn.SelectContext(ctx, &ids, query)
	if err != nil {
		db.log.Error("failed to fetch comic IDs", "error", err)
		return nil, err
	}

	return ids, nil
}

func (db *DB) Drop(ctx context.Context) error {
	const sqlStmt = `TRUNCATE TABLE comics`
	_, err := db.conn.ExecContext(ctx, sqlStmt)
	if err != nil {
		db.log.Error("failed to drop comics table", "error", err)
		return err
	}
	return nil
}
