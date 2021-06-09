package transport

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	hackathon "go-hackathon/api/hackathonservice"
	"go-hackathon/src/hackathonservice/pkg/hackathon/api/input"
)

func (s *server) AddHackathon(_ context.Context, request *hackathon.AddHackathonRequest) (*hackathon.AddHackathonResponse, error) {
	id, err := s.api.AddHackathon(input.AddHackathonInput{
		Name: request.Name,
		Type: request.Type,
	})

	if err != nil {
		return nil, err
	}

	return &hackathon.AddHackathonResponse{Id: id.String()}, nil
}

func (s *server) CloseHackathon(_ context.Context, request *hackathon.CloseHackathonRequest) (*empty.Empty, error) {
	err := s.api.CloseHackathon(input.CloseHackathonInput{HackathonID: request.Id})
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *server) AddHackathonParticipant(_ context.Context, request *hackathon.AddHackathonParticipantRequest) (*empty.Empty, error) {
	err := s.api.AddHackathonParticipant(input.AddHackathonParticipantInput{
		HackathonID: request.Id,
		Name:        request.Name,
		Endpoint:    request.Endpoint,
	})

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
