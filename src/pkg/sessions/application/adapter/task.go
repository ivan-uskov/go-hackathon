package adapter

type TaskAdapter interface {
	TranslateType(t string) (int, bool)
}
