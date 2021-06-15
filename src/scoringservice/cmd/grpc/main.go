package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	scoring "go-hackathon/api/scoringservice"
	"go-hackathon/src/common/cmd"
	transportUtils "go-hackathon/src/common/infrastructure/transport"
	expressionsApi "go-hackathon/src/scoringservice/pkg/expressions/api"
	"go-hackathon/src/scoringservice/pkg/scoringtask/api"
	"go-hackathon/src/scoringservice/pkg/transport"
	"google.golang.org/grpc"
	"os"
)

const appID = "scoring"

type config struct {
	cmd.GRPCConfig
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

func startServer(killSignalCh <-chan os.Signal, c config) {
	l := transportUtils.ListenTCP(c.ServerPort)
	defer transportUtils.CloseService(l, "socket")

	db := cmd.CreateDBConnection(c.DatabaseConfig)
	defer transportUtils.CloseService(db, "db connection")

	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()

	scoring.RegisterScoringServiceServer(grpcServer, transport.Server(api.NewApi(db, expressionsApi.NewApi())))

	go func() {
		log.WithField("port", c.ServerPort).Info("starting the server")
		if err := grpcServer.Serve(l); err != nil {
			log.Fatalf("Failed to listen and serve: %v", err)
		}
	}()

	cmd.WaitForKillSignal(killSignalCh)
}
