package transport

import (
	"github.com/gorilla/mux"
	"go-hackaton/src/pkg/common/transport"
	sessions "go-hackaton/src/pkg/sessions/api"
	"net/http"
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
	s.HandleFunc("/session/{ID:[0-9a-zA-Z-]+}", srv.closeSession).Methods(http.MethodDelete)
	s.HandleFunc("/session/{ID:[0-9a-zA-Z-]+}/participant", srv.addSessionParticipant).Methods(http.MethodPost)
	s.HandleFunc("/session/{ID:[0-9a-zA-Z-]+}/participants", srv.getSessionParticipants).Methods(http.MethodGet)

	return transport.LogMiddleware(r)
}
