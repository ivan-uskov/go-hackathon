package command

import (
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/errors"
	"testing"
)

const mockParticipantScore = 5
const mockInvalidParticipantScore = -1

func TestUpdateParticipantScore(t *testing.T) {
	uow := &mockUnitOfWork{}
	_ = uow.mockParticipantRepository.Add(mockParticipant)
	h := updateParticipantScoreCommandHandler{uow}
	err := h.Handle(UpdateParticipantScoreCommand{mockParticipant.ID, mockParticipantScore})
	if err != nil {
		t.Error("Update participant score not works")
	}

	p, _ := uow.mockParticipantRepository.Get(mockParticipant.ID)
	if p == nil {
		t.Error("Participant not exists after update score")
	} else if p.Score != mockParticipantScore {
		t.Error("Score not updated")
	}
}

func TestUpdateParticipantInvalidScore(t *testing.T) {
	uow := &mockUnitOfWork{}
	_ = uow.mockParticipantRepository.Add(mockParticipant)
	h := updateParticipantScoreCommandHandler{uow}

	err := h.Handle(UpdateParticipantScoreCommand{mockParticipant.ID, mockInvalidParticipantScore})
	if err != errors.InvalidParticipantScoreError {
		t.Error("Update participant invalid score works")
	}
}
