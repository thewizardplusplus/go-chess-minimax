package chessminimax

type searcherHolder struct {
	searcher BoundedMoveSearcher
}

func newSearcherHolder(
	searcher BoundedMoveSearcher,
) *searcherHolder {
	return &searcherHolder{searcher}
}

func (holder *searcherHolder) SetSearcher(
	searcher BoundedMoveSearcher,
) {
	holder.searcher = searcher
}
