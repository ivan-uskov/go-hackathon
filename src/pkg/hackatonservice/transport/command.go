package transport

import (
	"go-hackaton/src/pkg/common/application/errors"
	"go-hackaton/src/pkg/common/transport"
	"go-hackaton/src/pkg/sessions/api/input"
	"net/http"
)

type addSessionRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type addSessionResponse struct {
	ID string `json:"id"`
}

func (s *server) addSession(w http.ResponseWriter, r *http.Request) {
	var request addSessionRequest
	err := transport.ReadJson(r, &request)
	if err != nil {
		transport.ProcessError(w, err)
		return
	}

	id, err := s.api.AddSession(input.AddSessionInput{
		Code: request.Code,
		Name: request.Name,
		Type: request.Type,
	})
	if err != nil {
		transport.ProcessError(w, err)
		return
	}

	transport.RenderJson(w, addSessionResponse{ID: id.String()})
}

func (s *server) closeSession(w http.ResponseWriter, r *http.Request) {
	id, found := transport.Parameter(r, "ID")
	if !found {
		transport.ProcessError(w, errors.InvalidArgumentError)
	}

	err := s.api.CloseSession(input.CloseSessionInput{SessionID: id})
	if err != nil {
		transport.ProcessError(w, err)
		return
	}
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
