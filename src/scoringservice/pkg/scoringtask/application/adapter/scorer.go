package adapter

type Scorer interface {
	Score(url string) int
}

type ScorerFactory interface {
	GetScorer(taskType string) (Scorer, error)
}
