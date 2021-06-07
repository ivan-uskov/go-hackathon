package transport

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	hackathon "go-hackathon/api"
	"go-hackathon/src/common/application/errors"
	"go-hackathon/src/common/cmd/transport"
	"net/http"
	"text/template"
)

func (s *server) GetHackathons(_ context.Context, _ *empty.Empty) (*hackathon.HackathonsResponse, error) {
	hh, err := s.api.GetHackathons()
	if err != nil {
		return nil, err
	}

	hr := make([]*hackathon.HackathonsResponse_Hackathon, len(hh))
	for i, h := range hh {
		hr[i] = &hackathon.HackathonsResponse_Hackathon{
			ID:           h.ID,
			Name:         h.Name,
			Participants: int64(h.Participants),
			Type:         h.Type,
			CreatedAt:    h.CreatedAt.String(),
			ClosedAt:     transport.TimeToString(h.ClosedAt),
		}
	}

	return &hackathon.HackathonsResponse{Items: hr}, nil
}

func (s *server) GetHackathonParticipants(_ context.Context, request *hackathon.HackathonParticipantsRequest) (*hackathon.HackathonParticipantsResponse, error) {
	pp, err := s.api.GetHackathonParticipants(request.ID)
	if err != nil {
		return nil, err
	}

	pr := make([]*hackathon.HackathonParticipantsResponse_Participant, len(pp))
	for i, po := range pp {
		pr[i] = &hackathon.HackathonParticipantsResponse_Participant{
			ID:        po.ID,
			Name:      po.Name,
			Score:     int64(po.Score),
			CreatedAt: po.CreatedAt.String(),
			ScoredAt:  transport.TimeToString(po.ScoredAt),
		}
	}

	return &hackathon.HackathonParticipantsResponse{Items: pr}, nil
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

	h, err := s.api.GetHackathon(id)
	if err != nil {
		transport.ProcessError(w, err)
		return
	}

	args := hackathonParticipantsTemplateArgs{
		fmt.Sprintf("/api/v1/hackathon/%s/participants", id),
		fmt.Sprintf("/api/v1/hackathon/%s/participant", id),
		h.Name,
	}

	err = hackathonParticipantsTemplate.Execute(w, args)
	if err != nil {
		transport.ProcessError(w, err)
	}
}
