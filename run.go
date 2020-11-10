package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/utsavgupta/go-leader-election/jobs"

	"cloud.google.com/go/datastore"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/julienschmidt/httprouter"

	"github.com/utsavgupta/go-leader-election/globals"
	"github.com/utsavgupta/go-leader-election/handlers"
	"github.com/utsavgupta/go-leader-election/schedulers"
)

func main() {
	start := time.Now()

	ctx, fnCancel := context.WithCancel(context.Background())
	defer fnCancel()

	globals.InitConfig()
	globals.InitLogger(globals.APPLICATIONNAME, globals.Config.Stage)

	nodeName := namesgenerator.GetRandomName(1)

	client, e := datastore.NewClient(ctx, "")

	if e != nil {
		panic(e)
	}

	r := httprouter.New()

	r.GET("/.well-known/live", handlers.Wrap(handlers.Ok))
	r.GET("/.well-known/ready", handlers.Wrap(handlers.Ok))

	r.GET("/", handlers.Wrap(handlers.Default))

	intr := make(chan os.Signal)
	err := make(chan error)

	go func(e chan error) {
		err := http.ListenAndServe(fmt.Sprintf(":%d", globals.Config.Port), r)
		e <- err
	}(err)

	go func(c context.Context, cl *datastore.Client) {
		goPreacherScheduler := schedulers.NewScheduler(nodeName, cl)
		goPreacherScheduler(c, jobs.NewPreacher("go_preacher", "Go is a modern language for the web", os.Stdout), 1*time.Minute)
	}(ctx, client)

	go func(c context.Context, cl *datastore.Client) {
		scalaPreacherScheduler := schedulers.NewScheduler(nodeName, cl)
		scalaPreacherScheduler(c, jobs.NewPreacher("scala_preacher", "Scala is the most exciting language on the JVM", os.Stdout), 1*time.Minute)
	}(ctx, client)

	signal.Notify(intr, os.Interrupt)

	globals.Logger.Infof(ctx, "Started server in %dms. Listening to requests on port %d.", time.Since(start)/time.Millisecond, globals.Config.Port)

	select {
	case i := <-intr:
		globals.Logger.Infof(ctx, "Received interrupt %+v", i)
	case e := <-err:
		globals.Logger.Error(ctx, e)
	}

	globals.Logger.Infof(ctx, "Stopping server ...")
	fnCancel()
	time.Sleep(30 * time.Second)
	globals.Logger.Infof(ctx, "Bye!")
}
