package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	scoring "go-hackathon/api/scoringservice"
	"go-hackathon/src/common/cmd"
	transportUtils "go-hackathon/src/common/infrastructure/transport"
	"go-hackathon/src/hackathonservice/pkg/hackathon/api"
	"go-hackathon/src/hackathonservice/pkg/transport"
	"net/http"
	"os"
)

const appID = "hackathon"

type config struct {
	cmd.HTTPConfig
	cmd.DatabaseConfig

	ScoringGRPCAddress string `envconfig:"scoring_grpc_address"`
}

func main() {
	var c config
	if err := envconfig.Process(appID, &c); err != nil {
		log.Fatal(err)
	}

	cmd.SetupLogger()

	killSignalChan := cmd.GetKillSignalChan()
	startServer(killSignalChan, &c)
}

func startServer(killSignalChan <-chan os.Signal, c *config) {
	db := cmd.CreateDBConnection(c.DatabaseConfig)
	defer transportUtils.CloseService(db, "database connection")

	scoringConn := transportUtils.DialGRPC(c.ScoringGRPCAddress)
	defer transportUtils.CloseService(scoringConn, "scoring connection")

	handler := transport.Router(context.Background(), api.NewApi(db, scoring.NewScoringServiceClient(scoringConn)))
	srv := &http.Server{Addr: fmt.Sprintf(":%d", c.ServerPort), Handler: handler}

	go func() {
		log.WithField("port", c.ServerPort).Info("starting the server")
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
