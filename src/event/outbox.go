package event

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/inx51/howlite-resources/logger"
	_ "github.com/mattn/go-sqlite3"
)

type Outbox struct {
	mutex sync.Mutex
	db    *sql.DB
}

func NewOutbox(ctx context.Context, sqlitePath string) Outbox {
	setupSqliteFile(sqlitePath)
	db, err := sql.Open("sqlite3", sqlitePath)
	if err != nil {
		panic(err)
	}
	runMigrations(ctx, db)
	return Outbox{db: db}
}

func runMigrations(ctx context.Context, db *sql.DB) {
	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS outbox (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			payload TEXT NOT NULL,
			enqueued_utc TEXT NOT NULL
		)
	`); err != nil {
		_ = db.Close()
		panic(err)
	}
}

func setupSqliteFile(sqlitePath string) {
	if err := os.MkdirAll(filepath.Dir(sqlitePath), 0o755); err != nil {
		panic(err)
	}
}

func (outbox *Outbox) Enqueue(ctx context.Context, event []byte) {
	err := appendMessageToDb(ctx, outbox.db, event, time.Now().UTC())
	if err != nil {
		logger.Error(ctx, "failed to save to outbox", "error", err)
	}
}

func appendMessageToDb(ctx context.Context, db *sql.DB, event []byte, enqueuedUtc time.Time) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO outbox(payload, enqueued_utc) VALUES (?, ?)
	`, string(event), enqueuedUtc.Format(time.RFC3339Nano))
	return err
}

func (outbox *Outbox) Dequeue(ctx context.Context) []byte {
	outbox.mutex.Lock()
	defer outbox.mutex.Unlock()

	payload, err := getLatestMessageFromDb(ctx, outbox.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		logger.Error(ctx, "failed to dequeue event", "error", err)
		return nil
	}

	return payload
}

func getLatestMessageFromDb(ctx context.Context, db *sql.DB) ([]byte, error) {
	var payload string
	err := db.QueryRowContext(ctx, `
		DELETE FROM outbox
		WHERE id = (
			SELECT id
			FROM outbox
			ORDER BY id
			LIMIT 1
		)
		RETURNING payload
	`).Scan(&payload)
	if err != nil {
		return nil, err
	}

	return []byte(payload), nil
}

func (outbox *Outbox) Close(ctx context.Context) {
	if outbox == nil || outbox.db == nil {
		return
	}

	if err := outbox.db.Close(); err != nil {
		logger.Error(ctx, "failed to close outbox database", "error", err)
	}
}
