package output

import (
	"go-hackaton/src/pkg/sessions/application/query/data"
	"go-hackaton/src/pkg/sessions/model"
	"time"
)

const sessionTypeArithmeticExpression = "arithmetic expression"
const unexpectedSessionType = "none"

type SessionOutput struct {
	ID           string
	Name         string
	Participants int
	Type         string
	CreatedAt    time.Time
	ClosedAt     *time.Time
}

func sessionTypeToString(t int) string {
	switch t {
	case model.SessionTypeArithmeticExpression:
		return sessionTypeArithmeticExpression
	default:
		return unexpectedSessionType
	}
}

func NewSessionOutput(data data.SessionData) SessionOutput {
	return SessionOutput{
		data.ID,
		data.Name,
		data.Participants,
		sessionTypeToString(data.Type),
		data.CreatedAt,
		data.ClosedAt,
	}
}
