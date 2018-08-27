package indexer

import "synchronkartei.de/index/document"

type Index interface {
	SearchSprecher(query string) ([]*document.Sprecher, error)
	IndexSprecher(sprecher ...*document.Sprecher) error
}
