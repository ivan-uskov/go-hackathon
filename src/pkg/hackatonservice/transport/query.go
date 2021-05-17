package transport

import (
	"fmt"
	"go-hackaton/src/pkg/common/application/errors"
	"go-hackaton/src/pkg/common/transport"
	"net/http"
	"text/template"
	"time"
)

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
