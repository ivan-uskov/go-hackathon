package model

import "testing"

func TestTranslateType(t *testing.T) {
	if v, ok := TranslateType(TaskTypeArithmeticExpressionAlias); !ok || v != TaskTypeArithmeticExpression {
		t.Error("Translate expressions type not working")
	}

	if v, ok := TranslateType("another"); ok || v != TaskTypeInvalid {
		t.Error("Translate invalid type working")
	}
}
