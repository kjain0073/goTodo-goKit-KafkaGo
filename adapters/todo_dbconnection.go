package adapters

import (
	"context"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	mgo "gopkg.in/mgo.v2"

	"github.com/kjain0073/go-Todo/view"
)

var db *mgo.Database

func GetConnection() *mgo.Database {
	return db
}

func SetConnection(logger log.Logger, ctx context.Context) {
	// to connect with DB
	sess, err := mgo.Dial(view.HostName)
	if err != nil {
		level.Error(logger).Log("exit", err)
		os.Exit(-1)
	}
	sess.SetMode(mgo.Monotonic, true)
	db = sess.DB(view.DbName)

}
