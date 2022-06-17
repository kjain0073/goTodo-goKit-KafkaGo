package tasks

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateTodo endpoint.Endpoint
	GetTodos   endpoint.Endpoint
	DeleteTodo endpoint.Endpoint
	UpdateTodo endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateTodo: makeCreateTodoEndpoint(s),
		GetTodos:   makeGetTodosEndpoint(s),
		DeleteTodo: makeDeleteTodoEndpoint(s),
		UpdateTodo: makeUpdateTodoEndpoint(s),
	}
}

func makeCreateTodoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateTodoRequest)
		ok, err := s.CreateTodoToKafka(ctx, req.Title)
		return CreateTodoResponse{Ok: ok}, err
	}
}

func makeGetTodosEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		todoList, err := s.GetTodos(ctx)
		return GetTodosResponse{
			TodoList: todoList,
		}, err
	}
}

func makeDeleteTodoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteTodoRequest)
		ok, err := s.DeleteTodo(ctx, req.Id)
		return DeleteTodoResponse{Ok: ok}, err
	}
}

func makeUpdateTodoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateTodoRequest)
		ok, err := s.UpdateTodo(ctx, req.Id, req.Title, req.Completed)
		return UpdateTodoResponse{Ok: ok}, err
	}
}
