package output

import (
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/query/data"
	"go-hackathon/src/hackathonservice/pkg/hackathon/model"
	"time"
)

const typeArithmeticExpression = "arithmetic expression"
const unexpectedType = "none"

type HackathonOutput struct {
	ID           string
	Name         string
	Participants int
	Type         string
	CreatedAt    time.Time
	ClosedAt     *time.Time
}

func typeToString(t int) string {
	switch t {
	case model.HackathonTypeArithmeticExpression:
		return typeArithmeticExpression
	default:
		return unexpectedType
	}
}

func NewHackathonOutput(data data.HackathonData) HackathonOutput {
	return HackathonOutput{
		data.ID,
		data.Name,
		data.Participants,
		typeToString(data.Type),
		data.CreatedAt,
		data.ClosedAt,
	}
}
