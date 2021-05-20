package repository

import (
	"context"
	"database/sql"
	log "github.com/sirupsen/logrus"
)

type Database struct {
	*sql.DB
}

func (d *Database) Tx(job func(*sql.Tx) error) error {
	ctx := context.Background()
	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = job(tx)

	if err == nil {
		return tx.Commit()
	}

	log.Error(tx.Rollback())
	return err
}
