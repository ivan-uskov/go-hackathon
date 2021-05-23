package main

import (
	"context"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"go-hackaton/src/pkg/common/cmd"
	"go-hackaton/src/pkg/scoringservice/transport"
	"net/http"
)

const appID = "scoring"

type config struct {
	cmd.WebConfig
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

	router := transport.Router()
	srv := &http.Server{Addr: fmt.Sprintf(":%s", c.ServerPort), Handler: router}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	return srv
}
