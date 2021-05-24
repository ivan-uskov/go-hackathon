package repository

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"go-hackathon/src/common/infrastructure"
)

type Database struct {
	*sql.DB
}

func (d *Database) Tx(job func(*sql.Tx) error) error {
	tx, err := d.Begin()
	if err != nil {
		return infrastructure.InternalError(err)
	}

	err = job(tx)

	if err == nil {
		err2 := tx.Commit()
		if err2 != nil {
			log.Error(err2)
		}
	} else {
		err2 := tx.Rollback()
		if err2 != nil {
			log.Error(err2)
		}
	}

	return err
}
