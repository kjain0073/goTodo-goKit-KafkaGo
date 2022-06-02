package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type TodoEntity struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Title     string        `bson:"title"`
	Completed bool          `bson:"completed"`
	CreatedAt time.Time     `bson:"createAt"`
}
