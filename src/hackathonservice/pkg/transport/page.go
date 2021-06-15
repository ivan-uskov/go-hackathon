package transport

import (
	"fmt"
	"go-hackathon/src/common/application/errors"
	"go-hackathon/src/common/infrastructure/transport"
	"net/http"
	"text/template"
)

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
