package transport

import (
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	"go-hackathon/src/common/application/errors"
	"io"
	"net/http"
	"strings"
	"time"
)

var grpcServeMuxOptions = &runtime.JSONPb{
	EmitDefaults: true,
	OrigName:     true,
}

func NewServeMux() *runtime.ServeMux {
	return runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, grpcServeMuxOptions))
}

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
