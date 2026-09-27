package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/example/fullstack-assessment/backend/internal/model"
)

var ErrNotFound = errors.New("task not found")
var ErrDuplicateTitle = errors.New("task title already exists")

const baseSelect = `SELECT id, title, COALESCE(description,''), status, assignee_id, created_at, updated_at, deleted_at FROM tasks`

// TaskRepository contains the MySQL persistence layer.
type TaskRepository struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

func (r *TaskRepository) List(
	ctx context.Context,
	status, keyword string,
	assignee *uint64,
	page, limit int,
	sort string,
) (model.TaskList, error) {
	where := []string{"deleted_at IS NULL"}
	args := make([]any, 0)

	if status != "" {
		where = append(where, "status = ?")
		args = append(args, status)
	}

	if keyword != "" {
		where = append(where, "(title LIKE ? OR description LIKE ?)")
		k := "%" + keyword + "%"
		args = append(args, k, k)
	}

	if assignee != nil {
		where = append(where, "assignee_id = ?")
		args = append(args, *assignee)
	}

	order := "created_at DESC"

	switch sort {
	case "created_at_asc":
		order = "created_at ASC"
	case "updated_at_desc":
		order = "updated_at DESC"
	case "updated_at_asc":
		order = "updated_at ASC"
	case "title_asc":
		order = "title ASC"
	case "title_desc":
		order = "title DESC"
	}

	whereSQL := " WHERE " + strings.Join(where, " AND ")

	var total int

	if err := r.DB.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM tasks"+whereSQL,
		args...,
	).Scan(&total); err != nil {
		return model.TaskList{}, err
	}

	offset := (page - 1) * limit

	query := baseSelect +
		whereSQL +
		" ORDER BY " +
		order +
		" LIMIT ? OFFSET ?"

	queryArgs := append(args, limit, offset)

	rows, err := r.DB.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return model.TaskList{}, err
	}
	defer rows.Close()

	items := make([]model.Task, 0)

	for rows.Next() {
		var t model.Task

		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.AssigneeID,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.DeletedAt,
		); err != nil {
			return model.TaskList{}, err
		}

		items = append(items, t)
	}

	if err := rows.Err(); err != nil {
		return model.TaskList{}, err
	}

	totalPages := 0

	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return model.TaskList{
		Items:      items,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (r *TaskRepository) Get(ctx context.Context, id uint64) (model.Task, error) {
	var t model.Task

	err := r.DB.QueryRowContext(
		ctx,
		baseSelect+" WHERE id = ? AND deleted_at IS NULL",
		id,
	).Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.AssigneeID,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.DeletedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, ErrNotFound
	}

	return t, err
}

func (r *TaskRepository) Create(
	ctx context.Context,
	in model.TaskInput,
) (model.Task, error) {

	// Check apakah masih ada task aktif dengan title yang sama.
	var existingID uint64

	err := r.DB.QueryRowContext(
		ctx,
		"SELECT id FROM tasks WHERE title = ? AND deleted_at IS NULL LIMIT 1",
		in.Title,
	).Scan(&existingID)

	if err == nil {
		return model.Task{}, ErrDuplicateTitle
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, err
	}

	// Jika tidak ada duplicate, buat task baru.
	res, err := r.DB.ExecContext(
		ctx,
		"INSERT INTO tasks (title, description, status, assignee_id) VALUES (?, ?, ?, ?)",
		in.Title,
		in.Description,
		in.Status,
		in.AssigneeID,
	)

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return model.Task{}, ErrDuplicateTitle
		}

		return model.Task{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return model.Task{}, err
	}

	return r.Get(ctx, uint64(id))
}

func (r *TaskRepository) Update(
	ctx context.Context,
	id uint64,
	in model.TaskInput,
) (model.Task, error) {

	// Check apakah title sudah digunakan oleh task aktif lain.
	var existingID uint64

	err := r.DB.QueryRowContext(
		ctx,
		"SELECT id FROM tasks WHERE title = ? AND deleted_at IS NULL AND id <> ? LIMIT 1",
		in.Title,
		id,
	).Scan(&existingID)

	if err == nil {
		return model.Task{}, ErrDuplicateTitle
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, err
	}

	// Update task.
	res, err := r.DB.ExecContext(
		ctx,
		"UPDATE tasks SET title=?, description=?, status=?, assignee_id=? WHERE id=? AND deleted_at IS NULL",
		in.Title,
		in.Description,
		in.Status,
		in.AssigneeID,
		id,
	)

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return model.Task{}, ErrDuplicateTitle
		}

		return model.Task{}, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return model.Task{}, err
	}

	if n == 0 {
		return model.Task{}, ErrNotFound
	}

	return r.Get(ctx, id)
}

func (r *TaskRepository) SoftDelete(
	ctx context.Context,
	id uint64,
) error {

	res, err := r.DB.ExecContext(
		ctx,
		"UPDATE tasks SET deleted_at = NOW() WHERE id=? AND deleted_at IS NULL",
		id,
	)

	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *TaskRepository) ResetForTest(ctx context.Context) error {
	_, err := r.DB.ExecContext(ctx, "DELETE FROM tasks")
	return err
}

func ValidateInput(in model.TaskInput) error {
	if strings.TrimSpace(in.Title) == "" {
		return fmt.Errorf("title is required")
	}

	if in.Status != "todo" &&
		in.Status != "in_progress" &&
		in.Status != "done" {
		return fmt.Errorf("invalid status")
	}

	return nil
}