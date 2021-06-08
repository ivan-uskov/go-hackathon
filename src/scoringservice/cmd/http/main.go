package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"go-hackathon/src/common/cmd"
	transportUtils "go-hackathon/src/common/cmd/transport"
	"go-hackathon/src/scoringservice/pkg/scoringtask/api"
	"go-hackathon/src/scoringservice/pkg/transport"
	"net/http"
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

	ctx, stopServer := context.WithCancel(context.Background())
	startServer(ctx, &c)

	cmd.WaitForKillSignal(killSignalChan)
	stopServer()
}

func startServer(ctx context.Context, c *config) {
	db := cmd.CreateDBConnection(c.DatabaseConfig)
	router := transport.Router(ctx, api.NewApi(db))

	srv := &http.Server{Addr: fmt.Sprintf(":%s", c.ServerPort), Handler: router}

	go func() {
		log.WithFields(log.Fields{"port": c.ServerPort}).Info("starting the server")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Failed to listen and serve: %v", err)
		}
	}()

	go func() {
		<-ctx.Done()
		log.Info("Shutting down the http gateway server")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Errorf("Failed to shutdown http gateway server: %v", err)
		}

		transportUtils.CloseService(db, "database connection")
	}()
}
