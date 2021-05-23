package transport

import (
	"github.com/gorilla/mux"
	"go-hackaton/src/pkg/common/transport"
	"net/http"
)

type server struct{}

func (s *server) addTask(_ http.ResponseWriter, _ *http.Request) {
}

func (s *server) removeTasks(_ http.ResponseWriter, _ *http.Request) {
}

func Router() http.Handler {
	srv := &server{}

	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()

	s.HandleFunc("/task", srv.addTask).Methods(http.MethodPost)
	s.HandleFunc("/tasks", srv.removeTasks).Methods(http.MethodDelete)

	return transport.LogMiddleware(r)
}
