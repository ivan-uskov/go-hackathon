package transport

import (
	"fmt"
	"go-hackathon/src/common/application/errors"
	"go-hackathon/src/common/transport"
	"net/http"
	"text/template"
	"time"
)

type hackathonResponse struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Participants int        `json:"participants"`
	Type         string     `json:"type"`
	CreatedAt    time.Time  `json:"created_at"`
	ClosedAt     *time.Time `json:"closed_at"`
}

func (s *server) hackathonsList(w http.ResponseWriter, _ *http.Request) {
	hh, err := s.api.GetHackathons()
	if err != nil {
		transport.ProcessError(w, err)
		return
	}

	hr := make([]hackathonResponse, len(hh))
	for i, h := range hh {
		hr[i] = hackathonResponse{
			ID:           h.ID,
			Name:         h.Name,
			Participants: h.Participants,
			Type:         h.Type,
			CreatedAt:    h.CreatedAt,
			ClosedAt:     h.ClosedAt,
		}
	}

	transport.RenderJson(w, hr)
}

type participantResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Score     int        `json:"score"`
	CreatedAt time.Time  `json:"created_at"`
	ScoredAt  *time.Time `json:"scored_at"`
}

func (s *server) getHackathonParticipants(w http.ResponseWriter, r *http.Request) {
	id, found := transport.Parameter(r, "ID")
	if !found {
		transport.ProcessError(w, errors.InvalidArgumentError)
		return
	}

	pp, err := s.api.GetHackathonParticipants(id)
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

var hackathonParticipantsTemplate = template.Must(template.ParseFiles("/app/templates/hackathon_participants.html"))

type hackathonParticipantsTemplateArgs struct {
	LoadUrl       string
	AddUrl        string
	HackathonName string
}

func (s *server) getHackathonParticipantsPage(w http.ResponseWriter, r *http.Request) {
	id, found := transport.Parameter(r, "ID")
	if !found {
		transport.ProcessError(w, errors.InvalidArgumentError)
		return
	}

	hackathon, err := s.api.GetHackathon(id)
	if err != nil {
		transport.ProcessError(w, err)
		return
	}

	args := hackathonParticipantsTemplateArgs{
		fmt.Sprintf("/api/v1/hackathon/%s/participants", id),
		fmt.Sprintf("/api/v1/hackathon/%s/participant", id),
		hackathon.Name,
	}

	err = hackathonParticipantsTemplate.Execute(w, args)
	if err != nil {
		transport.ProcessError(w, err)
	}
}
