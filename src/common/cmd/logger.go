package cmd

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func SetupLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}
