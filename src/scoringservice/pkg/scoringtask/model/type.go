package model

const TaskTypeArithmeticExpression = "expressions"

func IsTypeValid(t string) bool {
	return t == TaskTypeArithmeticExpression
}
