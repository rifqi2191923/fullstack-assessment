package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/example/fullstack-assessment/backend/internal/model"
	"github.com/example/fullstack-assessment/backend/internal/repository"
	"github.com/redis/go-redis/v9"
)

type TaskService struct {
	Repo  *repository.TaskRepository
	Redis *redis.Client
	TTL   time.Duration
}

func NewTaskService(repo *repository.TaskRepository, rc *redis.Client, ttl time.Duration) *TaskService {
	return &TaskService{Repo: repo, Redis: rc, TTL: ttl}
}

func cacheKey(status, keyword string, assignee *uint64, page, limit int, sort string) string {
	a := ""
	if assignee != nil {
		a = fmt.Sprint(*assignee)
	}
	return fmt.Sprintf("tasks:v1:status=%s:keyword=%s:assignee=%s:page=%d:limit=%d:sort=%s", status, keyword, a, page, limit, sort)
}
func invalidate(ctx context.Context, rc *redis.Client) {
	if rc != nil {
		_ = rc.Del(ctx, "tasks:cache:index").Err()
	}
}

func (s *TaskService) List(ctx context.Context, status, keyword string, assignee *uint64, page, limit int, sort string) (model.TaskList, error) {
	key := cacheKey(status, keyword, assignee, page, limit, sort)
	if s.Redis != nil {
		if raw, err := s.Redis.Get(ctx, key).Result(); err == nil {
			var out model.TaskList
			if json.Unmarshal([]byte(raw), &out) == nil {
				return out, nil
			}
		}
	}
	out, err := s.Repo.List(ctx, status, keyword, assignee, page, limit, sort)
	if err != nil {
		return model.TaskList{}, err
	}
	if s.Redis != nil {
		if raw, err := json.Marshal(out); err == nil {
			_ = s.Redis.Set(ctx, key, raw, s.TTL).Err()
			_ = s.Redis.SAdd(ctx, "tasks:cache:index", key).Err()
		}
	}
	return out, nil
}
func (s *TaskService) Create(ctx context.Context, in model.TaskInput) (model.Task, error) {
	if err := repository.ValidateInput(in); err != nil {
		return model.Task{}, err
	}
	t, err := s.Repo.Create(ctx, in)
	if err == nil {
		s.Invalidate(ctx)
	}
	return t, err
}
func (s *TaskService) Update(ctx context.Context, id uint64, in model.TaskInput) (model.Task, error) {
	if err := repository.ValidateInput(in); err != nil {
		return model.Task{}, err
	}
	t, err := s.Repo.Update(ctx, id, in)
	if err == nil {
		s.Invalidate(ctx)
	}
	return t, err
}
func (s *TaskService) Delete(ctx context.Context, id uint64) error {
	err := s.Repo.SoftDelete(ctx, id)
	if err == nil {
		s.Invalidate(ctx)
	}
	return err
}
func (s *TaskService) Invalidate(ctx context.Context) {
	if s.Redis == nil {
		return
	}
	keys, err := s.Redis.SMembers(ctx, "tasks:cache:index").Result()
	if err == nil && len(keys) > 0 {
		_ = s.Redis.Del(ctx, keys...).Err()
	}
	invalidate(ctx, s.Redis)
}
