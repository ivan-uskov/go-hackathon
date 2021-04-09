package transport

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go-hackaton/src/pkg/common/application/errors"
	"go-hackaton/src/pkg/common/transport"
	sessions "go-hackaton/src/pkg/sessions/api"
	"go-hackaton/src/pkg/sessions/api/input"
	"html/template"
	"net/http"
	"time"
)

type server struct {
	api sessions.Api
}

type sessionResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Participants int       `json:"participants"`
	Type         string    `json:"type"`
	CreatedAt    time.Time `json:"created_at"`
}

func (s *server) sessionsList(w http.ResponseWriter, _ *http.Request) {
	ss, err := s.api.GetSessions()
	if err != nil {
		transport.ProcessError(w, err)
		return
	}

	sr := make([]sessionResponse, len(ss))
	for i, so := range ss {
		sr[i] = sessionResponse{
			ID:           so.ID,
			Name:         so.Name,
			Participants: so.Participants,
			Type:         so.Type,
			CreatedAt:    so.CreatedAt,
		}
	}

	transport.RenderJson(w, sr)
}

type addSessionParticipantRequest struct {
	SessionCode string `json:"session_code"`
	Endpoint    string `json:"endpoint"`
	Name        string `json:"name"`
}

func (s *server) addSessionParticipant(w http.ResponseWriter, r *http.Request) {
	id, found := transport.Parameter(r, "ID")
	if !found {
		transport.ProcessError(w, errors.InvalidArgumentError)
	}

	var request addSessionParticipantRequest
	err := transport.ReadJson(r, &request)
	if err != nil {
		transport.ProcessError(w, err)
		return
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

type participantResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Score     int        `json:"score"`
	CreatedAt time.Time  `json:"created_at"`
	ScoredAt  *time.Time `json:"scored_at"`
}

func (s *server) getSessionParticipants(w http.ResponseWriter, r *http.Request) {
	id, found := transport.Parameter(r, "ID")
	if !found {
		transport.ProcessError(w, errors.InvalidArgumentError)
	}

	pp, err := s.api.GetSessionParticipants(id)
	if err != nil {
		transport.ProcessError(w, err)
		return
	}

	pr := make([]participantResponse, len(pp))
	for i, po := range pp {
		pr[i] = participantResponse{
			ID:        po.ID,
			Name:      po.Name,
			Score:     po.Score,
			CreatedAt: po.CreatedAt,
			ScoredAt:  po.ScoredAt,
		}
	}

	transport.RenderJson(w, pr)
}

var sessionParticipantsTemplate = template.Must(template.ParseFiles("/app/templates/session_participants.html"))

type sessionParticipantsTemplateArgs struct {
	LoadUrl     string
	AddUrl      string
	SessionName string
}

func (s *server) getSessionParticipantsPage(w http.ResponseWriter, r *http.Request) {
	id, found := transport.Parameter(r, "ID")
	if !found {
		transport.ProcessError(w, errors.InvalidArgumentError)
	}

	session, err := s.api.GetSession(id)
	if err != nil {
		transport.ProcessError(w, err)
		return
	}

	args := sessionParticipantsTemplateArgs{
		fmt.Sprintf("/api/v1/session/%s/participants", id),
		fmt.Sprintf("/api/v1/session/%s/participant", id),
		session.Name,
	}

	err = sessionParticipantsTemplate.Execute(w, args)
	if err != nil {
		transport.ProcessError(w, err)
	}
}

func Router(api sessions.Api) http.Handler {
	srv := &server{api: api}

	r := mux.NewRouter()
	r.HandleFunc("/session/{ID:[0-9a-zA-Z-]+}/participants", srv.getSessionParticipantsPage).Methods(http.MethodGet)

	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/sessions", srv.sessionsList).Methods(http.MethodGet)
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
