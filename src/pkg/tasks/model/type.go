package model

const TaskTypeArithmeticExpressionAlias = "expressions"

const TaskTypeArithmeticExpression = 1
const TaskTypeInvalid = 0

func TranslateType(t string) (int, bool) {
	if t == TaskTypeArithmeticExpressionAlias {
		return TaskTypeArithmeticExpression, true
	}

	return TaskTypeInvalid, false
}
