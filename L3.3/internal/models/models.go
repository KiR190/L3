package models

import "time"

type Comment struct {
	ID        int       `json:"id"`
	ParentID  *int      `json:"parent_id"`
	Author    string    `json:"author"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentResponse struct {
	ID        int                `json:"id"`
	ParentID  *int               `json:"parent_id"`
	Author    string             `json:"author"`
	Text      string             `json:"text"`
	CreatedAt time.Time          `json:"created_at"`
	Children  []*CommentResponse `json:"children,omitempty"`
}
