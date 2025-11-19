package repository

import (
	"context"
	"fmt"

	"comment-tree/internal/models"

	"github.com/wb-go/wbf/dbpg"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *models.Comment) error
	GetByID(ctx context.Context, id int) (*models.Comment, error)
	GetTree(ctx context.Context, parentID *int, limit, offset int, sort string) ([]*models.Comment, error)
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, query string, limit, offset int, sort string) ([]*models.Comment, error)
}

type CommentRepo struct {
	DB *dbpg.DB
}

func NewCommentRepo(db *dbpg.DB) CommentRepository {
	return &CommentRepo{
		DB: db,
	}
}

func (r *CommentRepo) Create(ctx context.Context, c *models.Comment) error {
	query := `INSERT INTO comments (parent_id, author, text) 
	VALUES ($1, $2, $3) 
	RETURNING id, created_at`
	return r.DB.QueryRowContext(ctx, query, c.ParentID, c.Author, c.Text).Scan(&c.ID, &c.CreatedAt)
}

func (r *CommentRepo) GetByID(ctx context.Context, id int) (*models.Comment, error) {
	query := `SELECT id, parent_id, author, text, created_at FROM comments WHERE id = $1;`
	var c models.Comment
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.ParentID, &c.Author, &c.Text, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, err
}

func (r *CommentRepo) GetTree(ctx context.Context, parentID *int, limit, offset int, sort string) ([]*models.Comment, error) {
	var comments []*models.Comment

	if parentID == nil {
		// Только корневые id
		rootRows, err := r.DB.QueryContext(ctx, fmt.Sprintf(`
			SELECT id
			FROM comments
			WHERE parent_id IS NULL
			ORDER BY created_at %s
			LIMIT $1 OFFSET $2
		`, sort), limit, offset)
		if err != nil {
			return nil, err
		}
		defer rootRows.Close()

		var rootIDs []int
		for rootRows.Next() {
			var id int
			if err := rootRows.Scan(&id); err != nil {
				return nil, err
			}
			rootIDs = append(rootIDs, id)
		}

		// Если корневых комментариев нет возвращаем пустой список
		if len(rootIDs) == 0 {
			return []*models.Comment{}, nil
		}

		// Берём полные деревья для этих корней
		for _, rootID := range rootIDs {
			rows, err := r.DB.QueryContext(ctx, `
				WITH RECURSIVE tree AS (
					SELECT id, parent_id, author, text, created_at
					FROM comments
					WHERE id = $1
					UNION ALL
					SELECT c.id, c.parent_id, c.author, c.text, c.created_at
					FROM comments c
					JOIN tree t ON c.parent_id = t.id
				)
				SELECT id, parent_id, author, text, created_at
				FROM tree
				ORDER BY created_at ASC
			`, rootID)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			for rows.Next() {
				var c models.Comment
				if err := rows.Scan(&c.ID, &c.ParentID, &c.Author, &c.Text, &c.CreatedAt); err != nil {
					return nil, err
				}
				comments = append(comments, &c)
			}
		}

	} else {
		// Берём всё поддерево от parentID
		rows, err := r.DB.QueryContext(ctx, `
			WITH RECURSIVE tree AS (
				SELECT id, parent_id, author, text, created_at
				FROM comments
				WHERE id = $1
				UNION ALL
				SELECT c.id, c.parent_id, c.author, c.text, c.created_at
				FROM comments c
				JOIN tree t ON c.parent_id = t.id
			)
			SELECT id, parent_id, author, text, created_at
			FROM tree
			ORDER BY created_at ASC
		`, *parentID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var c models.Comment
			if err := rows.Scan(&c.ID, &c.ParentID, &c.Author, &c.Text, &c.CreatedAt); err != nil {
				return nil, err
			}
			comments = append(comments, &c)
		}
	}

	return comments, nil
}

func (r *CommentRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM comments WHERE id = $1;`
	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}

func (r *CommentRepo) Search(ctx context.Context, query string, limit, offset int, sort string) ([]*models.Comment, error) {
	rows, err := r.DB.QueryContext(ctx, fmt.Sprintf(`
        SELECT id, parent_id, author, text, created_at
        FROM comments
        WHERE to_tsvector('russian', text) @@ websearch_to_tsquery('russian', $1)
        ORDER BY created_at %s
        LIMIT $2 OFFSET $3
    `, sort), query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.ParentID, &c.Author, &c.Text, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}

	return comments, nil
}
