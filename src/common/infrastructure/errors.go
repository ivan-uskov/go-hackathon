package infrastructure

import (
	log "github.com/sirupsen/logrus"
	"go-hackathon/src/common/application/errors"
)

func InternalError(e error) error {
	if e != nil {
		log.Error(e)
	}

	return errors.InternalError
}
