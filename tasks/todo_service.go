package tasks

import (
	"context"

	"github.com/kjain0073/go-Todo/models"
)

type Service interface {
	CreateTodoToKafka(ctx context.Context, Title string) (string, error)
	CreateTodoToRepo(ctx context.Context, todo models.TodoEntity) (string, error)
	GetTodos(ctx context.Context) ([]models.TodoDto, error)
	DeleteTodo(ctx context.Context, Id string) (string, error)
	UpdateTodo(ctx context.Context, Id string, Title string, Completed bool) (string, error)
}
