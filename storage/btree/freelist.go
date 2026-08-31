package engine

type freelist struct {
	pager *pager
	head  pageID
}
