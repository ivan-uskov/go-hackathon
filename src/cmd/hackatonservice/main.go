package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"go-hackaton/src/pkg/common/cmd"
	"go-hackaton/src/pkg/hackatonservice/transport"
	sessions "go-hackaton/src/pkg/sessions/api"
	tasks "go-hackaton/src/pkg/tasks/api"
	"net/http"
)

const appID = "hackaton"

type config struct {
	cmd.WebConfig
	cmd.DatabaseConfig
}

func main() {
	var c config
	if err := envconfig.Process(appID, &c); err != nil {
		log.Fatal(err)
	}

	cmd.SetupLogger()

	killSignalChan := cmd.GetKillSignalChan()
	srv := startServer(&c)

	cmd.WaitForKillSignal(killSignalChan)
	log.Fatal(srv.Shutdown(context.Background()))
}

func startServer(c *config) *http.Server {
	log.WithFields(log.Fields{"port": c.ServerPort}).Info("starting the server")
	db := cmd.CreateDBConnection(c.DatabaseConfig)
	sessionsApi := sessions.NewApi(db, tasks.NewApi())

	router := transport.Router(sessionsApi)
	srv := &http.Server{Addr: fmt.Sprintf(":%s", c.ServerPort), Handler: router}
	go func() {
		log.Fatal(srv.ListenAndServe())
		log.Fatal(db.Close())
	}()

	return srv
}
