package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"myapp/internal/model"
)

type mockRepo struct {
	listResp            []*model.Todo
	getResp             *model.Todo
	createResp          *model.Todo
	updateTitleResp     *model.Todo
	updateCompletedResp *model.Todo
	deleteErr           error

	// capture args
	lastGetID             int64
	lastCreateTitle       string
	lastUpdateTitleID     int64
	lastUpdateTitle       string
	lastUpdateCompletedID int64
	lastUpdateCompleted   bool
	lastDeleteID          int64
}

func (m *mockRepo) ListTodos(ctx context.Context) ([]*model.Todo, error) {
	return m.listResp, nil
}
func (m *mockRepo) GetTodo(ctx context.Context, id int64) (*model.Todo, error) {
	m.lastGetID = id
	if m.getResp == nil {
		return nil, errors.New("not found")
	}
	return m.getResp, nil
}
func (m *mockRepo) CreateTodo(ctx context.Context, title string) (*model.Todo, error) {
	m.lastCreateTitle = title
	return m.createResp, nil
}
func (m *mockRepo) UpdateTitle(ctx context.Context, id int64, title string) (*model.Todo, error) {
	m.lastUpdateTitleID = id
	m.lastUpdateTitle = title
	return m.updateTitleResp, nil
}
func (m *mockRepo) UpdateCompleted(ctx context.Context, id int64, completed bool) (*model.Todo, error) {
	m.lastUpdateCompletedID = id
	m.lastUpdateCompleted = completed
	return m.updateCompletedResp, nil
}
func (m *mockRepo) DeleteTodo(ctx context.Context, id int64) error {
	m.lastDeleteID = id
	return m.deleteErr
}

func sampleTodo(id int64, title string) *model.Todo {
	now := time.Now().UTC()
	return &model.Todo{ID: id, Title: title, Completed: false, CreatedAt: now, UpdatedAt: now}
}

func TestListTodos(t *testing.T) {
	ctx := context.Background()
	m := &mockRepo{listResp: []*model.Todo{sampleTodo(1, "a"), sampleTodo(2, "b")}}
	s := NewTodoService(m)

	got, err := s.ListTodos(ctx)
	if err != nil {
		t.Fatalf("ListTodos error: %v", err)
	}
	if !reflect.DeepEqual(got, m.listResp) {
		t.Fatalf("unexpected result: got=%v want=%v", got, m.listResp)
	}
}

func TestGetTodo_Validation(t *testing.T) {
	ctx := context.Background()
	s := NewTodoService(&mockRepo{})
	if _, err := s.GetTodo(ctx, 0); err == nil {
		t.Fatalf("expected error for invalid id")
	}
}

func TestGetTodo_Delegates(t *testing.T) {
	ctx := context.Background()
	want := sampleTodo(5, "hello")
	m := &mockRepo{getResp: want}
	s := NewTodoService(m)
	got, err := s.GetTodo(ctx, 5)
	if err != nil {
		t.Fatalf("GetTodo error: %v", err)
	}
	if got.ID != want.ID || got.Title != want.Title {
		t.Fatalf("unexpected todo: got=%v want=%v", got, want)
	}
	if m.lastGetID != 5 {
		t.Fatalf("repo received wrong id: %d", m.lastGetID)
	}
}

func TestCreateTodo_Validation(t *testing.T) {
	ctx := context.Background()
	s := NewTodoService(&mockRepo{})
	if _, err := s.CreateTodo(ctx, ""); err == nil {
		t.Fatalf("expected error for empty title")
	}
}

func TestCreateTodo_Delegates(t *testing.T) {
	ctx := context.Background()
	want := sampleTodo(10, "create")
	m := &mockRepo{createResp: want}
	s := NewTodoService(m)
	got, err := s.CreateTodo(ctx, "create")
	if err != nil {
		t.Fatalf("CreateTodo error: %v", err)
	}
	if got.ID != want.ID {
		t.Fatalf("unexpected id: %d", got.ID)
	}
	if m.lastCreateTitle != "create" {
		t.Fatalf("repo received wrong title: %s", m.lastCreateTitle)
	}
}

func TestUpdateTitle_Validation(t *testing.T) {
	ctx := context.Background()
	s := NewTodoService(&mockRepo{})
	if _, err := s.UpdateTitle(ctx, 0, "x"); err == nil {
		t.Fatalf("expected error for invalid id")
	}
	if _, err := s.UpdateTitle(ctx, 1, ""); err == nil {
		t.Fatalf("expected error for empty title")
	}
}

func TestUpdateTitle_Delegates(t *testing.T) {
	ctx := context.Background()
	want := sampleTodo(3, "updated")
	m := &mockRepo{updateTitleResp: want}
	s := NewTodoService(m)
	got, err := s.UpdateTitle(ctx, 3, "updated")
	if err != nil {
		t.Fatalf("UpdateTitle error: %v", err)
	}
	if got.Title != "updated" {
		t.Fatalf("unexpected title: %v", got.Title)
	}
	if m.lastUpdateTitleID != 3 || m.lastUpdateTitle != "updated" {
		t.Fatalf("repo received wrong args: id=%d title=%s", m.lastUpdateTitleID, m.lastUpdateTitle)
	}
}

func TestUpdateCompleted_Validation(t *testing.T) {
	ctx := context.Background()
	s := NewTodoService(&mockRepo{})
	if _, err := s.UpdateCompleted(ctx, 0, true); err == nil {
		t.Fatalf("expected error for invalid id")
	}
}

func TestUpdateCompleted_Delegates(t *testing.T) {
	ctx := context.Background()
	want := sampleTodo(4, "c")
	want.Completed = true
	m := &mockRepo{updateCompletedResp: want}
	s := NewTodoService(m)
	got, err := s.UpdateCompleted(ctx, 4, true)
	if err != nil {
		t.Fatalf("UpdateCompleted error: %v", err)
	}
	if !got.Completed {
		t.Fatalf("expected completed=true")
	}
	if m.lastUpdateCompletedID != 4 || m.lastUpdateCompleted != true {
		t.Fatalf("repo received wrong args: id=%d completed=%v", m.lastUpdateCompletedID, m.lastUpdateCompleted)
	}
}

func TestDeleteTodo_Validation(t *testing.T) {
	ctx := context.Background()
	s := NewTodoService(&mockRepo{})
	if err := s.DeleteTodo(ctx, 0); err == nil {
		t.Fatalf("expected error for invalid id")
	}
}

func TestDeleteTodo_Delegates(t *testing.T) {
	ctx := context.Background()
	m := &mockRepo{deleteErr: nil}
	s := NewTodoService(m)
	if err := s.DeleteTodo(ctx, 7); err != nil {
		t.Fatalf("DeleteTodo error: %v", err)
	}
	if m.lastDeleteID != 7 {
		t.Fatalf("repo received wrong id: %d", m.lastDeleteID)
	}
}
