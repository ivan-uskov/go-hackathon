package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	scoring "go-hackathon/api/scoringservice"
	"go-hackathon/src/common/cmd"
	transportUtils "go-hackathon/src/common/cmd/transport"
	"go-hackathon/src/scoringservice/pkg/scoringtask/api"
	"go-hackathon/src/scoringservice/pkg/transport"
	"google.golang.org/grpc"
	"net"
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

	ctx, stopServer := context.WithCancel(context.Background())
	startServer(ctx, &c)

	cmd.WaitForKillSignal(killSignalChan)
	stopServer()
}

func startServer(ctx context.Context, c *config) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", c.ServerPort))
	if err != nil {
		log.Fatal(err)
	}

	db := cmd.CreateDBConnection(c.DatabaseConfig)
	server := transport.Server(api.NewApi(db))

	grpcServer := grpc.NewServer()
	scoring.RegisterScoringServiceServer(grpcServer, server)

	go func() {
		log.WithField("port", c.ServerPort).Info("starting the server")
		if err := grpcServer.Serve(l); err != nil {
			log.Fatalf("Failed to listen and serve: %v", err)
		}
	}()

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
		transportUtils.CloseService(db, "db connection")
		transportUtils.CloseService(l, "socket")
	}()
}
