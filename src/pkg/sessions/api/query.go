package api

import (
	"go-hackaton/src/pkg/sessions/api/output"
	"time"
)

func (a *api) GetSessions() ([]output.SessionOutput, error) {
	sessions, err := a.sqs.GetSessions()
	if err != nil {
		return nil, err
	}

	sessionsOutput := make([]output.SessionOutput, len(sessions))
	for i, session := range sessions {
		sessionsOutput[i] = output.NewSessionOutput(session)
	}

	return sessionsOutput, nil
}

func (a *api) GetSession(id string) (*output.SessionOutput, error) {
	session, err := a.sqs.GetSession(id)
	if err != nil {
		return nil, err
	}

	out := output.NewSessionOutput(*session)

	return &out, nil
}

func (a *api) GetSessionParticipants(sessionId string) ([]output.ParticipantOutput, error) {
	participants, err := a.pqs.GetParticipants(sessionId)
	if err != nil {
		return nil, err
	}

	participantsOutput := make([]output.ParticipantOutput, len(participants))
	for i, participant := range participants {
		participantsOutput[i] = output.NewParticipantOutput(participant)
	}

	return participantsOutput, nil
}

func (a *api) GetFirstScoredParticipantBefore(time time.Time) (*output.ParticipantOutput, error) {
	participant, err := a.pqs.GetFirstScoredParticipantBefore(time)
	if err != nil {
		return nil, err
	}

	var participantOutput *output.ParticipantOutput
	if participant != nil {
		out := output.NewParticipantOutput(*participant)
		participantOutput = &out
	}

	return participantOutput, nil
}
