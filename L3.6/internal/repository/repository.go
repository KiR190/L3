package repository

import (
	"context"

	"sales-tracker/internal/models"

	"github.com/wb-go/wbf/dbpg"
)

type ItemRepository interface {
	Create(ctx context.Context, it *models.Item) error
	GetByID(ctx context.Context, id string) (*models.Item, error)
	GetAll(ctx context.Context) ([]models.Item, error)
	Update(ctx context.Context, it *models.Item) error
	Delete(ctx context.Context, id string) error
	GetAnalytics(ctx context.Context, from, to string) (models.Analytics, error)
}

type Repository struct {
	db *dbpg.DB
}

func NewRepository(db *dbpg.DB) ItemRepository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, item *models.Item) error {
	query := `
		INSERT INTO items (id, type, amount, currency, category_id, note, occurred_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW())
	`

	_, err := r.db.ExecContext(ctx, query,
		item.ID, item.Type, item.Amount, item.Currency,
		item.CategoryID, item.Description, item.OccurredAt,
	)
	return err
}

func (r *Repository) GetAll(ctx context.Context) ([]models.Item, error) {
	query := `SELECT id, type, amount, currency, category_id, note, occurred_at, created_at, updated_at FROM items ORDER BY occurred_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var it models.Item
		err = rows.Scan(&it.ID, &it.Type, &it.Amount, &it.Currency, &it.CategoryID,
			&it.Description, &it.OccurredAt, &it.CreatedAt, &it.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, it)
	}

	return items, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*models.Item, error) {
	var it models.Item
	query := `SELECT id, type, amount, currency, category_id, note, occurred_at, created_at, updated_at FROM items WHERE id=$1`
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&it.ID, &it.Type, &it.Amount, &it.Currency, &it.CategoryID,
			&it.Description, &it.OccurredAt, &it.CreatedAt, &it.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &it, nil
}

func (r *Repository) Update(ctx context.Context, item *models.Item) error {
	query := `
		UPDATE items 
		SET type=$1, amount=$2, category_id=$3, note=$4, occurred_at=$5, updated_at=NOW()
		WHERE id=$6
	`
	_, err := r.db.ExecContext(ctx, query,
		item.Type, item.Amount, item.CategoryID, item.Description,
		item.OccurredAt, item.ID,
	)
	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM items WHERE id=$1`, id)
	return err
}

func (r *Repository) GetAnalytics(ctx context.Context, from, to string) (models.Analytics, error) {
	var result models.Analytics

	if from == "" {
		from = "0001-01-01"
	}
	if to == "" {
		to = "9999-12-31"
	}

	err := r.db.QueryRowContext(ctx, `
        SELECT
            COALESCE(SUM(amount),0),
            COALESCE(AVG(amount),0),
            COUNT(*),
            COALESCE(percentile_cont(0.5) WITHIN GROUP (ORDER BY amount), 0),
            COALESCE(percentile_cont(0.9) WITHIN GROUP (ORDER BY amount), 0)
        FROM items
        WHERE occurred_at >= $1 AND occurred_at <= $2
    `, from, to).Scan(&result.Sum, &result.Avg, &result.Count, &result.Median, &result.P90)

	return result, err
}
