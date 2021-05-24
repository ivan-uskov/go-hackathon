package transport

import (
	"github.com/gorilla/mux"
	"go-hackathon/src/common/transport"
	"go-hackathon/src/hackathonservice/pkg/hackathon/api"
	"net/http"
)

type server struct {
	api api.Api
}

func Router(api api.Api) http.Handler {
	srv := &server{api: api}

	r := mux.NewRouter()
	r.HandleFunc("/hackathon/{ID:[0-9a-zA-Z-]+}/participants", srv.getHackathonParticipantsPage).Methods(http.MethodGet)

	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/hackathons", srv.hackathonsList).Methods(http.MethodGet)
	s.HandleFunc("/hackathon", srv.addHackathon).Methods(http.MethodPost)
	s.HandleFunc("/hackathon/{ID:[0-9a-zA-Z-]+}", srv.closeHackathon).Methods(http.MethodDelete)
	s.HandleFunc("/hackathon/{ID:[0-9a-zA-Z-]+}/participant", srv.addHackathonParticipant).Methods(http.MethodPost)
	s.HandleFunc("/hackathon/{ID:[0-9a-zA-Z-]+}/participants", srv.getHackathonParticipants).Methods(http.MethodGet)

	return transport.LogMiddleware(r)
}
