package view

import (
	"context"
	"errors"

	"github.com/go-kit/log"
	"github.com/kjain0073/go-Todo/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	HostName       string = "mongodb://localhost:27017"
	DbName         string = "demo_todo"
	CollectionName string = "todo"
	Port           string = ":9000"
)

//Repository layer
type MongoDbRepository interface {
	CreateTodo(ctx context.Context, todo models.TodoEntity) error
	GetTodos(ctx context.Context) ([]models.TodoEntity, error)
	DeleteTodo(ctx context.Context, id string) error
	UpdateTodo(ctx context.Context, id string, title string, completed bool) error
}

type mongodbrepo struct {
	db     *mgo.Database
	logger log.Logger
}

func NewMongoDbRepo(db *mgo.Database, logger log.Logger) MongoDbRepository {
	return &mongodbrepo{
		db:     db,
		logger: log.With(logger, "repo", "mongoDB"),
	}
}

func (repo *mongodbrepo) GetTodos(ctx context.Context) ([]models.TodoEntity, error) {
	todos := []models.TodoEntity{}
	if err := repo.db.C(CollectionName).Find(bson.M{}).All(&todos); err != nil {
		return nil, errors.New("unable to Fetch Todos")
	}

	return todos, nil
}

func (repo *mongodbrepo) CreateTodo(ctx context.Context, todo models.TodoEntity) error {

	if todo.Title == "" {
		return errors.New("title is required")
	}
	if err := repo.db.C(CollectionName).Insert(&todo); err != nil {
		return err
	}
	return nil
}

func (repo *mongodbrepo) DeleteTodo(ctx context.Context, id string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("invalid Task Id")
	}

	if err := repo.db.C(CollectionName).RemoveId(bson.ObjectIdHex(id)); err != nil {
		return errors.New("failed to delete Todo")
	}

	return nil
}

func (repo *mongodbrepo) UpdateTodo(ctx context.Context, id string, title string, completed bool) error {

	if !bson.IsObjectIdHex(id) {
		return errors.New("id is invalid")
	}

	if title == "" {
		return errors.New("title field is required")
	}

	if err := repo.db.C(CollectionName).
		Update(
			bson.M{"_id": bson.ObjectIdHex(id)},
			bson.M{"title": title, "completed": completed},
		); err != nil {
		return errors.New("failed to update Todo")
	}
	return nil
}
