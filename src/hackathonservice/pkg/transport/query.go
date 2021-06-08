package transport

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	hackathon "go-hackathon/api/hackathonservice"
	"go-hackathon/src/common/cmd/transport"
)

func (s *server) GetHackathons(_ context.Context, _ *empty.Empty) (*hackathon.HackathonsResponse, error) {
	hh, err := s.api.GetHackathons()
	if err != nil {
		return nil, err
	}

	hr := make([]*hackathon.HackathonsResponse_Hackathon, len(hh))
	for i, h := range hh {
		hr[i] = &hackathon.HackathonsResponse_Hackathon{
			Id:           h.ID,
			Name:         h.Name,
			Participants: int32(h.Participants),
			Type:         h.Type,
			CreatedAt:    h.CreatedAt.String(),
			ClosedAt:     transport.TimeToString(h.ClosedAt),
		}
	}

	return &hackathon.HackathonsResponse{Items: hr}, nil
}

func (s *server) GetHackathonParticipants(_ context.Context, request *hackathon.HackathonParticipantsRequest) (*hackathon.HackathonParticipantsResponse, error) {
	pp, err := s.api.GetHackathonParticipants(request.Id)
	if err != nil {
		return nil, err
	}

	pr := make([]*hackathon.HackathonParticipantsResponse_Participant, len(pp))
	for i, po := range pp {
		pr[i] = &hackathon.HackathonParticipantsResponse_Participant{
			Id:        po.ID,
			Name:      po.Name,
			Score:     int32(po.Score),
			CreatedAt: po.CreatedAt.String(),
			ScoredAt:  transport.TimeToString(po.ScoredAt),
		}
	}

	return &hackathon.HackathonParticipantsResponse{Items: pr}, nil
}
