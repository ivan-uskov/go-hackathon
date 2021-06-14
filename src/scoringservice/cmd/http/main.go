package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"go-hackathon/src/common/cmd"
	transportUtils "go-hackathon/src/common/cmd/transport"
	expressionsApi "go-hackathon/src/scoringservice/pkg/expressions/api"
	"go-hackathon/src/scoringservice/pkg/scoringtask/api"
	"go-hackathon/src/scoringservice/pkg/transport"
	"net/http"
	"os"
)

const appID = "scoring"

type config struct {
	cmd.HTTPConfig
	cmd.DatabaseConfig
}

func main() {
	var c config
	if err := envconfig.Process(appID, &c); err != nil {
		log.Fatal(err)
	}

	cmd.SetupLogger()

	killSignalChan := cmd.GetKillSignalChan()
	startServer(killSignalChan, c)
}

func startServer(killSignalChan <-chan os.Signal, c config) {
	db := cmd.CreateDBConnection(c.DatabaseConfig)
	defer transportUtils.CloseService(db, "database connection")

	router := transport.Router(context.Background(), api.NewApi(db, expressionsApi.NewApi()))
	srv := &http.Server{Addr: fmt.Sprintf(":%d", c.ServerPort), Handler: router}

	go func() {
		log.WithFields(log.Fields{"port": c.ServerPort}).Info("starting the server")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Failed to listen and serve: %v", err)
		}
	}()

	cmd.WaitForKillSignal(killSignalChan)

	log.Info("Shutting down the http gateway server")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Errorf("Failed to shutdown http gateway server: %v", err)
	}
}
