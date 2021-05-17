package transport

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	sessions "go-hackaton/src/pkg/sessions/api"
	"net/http"
	"time"
)

type server struct {
	api sessions.Api
}

func Router(api sessions.Api) http.Handler {
	srv := &server{api: api}

	r := mux.NewRouter()
	r.HandleFunc("/session/{ID:[0-9a-zA-Z-]+}/participants", srv.getSessionParticipantsPage).Methods(http.MethodGet)

	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/sessions", srv.sessionsList).Methods(http.MethodGet)
	s.HandleFunc("/session", srv.addSession).Methods(http.MethodPost)
	s.HandleFunc("/session/{ID:[0-9a-zA-Z-]+}/participant", srv.addSessionParticipant).Methods(http.MethodPost)
	s.HandleFunc("/session/{ID:[0-9a-zA-Z-]+}/participants", srv.getSessionParticipants).Methods(http.MethodGet)

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
