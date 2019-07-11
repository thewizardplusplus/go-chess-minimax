package chessminimax

type searcherHolder struct {
	searcher MoveSearcher
}

func (holder *searcherHolder) SetSearcher(
	searcher MoveSearcher,
) {
	holder.searcher = searcher
}
