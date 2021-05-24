package api

import (
	"go-hackathon/src/hackathonservice/pkg/hackathon/api/errors"
	"go-hackathon/src/hackathonservice/pkg/hackathon/api/output"
)

func (a *api) GetHackathons() ([]output.HackathonOutput, error) {
	hackathons, err := a.sqs.GetHackathons()
	if err != nil {
		return nil, errors.WrapError(err)
	}

	hackathonOutputs := make([]output.HackathonOutput, len(hackathons))
	for i, hackathon := range hackathons {
		hackathonOutputs[i] = output.NewHackathonOutput(hackathon)
	}

	return hackathonOutputs, nil
}

func (a *api) GetHackathon(id string) (*output.HackathonOutput, error) {
	hackathon, err := a.sqs.GetHackathon(id)
	if err != nil {
		return nil, errors.WrapError(err)
	}

	var hackathonOutput *output.HackathonOutput
	if hackathon != nil {
		out := output.NewHackathonOutput(*hackathon)
		hackathonOutput = &out
	}

	return hackathonOutput, nil
}

func (a *api) GetHackathonParticipants(hackathonID string) ([]output.ParticipantOutput, error) {
	participants, err := a.pqs.GetParticipants(hackathonID)
	if err != nil {
		return nil, errors.WrapError(err)
	}

	participantsOutput := make([]output.ParticipantOutput, len(participants))
	for i, participant := range participants {
		participantsOutput[i] = output.NewParticipantOutput(participant)
	}

	return participantsOutput, nil
}
