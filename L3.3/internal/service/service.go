package service

import (
	"context"

	"comment-tree/internal/models"
	"comment-tree/internal/repository"
)

type CommentService struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateComment(ctx context.Context, req *models.Comment) (*models.Comment, error) {
	comment := &models.Comment{
		ParentID: req.ParentID,
		Author:   req.Author,
		Text:     req.Text,
	}
	err := s.repo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *CommentService) GetComment(ctx context.Context, id int) (*models.Comment, error) {
	comment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) Search(ctx context.Context, q string, limit, offset int, sort string) ([]*models.CommentResponse, error) {
	comments, err := s.repo.Search(ctx, q, limit, offset, sort)
	if err != nil {
		return nil, err
	}

	return BuildTree(comments, nil), nil
}

func (s *CommentService) GetTree(ctx context.Context, parentID *int, limit, offset int, sort string) ([]*models.CommentResponse, error) {
	comments, err := s.repo.GetTree(ctx, parentID, limit, offset, sort)
	if err != nil {
		return nil, err
	}

	return BuildTree(comments, parentID), nil
}

func BuildTree(comments []*models.Comment, parentID *int) []*models.CommentResponse {
	// индекс по ID
	lookup := make(map[int]*models.CommentResponse)

	for _, c := range comments {
		lookup[c.ID] = &models.CommentResponse{
			ID:        c.ID,
			ParentID:  c.ParentID,
			Author:    c.Author,
			Text:      c.Text,
			CreatedAt: c.CreatedAt,
			Children:  []*models.CommentResponse{},
		}
	}

	// итоговые корни
	var roots []*models.CommentResponse

	// строим дерево
	for _, c := range comments {
		node := lookup[c.ID]

		if parentID != nil && c.ID == *parentID {
			roots = append(roots, node)
		} else if c.ParentID == nil {
			roots = append(roots, node)
		} else {
			parent := lookup[*c.ParentID]
			if parent != nil {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return roots
}

// Удаление комментария и всех вложенных
func (s *CommentService) DeleteComment(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
