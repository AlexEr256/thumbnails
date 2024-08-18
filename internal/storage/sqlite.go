package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	dto "github.com/AlexEr256/thumbnail/internal/domain"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to db - %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS videos (link TEXT, url TEXT)`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create table - %w", err)
	}
	return nil
}
func (s *Storage) SaveVideos(videosList []dto.VideoInfo) error {
	sqlStr := "INSERT INTO videos(link, url) VALUES "
	vals := []interface{}{}

	for _, row := range videosList {
		sqlStr += "(?, ?),"
		vals = append(vals, row.Link, row.Url)
	}

	sqlStr = strings.TrimSuffix(sqlStr, ",")

	stmt, err := s.db.Prepare(sqlStr)
	if err != nil {
		return fmt.Errorf("failed to prepare statement for insert videos - %w", err)
	}
	_, err = stmt.Exec(vals...)
	if err != nil {
		return fmt.Errorf("failed to insert videos - %w", err)
	}
	return nil
}

func (s *Storage) ListVideos(videosList []string) ([]dto.VideoInfo, error) {
	videos := make([]dto.VideoInfo, 0)
	query := `select * from videos where link in (`

	for i, id := range videosList {
		if i != 0 {
			query += `, `
		}
		query += "'" + id + "'"
	}
	query += `);`

	rows, err := s.db.Query(query)

	if err != nil {
		return videos, fmt.Errorf("failed to select videos - %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var link string
		var url string
		err = rows.Scan(&link, &url)
		if err != nil {

			return videos, fmt.Errorf("failed to handle query result - %w", err)
		}
		videos = append(videos, dto.VideoInfo{Link: link, Url: url})
	}

	err = rows.Err()
	if err != nil {
		return videos, fmt.Errorf("failed to handle query result - %w", err)
	}

	return videos, nil

}
