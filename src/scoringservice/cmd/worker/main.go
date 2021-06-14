package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"go-hackathon/src/common/cmd"
	"go-hackathon/src/common/cmd/supervisor"
	transportUtils "go-hackathon/src/common/cmd/transport"
	expressionsApi "go-hackathon/src/scoringservice/pkg/expressions/api"
	"go-hackathon/src/scoringservice/pkg/scoringtask/api"
	"os"
)

const appID = "scoring"

type config struct {
	cmd.DatabaseConfig

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
	defer transportUtils.CloseService(db, "database connection")

	s := supervisor.StartSupervisor(api.NewApi(db, expressionsApi.NewApi()).ScoreOnce, c.FailTimeoutSeconds)

	cmd.WaitForKillSignal(killSignalChan)

	s.Stop()
}
