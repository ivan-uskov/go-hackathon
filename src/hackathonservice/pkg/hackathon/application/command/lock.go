package command

import (
	"fmt"
	"github.com/google/uuid"
)

func getHackathonNameLock(name string) string {
	return fmt.Sprintf("hackathon_name_%s", name)
}

func getHackathonIDLock(id uuid.UUID) string {
	return fmt.Sprintf("hackathon_id_%s", id.String())
}

func getParticipantNameLock(name string) string {
	return fmt.Sprintf("hackathon_participant_name_%s", name)
}
