package repository

import (
	"context"

	. "shortener/internal/models"

	"github.com/wb-go/wbf/dbpg"
)

type ShortURLRepository interface {
	Save(ctx context.Context, url *ShortURL) error
	FindByID(ctx context.Context, ShortCode string) (*ShortURL, error)
	FindTopPopular(ctx context.Context, limit int) ([]ShortURL, error)
	ListLatest(ctx context.Context, limit int) ([]ShortURL, error)
}

type AnalyticsRepository interface {
	Save(ctx context.Context, event *ClickEvent) error
	GetStats(ctx context.Context, ShortCode string) ([]*ClickEvent, error)
}

type shortURLRepo struct {
	DB *dbpg.DB
}

type analyticsRepo struct {
	DB *dbpg.DB
}

func NewShortURLRepo(db *dbpg.DB) ShortURLRepository {
	return &shortURLRepo{
		DB: db,
	}
}

func NewAnalyticsRepo(db *dbpg.DB) AnalyticsRepository {
	return &analyticsRepo{
		DB: db,
	}
}

func (r *shortURLRepo) Save(ctx context.Context, u *ShortURL) error {
	query := `INSERT INTO short_urls (short_code, original) 
	VALUES ($1, $2) 
	RETURNING id, created_at`
	return r.DB.QueryRowContext(ctx, query, u.ShortCode, u.Original).Scan(&u.ID, &u.CreatedAt)
}

func (r *shortURLRepo) FindByID(ctx context.Context, code string) (*ShortURL, error) {
	query := `SELECT id, original FROM short_urls WHERE short_code = $1;`
	var u ShortURL
	err := r.DB.QueryRowContext(ctx, query, code).Scan(&u.ID, &u.Original)
	if err != nil {
		return nil, err
	}
	return &u, err
}

func (r *analyticsRepo) Save(ctx context.Context, u *ClickEvent) error {
	query := `INSERT INTO click_events (short_url_id, user_agent) VALUES ($1, $2)`
	_, err := r.DB.ExecContext(ctx, query, u.ShortID, u.UserAgent)
	return err
}

func (r *analyticsRepo) GetStats(ctx context.Context, shortCode string) ([]*ClickEvent, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT ce.id, ce.short_url_id, ce.user_agent, ce.timestamp
		FROM click_events ce
		JOIN short_urls su ON su.id = ce.short_url_id
		WHERE su.short_code = $1
		ORDER BY ce.timestamp DESC
	`, shortCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*ClickEvent
	for rows.Next() {
		var e ClickEvent
		if err := rows.Scan(&e.ID, &e.ShortID, &e.UserAgent, &e.Timestamp); err != nil {
			return nil, err
		}
		events = append(events, &e)
	}
	return events, nil
}

func (r *shortURLRepo) FindTopPopular(ctx context.Context, limit int) ([]ShortURL, error) {
	query := `
		SELECT 
			su.id,
			su.short_code,
			su.original,
			COUNT(ce.id) AS click_count
		FROM short_urls su
		LEFT JOIN click_events ce ON ce.short_url_id = su.id
		GROUP BY su.id
		ORDER BY click_count DESC
		LIMIT $1;
	`

	rows, err := r.DB.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []ShortURL
	for rows.Next() {
		var u ShortURL
		var clickCount int64
		if err := rows.Scan(&u.ID, &u.ShortCode, &u.Original, &clickCount); err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func (r *shortURLRepo) ListLatest(ctx context.Context, limit int) ([]ShortURL, error) {
	query := `
		SELECT id, short_code, original, created_at
		FROM short_urls
		ORDER BY created_at DESC
		LIMIT $1;
	`

	rows, err := r.DB.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ShortURL

	for rows.Next() {
		var s ShortURL
		if err := rows.Scan(&s.ID, &s.ShortCode, &s.Original, &s.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
