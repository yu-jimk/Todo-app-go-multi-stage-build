package repository

import (
	"context"

	"myapp/internal/db"
	"myapp/internal/model"
)

type TodoRepository interface {
	ListTodos(ctx context.Context) ([]*model.Todo, error)
	GetTodo(ctx context.Context, id int64) (*model.Todo, error)
	CreateTodo(ctx context.Context, title string) (*model.Todo, error)
	UpdateTitle(ctx context.Context, id int64, title string) (*model.Todo, error)
	UpdateCompleted(ctx context.Context, id int64, completed bool) (*model.Todo, error)
	DeleteTodo(ctx context.Context, id int64) error
}

type todoRepository struct {
	queries *db.Queries
}

func NewTodoRepository(queries *db.Queries) TodoRepository {
	return &todoRepository{queries: queries}
}

func toModel(d db.Todo) *model.Todo {
	return &model.Todo{
		ID:        d.ID,
		Title:     d.Title,
		Completed: d.Completed,
		CreatedAt: d.CreatedAt.Time,
		UpdatedAt: d.UpdatedAt.Time,
	}
}

func (r *todoRepository) ListTodos(ctx context.Context) ([]*model.Todo, error) {
	// DBから取得 (db.Todo型)
	dbItems, err := r.queries.ListTodos(ctx)
	if err != nil {
		return nil, err
	}

	// model.Todo型へ変換
	items := make([]*model.Todo, len(dbItems))
	for i, d := range dbItems {
		items[i] = toModel(d)
	}
	return items, nil
}

func (r *todoRepository) GetTodo(ctx context.Context, id int64) (*model.Todo, error) {
	item, err := r.queries.GetTodo(ctx, id)
	if err != nil {
		return nil, err
	}
	return toModel(item), nil
}

func (r *todoRepository) CreateTodo(ctx context.Context, title string) (*model.Todo, error) {
	arg := db.CreateTodoParams{
		Title:     title,
		Completed: false,
	}
	item, err := r.queries.CreateTodo(ctx, arg)
	if err != nil {
		return nil, err
	}
	return toModel(item), nil
}

func (r *todoRepository) UpdateTitle(ctx context.Context, id int64, title string) (*model.Todo, error) {
	arg := db.UpdateTodoTitleParams{
		ID:    id,
		Title: title,
	}

	item, err := r.queries.UpdateTodoTitle(ctx, arg)
	if err != nil {
		return nil, err
	}
	return toModel(item), nil
}

func (r *todoRepository) UpdateCompleted(ctx context.Context, id int64, completed bool) (*model.Todo, error) {
	arg := db.UpdateTodoCompletedParams{
		ID:        id,
		Completed: completed,
	}

	item, err := r.queries.UpdateTodoCompleted(ctx, arg)
	if err != nil {
		return nil, err
	}
	return toModel(item), nil
}

func (r *todoRepository) DeleteTodo(ctx context.Context, id int64) error {
	// DELETE文を実行して、エラーがなければOKとするのが一般的です。
	// 存在しないIDをDELETEしてもDB的にはエラーにならないため、削除対象が存在するか確認してから消す必要はないです。
	return r.queries.DeleteTodo(ctx, id)
}
