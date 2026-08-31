package btree

type freelist struct {
	pager *pager
	head  pageID
}
