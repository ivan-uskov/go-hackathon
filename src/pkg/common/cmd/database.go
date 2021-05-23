package cmd

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func CreateDBConnection(c DatabaseConfig) *sql.DB {
	arguments := c.DatabaseArguments
	if len(arguments) > 0 {
		arguments = "?" + arguments
	}

	dsn := fmt.Sprintf("%s:%s@%s/%s%s", c.DatabaseUser, c.DatabasePassword, c.DatabaseAddress, c.DatabaseName, arguments)
	db, err := sql.Open(c.DatabaseDriver, dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Debugf("Connection to %s established", dsn)

	return db
}
