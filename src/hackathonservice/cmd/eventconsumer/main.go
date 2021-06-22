package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	scoring "go-hackathon/api/scoringservice"
	"go-hackathon/src/common/cmd"
	"go-hackathon/src/common/infrastructure/rabbitmq"
	transportUtils "go-hackathon/src/common/infrastructure/transport"
	"go-hackathon/src/hackathonservice/pkg/hackathon/api"
	"os"
)

const appID = "hackathon"

type config struct {
	cmd.DatabaseConfig
	cmd.AMQPConfig

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

	amqp := rabbitmq.NewConsumer(c.AMQPConfig)
	defer amqp.Close()

	go amqp.Consume(api.NewApi(db, scoring.NewScoringServiceClient(scoringConn)).ProcessEvent)

	cmd.WaitForKillSignal(killSignalChan)
}
