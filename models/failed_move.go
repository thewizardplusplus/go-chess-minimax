package models

// FailedMove ...
type FailedMove struct {
	Move  ScoredMove
	Error error
}
