package transport

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go-hackathon/src/common/application/errors"
	"io"
	"net/http"
	"strings"
	"time"
)

func ProcessError(w http.ResponseWriter, e error) {
	if e == errors.InternalError {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	} else {
		http.Error(w, e.Error(), http.StatusBadRequest)
	}
}

func Parameter(r *http.Request, key string) (string, bool) {
	val, found := mux.Vars(r)[key]
	return val, found
}

func CloseService(closer io.Closer, subject ...string) {
	log.Infof("Close %v", subject)
	Close(closer, subject...)
}

func Close(closer io.Closer, subject ...string) {
	subjectStr := strings.Join(subject, "")

	err := closer.Close()
	if err != nil {
		log.Errorf("Failed to close %v: %v", subjectStr, err)
	}
}

func TimeToString(t *time.Time) string {
	if t == nil {
		return ""
	}

	return t.String()
}
