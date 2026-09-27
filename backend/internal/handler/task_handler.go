package handler

import (
	"errors"
	"strconv"
	"strings"

	"github.com/example/fullstack-assessment/backend/internal/model"
	"github.com/example/fullstack-assessment/backend/internal/repository"
	"github.com/example/fullstack-assessment/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct{ Service *service.TaskService }

func NewTaskHandler(s *service.TaskService) *TaskHandler { return &TaskHandler{Service: s} }

func errorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": gin.H{"code": status, "message": message}})
}
func parseID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, 400, "invalid task id")
		return 0, false
	}
	return id, true
}

func (h *TaskHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	var assignee *uint64
	if v := c.Query("assignee"); v != "" {
		n, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			errorResponse(c, 400, "invalid assignee")
			return
		}
		assignee = &n
	}
	out, err := h.Service.List(c.Request.Context(), c.Query("status"), c.Query("keyword"), assignee, page, limit, c.DefaultQuery("sort", "created_at_desc"))
	if err != nil {
		errorResponse(c, 500, "failed to fetch tasks")
		return
	}
	c.JSON(200, out)
}
func (h *TaskHandler) Create(c *gin.Context) {
	var in model.TaskInput
	if err := c.ShouldBindJSON(&in); err != nil {
		errorResponse(c, 400, "invalid request body")
		return
	}
	t, err := h.Service.Create(c.Request.Context(), in)
	if errors.Is(err, repository.ErrDuplicateTitle) {
		errorResponse(c, 409, "task title already exists")
		return
	}
	if err != nil {
		errorResponse(c, 400, err.Error())
		return
	}
	c.JSON(201, t)
}
func (h *TaskHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in model.TaskInput
	if err := c.ShouldBindJSON(&in); err != nil {
		errorResponse(c, 400, "invalid request body")
		return
	}
	t, err := h.Service.Update(c.Request.Context(), id, in)
	if errors.Is(err, repository.ErrNotFound) {
		errorResponse(c, 404, "task not found")
		return
	}
	if errors.Is(err, repository.ErrDuplicateTitle) {
		errorResponse(c, 409, "task title already exists")
		return
	}
	if err != nil {
		errorResponse(c, 400, err.Error())
		return
	}
	c.JSON(200, t)
}
func (h *TaskHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	err := h.Service.Delete(c.Request.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		errorResponse(c, 404, "task not found")
		return
	}
	if err != nil {
		errorResponse(c, 500, "failed to delete task")
		return
	}
	c.JSON(200, gin.H{"message": "task deleted"})
}

func Register(r *gin.RouterGroup, h *TaskHandler) {
	r.GET("/tasks", h.List)
	r.POST("/tasks", h.Create)
	r.PUT("/tasks/:id", h.Update)
	r.DELETE("/tasks/:id", h.Delete)
}

var _ = strings.TrimSpace
