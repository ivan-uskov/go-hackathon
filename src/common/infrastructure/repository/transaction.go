package repository

import (
	"database/sql"
	"errors"
	log "github.com/sirupsen/logrus"
	"go-hackathon/src/common/infrastructure"
)

const lockTimeoutSeconds = 20

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

func (d *Database) Lock(name string, job func() error) error {
	err := d.getLock(name)
	if err != nil {
		return err
	}

	defer func() {
		err := d.releaseLock(name)
		if err != nil {
			log.Error(err)
		}
	}()

	return job()
}

func (d *Database) getLock(name string) error {
	rows, err := d.Query("SELECT GET_LOCK(MD5(?), ?)", name, lockTimeoutSeconds)
	if err != nil {
		return err
	}
	defer infrastructure.Close(rows)

	var ok int
	if rows.Next() {
		err = rows.Scan(&ok)
		if err != nil {
			return err
		}
	}

	if ok != 1 {
		return errors.New("get lock timeout or error")
	}

	return nil
}

func (d *Database) releaseLock(name string) error {
	rows, err := d.Query("SELECT RELEASE_LOCK(MD5(?))", name)
	if err != nil {
		return err
	}
	defer infrastructure.Close(rows)

	var ok int
	if rows.Next() {
		err = rows.Scan(&ok)
		if err != nil {
			return err
		}
	}

	if ok != 1 {
		return errors.New("get lock timeout or error")
	}

	return nil
}
