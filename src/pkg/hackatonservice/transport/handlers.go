package transport

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go-hackaton/src/pkg/common/application/errors"
	"go-hackaton/src/pkg/common/transport"
	sessions "go-hackaton/src/pkg/sessions/api"
	"go-hackaton/src/pkg/sessions/api/input"
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

type addSessionParticipantRequest struct {
	SessionCode string `json:"session_code"`
	Endpoint    string `json:"endpoint"`
	Name        string `json:"name"`
}

func (s *server) addSessionParticipant(w http.ResponseWriter, r *http.Request) {
	var request addSessionParticipantRequest
	err := transport.ReadJson(r, &request)
	if err != nil {
		transport.ProcessError(w, err)
		return
	}

	id, found := transport.Parameter(r, "ID")
	if !found {
		transport.ProcessError(w, errors.InvalidArgumentError)
	}

	err = s.api.AddSessionParticipant(input.AddSessionParticipantInput{
		SessionID:   id,
		SessionCode: request.SessionCode,
		Name:        request.Name,
		Endpoint:    request.Endpoint,
	})
	if err != nil {
		transport.ProcessError(w, err)
		return
	}
}

func Router(api sessions.Api) http.Handler {
	srv := &server{api: api}

	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/sessions", srv.sessionsList).Methods(http.MethodGet)
	s.HandleFunc("/session/{ID:[0-9a-zA-Z-]+}/participant", srv.addSessionParticipant).Methods(http.MethodPost)

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
