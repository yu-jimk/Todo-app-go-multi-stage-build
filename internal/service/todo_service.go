package service

import (
	"context"
	"errors"

	"myapp/internal/model"
	"myapp/internal/repository"
)

type TodoService interface {
	// 戻り値をポインタ (*model.Todo) にすると、nilチェックなどがしやすくなります
	ListTodos(ctx context.Context) ([]*model.Todo, error)
	GetTodo(ctx context.Context, id int64) (*model.Todo, error)
	CreateTodo(ctx context.Context, title string) (*model.Todo, error)
	UpdateTitle(ctx context.Context, id int64, title string) (*model.Todo, error)
	UpdateCompleted(ctx context.Context, id int64, completed bool) (*model.Todo, error)
	DeleteTodo(ctx context.Context, id int64) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(r repository.TodoRepository) TodoService {
	return &todoService{repo: r}
}

func (s *todoService) ListTodos(ctx context.Context) ([]*model.Todo, error) {
	// 変換はRepoが終わらせているので、そのまま返すだけ
	return s.repo.ListTodos(ctx)
}

func (s *todoService) GetTodo(ctx context.Context, id int64) (*model.Todo, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	return s.repo.GetTodo(ctx, id)
}

func (s *todoService) CreateTodo(ctx context.Context, title string) (*model.Todo, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	return s.repo.CreateTodo(ctx, title)
}

func (s *todoService) UpdateTitle(ctx context.Context, id int64, title string) (*model.Todo, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	if title == "" {
		return nil, errors.New("title is required")
	}
	return s.repo.UpdateTitle(ctx, id, title)
}

func (s *todoService) UpdateCompleted(ctx context.Context, id int64, completed bool) (*model.Todo, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	return s.repo.UpdateCompleted(ctx, id, completed)
}

func (s *todoService) DeleteTodo(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	return s.repo.DeleteTodo(ctx, id)
}
