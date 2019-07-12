package chessminimax

type searcherHolder struct {
	searcher MoveSearcher
}

// SetSearcher ...
func (holder *searcherHolder) SetSearcher(
	searcher MoveSearcher,
) {
	holder.searcher = searcher
}
