package chessminimax

type searcherHolder struct {
	searcher BoundedMoveSearcher
}

func (holder *searcherHolder) SetSearcher(
	searcher BoundedMoveSearcher,
) {
	holder.searcher = searcher
}
