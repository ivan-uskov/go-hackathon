package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"go-hackathon/src/common/cmd"
	"go-hackathon/src/common/cmd/events"
	"go-hackathon/src/common/cmd/supervisor"
	"go-hackathon/src/common/infrastructure/rabbitmq"
	"go-hackathon/src/common/infrastructure/transport"
	"os"
)

const appID = "scoring"

type config struct {
	cmd.DatabaseConfig
	cmd.AMQPConfig

	FailTimeoutSeconds int `envconfig:"fail_timeout_seconds"`
}

func main() {
	var c config
	if err := envconfig.Process(appID, &c); err != nil {
		log.Fatal(err)
	}

	cmd.SetupLogger()

	killSignalChan := cmd.GetKillSignalChan()
	startWorker(killSignalChan, c)
}

func startWorker(killSignalChan <-chan os.Signal, c config) {
	db := cmd.CreateDBConnection(c.DatabaseConfig)
	defer transport.CloseService(db, "database connection")

	amqp := rabbitmq.NewProducer(c.AMQPConfig)
	defer amqp.Close()

	s := supervisor.StartSupervisor(events.NewProducer(db, amqp).Produce, c.FailTimeoutSeconds)

	cmd.WaitForKillSignal(killSignalChan)

	s.Stop()
}
