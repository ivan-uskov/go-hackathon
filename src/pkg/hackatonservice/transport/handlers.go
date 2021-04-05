package transport

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go-hackaton/src/pkg/common/transport"
	sessions "go-hackaton/src/pkg/sessions/api"
	"net/http"
	"time"
)

type server struct {
	api sessions.Api
}

func (s *server) sessionsList(w http.ResponseWriter, _ *http.Request) {
	ss, err := s.api.GetSessions()
	if err != nil {
		transport.ProcessError(w, err)
		return
	}

	transport.RenderJson(w, ss)
}

func Router(api sessions.Api) http.Handler {
	srv := &server{api: api}

	r := mux.NewRouter()
	r.HandleFunc("/sessions", srv.sessionsList).Methods(http.MethodGet)

	return logMiddleware(r)
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL,
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
			"duration":   time.Since(startTime).String(),
			"at":         startTime,
		}).Info("got request")
	})
}
