package transport

import (
	"go-hackathon/src/common/application/errors"
	"go-hackathon/src/common/cmd/transport"
	"go-hackathon/src/hackathonservice/pkg/hackathon/api/input"
	"net/http"
)

type addHackathonRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type addHackathonResponse struct {
	ID string `json:"id"`
}

func (s *server) addHackathon(w http.ResponseWriter, r *http.Request) {
	var request addHackathonRequest
	err := transport.ReadJson(r, &request)
	if err != nil {
		transport.ProcessError(w, err)
		return
	}

	id, err := s.api.AddHackathon(input.AddHackathonInput{
		Code: request.Code,
		Name: request.Name,
		Type: request.Type,
	})
	if err != nil {
		transport.ProcessError(w, err)
		return
	}

	transport.RenderJson(w, addHackathonResponse{ID: id.String()})
}

func (s *server) closeHackathon(w http.ResponseWriter, r *http.Request) {
	id, found := transport.Parameter(r, "ID")
	if !found {
		transport.ProcessError(w, errors.InvalidArgumentError)
	}

	err := s.api.CloseHackathon(input.CloseHackathonInput{HackathonID: id})
	if err != nil {
		transport.ProcessError(w, err)
		return
	}
}

type addHackathonParticipantRequest struct {
	HackathonCode string `json:"hackathon_code"`
	Endpoint      string `json:"endpoint"`
	Name          string `json:"name"`
}

func (s *server) addHackathonParticipant(w http.ResponseWriter, r *http.Request) {
	id, found := transport.Parameter(r, "ID")
	if !found {
		transport.ProcessError(w, errors.InvalidArgumentError)
	}

	var request addHackathonParticipantRequest
	err := transport.ReadJson(r, &request)
	if err != nil {
		transport.ProcessError(w, err)
		return
	}

	err = s.api.AddHackathonParticipant(input.AddHackathonParticipantInput{
		HackathonID:   id,
		HackathonCode: request.HackathonCode,
		Name:          request.Name,
		Endpoint:      request.Endpoint,
	})
	if err != nil {
		transport.ProcessError(w, err)
		return
	}
}
