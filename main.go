package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log/level"
	mgo "gopkg.in/mgo.v2"

	"github.com/kjain0073/go-Todo/adapters"
	"github.com/kjain0073/go-Todo/tasks"
	"github.com/kjain0073/go-Todo/view"
)

func main() {
	var httpAddr = flag.String("http", ":8080", "http listen address")

	//init logger
	logger := adapters.InitLogger()
	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")
	ctx := context.Background()

	var db *mgo.Database
	adapters.SetConnection(logger, ctx)
	db = adapters.GetConnection()

	flag.Parse()

	//init a service
	var srv tasks.Service
	{
		mongodbrepo := view.NewMongoDbRepo(db, logger)
		srv = view.NewService(mongodbrepo, logger)
		go func() {
			view.PollFromKafka(srv, ctx, logger)
		}()
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	endpoints := tasks.MakeEndpoints(srv)

	go func() {
		fmt.Println("listening on port", *httpAddr)
		handler := tasks.NewHTTPServer(ctx, endpoints)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	level.Error(logger).Log("exit", <-errs)
}
