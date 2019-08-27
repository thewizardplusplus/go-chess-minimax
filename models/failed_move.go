package models

// FailedMove ...
type FailedMove struct {
	Move  ScoredMove
	Deep  int
	Error error
}
