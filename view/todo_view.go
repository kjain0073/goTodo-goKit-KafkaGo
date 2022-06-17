package view

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gofrs/uuid"

	"gopkg.in/mgo.v2/bson"

	"github.com/kjain0073/go-Todo/models"
	"github.com/kjain0073/go-Todo/tasks"
)

type service struct {
	mongoDbRepo MongoDbRepository
	logger      log.Logger
}

func NewService(mongodbrepo MongoDbRepository, logger log.Logger) tasks.Service {
	return &service{
		mongoDbRepo: mongodbrepo,
		logger:      logger,
	}
}

func (s service) CreateTodoToKafka(ctx context.Context, title string) (string, error) {
	logger := log.With(s.logger, "method", "CreateTodoToKafka")

	uuid, _ := uuid.NewV4()
	id := uuid.String()
	todo := models.TodoEntity{
		ID:        bson.NewObjectId(),
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	if err := SaveToKafka(ctx, todo); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("create todo", id)

	return "Success", nil
}

func (s service) CreateTodoToRepo(ctx context.Context, todo models.TodoEntity) (string, error) {
	logger := log.With(s.logger, "method", "CreateTodoToRepo")

	if err := s.mongoDbRepo.CreateTodo(ctx, todo); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	return "Success", nil
}

func (s service) GetTodos(ctx context.Context) ([]models.TodoDto, error) {
	logger := log.With(s.logger, "method", "GetTodo")

	todos, err := s.mongoDbRepo.GetTodos(ctx)

	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	todoList := []models.TodoDto{}
	todoIds := ""
	for _, t := range todos {
		todoIds += t.ID.Hex()
		todoIds += ", "
		todoList = append(todoList, models.TodoDto{
			ID:        t.ID.Hex(),
			Title:     t.Title,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt,
		})
	}

	logger.Log("Get Todos", todoIds)

	return todoList, nil
}

func (s service) DeleteTodo(ctx context.Context, id string) (string, error) {
	logger := log.With(s.logger, "method", "DeleteTodo")

	if err := s.mongoDbRepo.DeleteTodo(ctx, id); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("Delete Todo", id)

	return "Success", nil
}

func (s service) UpdateTodo(ctx context.Context, id string, title string, completed bool) (string, error) {
	logger := log.With(s.logger, "method", "UpdateTodo")

	if err := s.mongoDbRepo.UpdateTodo(ctx, id, title, completed); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("Update todo", id)

	return "Success", nil
}
