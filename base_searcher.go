package chessminimax

type baseSearcher struct {
	searcher MoveSearcher
}

// SetSearcher ...
func (searcher *baseSearcher) SetSearcher(
	innerSearcher MoveSearcher,
) {
	searcher.searcher = innerSearcher
}
