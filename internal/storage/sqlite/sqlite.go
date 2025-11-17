package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"urlshortener/internal/storage"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %s", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS urls(
	id INTEGER PRIMARY KEY,
	alias TEXT NOT NULL UNIQUE,
	url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON urls(alias);`)
	if err != nil {
		return nil, fmt.Errorf("%s : %s", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s : %s", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(targetUrl, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"
	stmt, err := s.db.Prepare(`INSERT INTO urls(url, alias) VALUES(?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("%s: %s", op, err)
	}

	result, err := stmt.Exec(targetUrl, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s : %s", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: %s", op, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last inserted id: %s", op, err)
	}
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"
	stmt, err := s.db.Prepare(`SELECT url FROM urls WHERE alias = ?`)
	if err != nil {
		return "", fmt.Errorf("%s: %s", op, err)
	}
	var resURL string
	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: %s", op, err)
	}
	return resURL, nil
}
