package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

type Store struct {
	db   *sql.DB
	path string
}

func Open(embedded []byte) (*Store, error) {
	f, err := os.CreateTemp("", "lexigen-*.sqlite")
	if err != nil {
		return nil, fmt.Errorf("create temp database: %w", err)
	}
	path := f.Name()
	if _, err := f.Write(embedded); err != nil {
		f.Close()
		os.Remove(path)
		return nil, fmt.Errorf("write temp database: %w", err)
	}
	if err := f.Close(); err != nil {
		os.Remove(path)
		return nil, fmt.Errorf("close temp database: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		os.Remove(path)
		return nil, fmt.Errorf("open database: %w", err)
	}

	return &Store{db: db, path: path}, nil
}

func (s *Store) Close() error {
	err := s.db.Close()
	if rmErr := os.Remove(s.path); err == nil {
		err = rmErr
	}
	return err
}

func (s *Store) MaxPhrases() (int, error) {
	var max int
	err := s.db.QueryRow(`
		SELECT MIN(
			(SELECT COUNT(*) FROM adjectives),
			(SELECT COUNT(*) FROM nouns)
		)`).Scan(&max)
	if err != nil {
		return 0, fmt.Errorf("count words: %w", err)
	}
	return max, nil
}

const phrasesQuery = `
WITH
    A AS (
        SELECT
            ROW_NUMBER() OVER (ORDER BY rnd) AS num,
            word
        FROM (
            SELECT word, RANDOM() AS rnd
            FROM adjectives
            ORDER BY rnd
            LIMIT ?
        )
    ),
    N AS (
        SELECT
            ROW_NUMBER() OVER (ORDER BY rnd) AS num,
            word
        FROM (
            SELECT word, RANDOM() AS rnd
            FROM nouns
            ORDER BY rnd
            LIMIT ?
        )
    )
SELECT A.word, N.word
FROM A, N
WHERE A.num = N.num
`

func (s *Store) Phrases(n int, delimiter string) ([]string, error) {
	rows, err := s.db.Query(phrasesQuery, n, n)
	if err != nil {
		return nil, fmt.Errorf("query phrases: %w", err)
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var adj, noun string
		if err := rows.Scan(&adj, &noun); err != nil {
			return nil, fmt.Errorf("scan phrase: %w", err)
		}
		out = append(out, adj+delimiter+noun)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read phrases: %w", err)
	}
	return out, nil
}
