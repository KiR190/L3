package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"warehouse/internal/auth"
	"warehouse/internal/models"

	"github.com/wb-go/wbf/dbpg"
)

type ItemRepository interface {
	Create(ctx context.Context, item *models.Item) error
	GetAll(ctx context.Context, limit, offset int) ([]models.Item, error)
	GetByID(ctx context.Context, id string) (*models.Item, error)
	Update(ctx context.Context, item *models.Item) error
	Delete(ctx context.Context, id string) error
	GetHistoryByItemID(ctx context.Context, itemID string, limit, offset int) ([]models.ItemHistory, error)
}

type Repository struct {
	db *dbpg.DB
}

func NewRepository(db *dbpg.DB) ItemRepository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, item *models.Item) error {
	userID, _ := auth.GetUserID(ctx)
	username, _ := auth.GetUserEmail(ctx)
	role, _ := auth.GetUserRole(ctx)

	query := `
		WITH _ctx AS (
			SELECT
				set_config('application.user_id', $1, true),
				set_config('application.username', $2, true),
				set_config('application.role', $3, true)
		)
		INSERT INTO items (id, sku, name, description, quantity, location)
		SELECT $4,$5,$6,$7,$8,$9 FROM _ctx
	`

	_, err := r.db.ExecContext(ctx, query,
		userID,
		username,
		role,
		item.ID,
		item.SKU,
		item.Name,
		item.Description,
		item.Quantity,
		item.Location,
	)
	return err
}

func (r *Repository) GetAll(ctx context.Context, limit, offset int) ([]models.Item, error) {
	query := `
	SELECT id, sku, name, description, quantity, location, created_at, updated_at   
	FROM items
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2                                                             
`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.Item, 0)
	for rows.Next() {
		var it models.Item
		err = rows.Scan(
			&it.ID,
			&it.SKU,
			&it.Name,
			&it.Description,
			&it.Quantity,
			&it.Location,
			&it.CreatedAt,
			&it.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, it)
	}

	return items, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*models.Item, error) {
	var it models.Item

	query := `
		SELECT id, sku, name, description, quantity, location, created_at, updated_at   -- изменено
		FROM items
		WHERE id=$1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&it.ID,
		&it.SKU,
		&it.Name,
		&it.Description,
		&it.Quantity,
		&it.Location,
		&it.CreatedAt,
		&it.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &it, nil
}

func (r *Repository) Update(ctx context.Context, item *models.Item) error {
	userID, _ := auth.GetUserID(ctx)
	username, _ := auth.GetUserEmail(ctx)
	role, _ := auth.GetUserRole(ctx)

	query := `
		WITH _ctx AS (
			SELECT
				set_config('application.user_id', $1, true),
				set_config('application.username', $2, true),
				set_config('application.role', $3, true)
		)
		UPDATE items
		SET sku=$4, name=$5, description=$6, quantity=$7, location=$8, updated_at=NOW()
		FROM _ctx
		WHERE id=$9
	`

	_, err := r.db.ExecContext(ctx, query,
		userID,
		username,
		role,
		item.SKU,
		item.Name,
		item.Description,
		item.Quantity,
		item.Location,
		item.ID,
	)
	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	userID, _ := auth.GetUserID(ctx)
	username, _ := auth.GetUserEmail(ctx)
	role, _ := auth.GetUserRole(ctx)

	query := `
		WITH _ctx AS (
			SELECT
				set_config('application.user_id', $1, true),
				set_config('application.username', $2, true),
				set_config('application.role', $3, true)
		)
		DELETE FROM items USING _ctx WHERE id=$4
	`

	_, err := r.db.ExecContext(ctx, query, userID, username, role, id)
	return err
}

func (r *Repository) GetHistoryByItemID(ctx context.Context, itemID string, limit, offset int) ([]models.ItemHistory, error) {
	query := `
		SELECT
			id,
			item_id,
			action,
			old_data::text,
			new_data::text,
			user_id,
			username,
			role,
			created_at
		FROM item_history
		WHERE item_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, itemID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.ItemHistory, 0)
	for rows.Next() {
		var h models.ItemHistory
		var oldText sql.NullString
		var newText sql.NullString
		err = rows.Scan(
			&h.ID,
			&h.ItemID,
			&h.Action,
			&oldText,
			&newText,
			&h.UserID,
			&h.Username,
			&h.Role,
			&h.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if oldText.Valid && oldText.String != "" {
			h.OldData = json.RawMessage([]byte(oldText.String))
		}
		if newText.Valid && newText.String != "" {
			h.NewData = json.RawMessage([]byte(newText.String))
		}

		out = append(out, h)
	}

	return out, nil
}
