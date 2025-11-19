package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"comment-tree/internal/models"
	"comment-tree/internal/service"

	"github.com/wb-go/wbf/ginext"
)

type Handler struct {
	service *service.CommentService
}

func NewHandler(s *service.CommentService) *Handler {
	return &Handler{service: s}
}

// POST /comments
func (h *Handler) CreateComment(c *ginext.Context) {
	var req models.Comment
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	comment, err := h.service.CreateComment(c, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, comment)
}

// GET /comments?parent={id}&limit=&offset=&sort=
func (h *Handler) GetComments(c *ginext.Context) {
	// парсинг параметра parent в *int
	parent, err := parseParentParamFromQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	limit := queryInt(c, "limit", 20)
	offset := queryInt(c, "offset", 0)
	sort := c.DefaultQuery("sort", "ASC")

	ctx := c.Request.Context()

	tree, err := h.service.GetTree(ctx, parent, limit, offset, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tree)
}

// GET /comments/search?query=&limit=&offset=&sort=
func (h *Handler) SearchComments(c *ginext.Context) {
	q := c.Query("query")
	if q == "" {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "query is required"})
		return
	}

	limit := queryInt(c, "limit", 20)
	offset := queryInt(c, "offset", 0)
	sort := c.DefaultQuery("sort", "ASC")

	ctx := c.Request.Context()

	res, err := h.service.Search(ctx, q, limit, offset, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func parseParentParamFromQuery(c *ginext.Context) (*int, error) {
	v := c.Query("parent")
	if v == "" {
		return nil, nil
	}
	id, err := strconv.Atoi(v)
	if err != nil {
		return nil, fmt.Errorf("invalid parent: %v", err)
	}
	return &id, nil
}

func queryInt(c *ginext.Context, name string, def int) int {
	val := c.Query(name)
	if val == "" {
		return def
	}
	n, err := strconv.Atoi(val)
	if err != nil || n < 0 {
		return def
	}
	return n
}

// DELETE /comments/{id}
func (h *Handler) DeleteComment(c *ginext.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid id"})
		return
	}

	ctx := c.Request.Context()

	if err := h.service.DeleteComment(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
