package transport

import (
	"database/sql"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type server struct {
}

func (s *server) sessionsList(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("hello"))
	if err != nil {
		log.Error(err)
	}
}

func Router(db *sql.DB) http.Handler {
	srv := makeServer(db)

	r := mux.NewRouter()
	r.HandleFunc("/sessions", srv.sessionsList).Methods(http.MethodGet)

	return logMiddleware(r)
}

func makeServer(db *sql.DB) *server {
	return &server{}
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
