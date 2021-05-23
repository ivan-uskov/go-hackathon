package infrastructure

import (
	log "github.com/sirupsen/logrus"
	"go-hackaton/src/pkg/common/application/errors"
)

func InternalError(e error) error {
	log.Error(e)
	return errors.InternalError
}
