package indexer

import (
	"log"
	"time"

	"github.com/blevesearch/bleve"
	"synchronkartei.de/index/document"
	"synchronkartei.de/index/util"
)

type BleveIndex struct {
	index bleve.Index
}

func NewBleveIndex(path string) (Index, error) {
	index, err := bleve.Open(path)
	if err != nil {
		return nil, err
	}

	return &BleveIndex{index: index}, nil
}

func CreateBleveIndex(path string) (Index, error) {
	indexMapping := bleve.NewIndexMapping()

	index, err := bleve.New(path, indexMapping)
	if err != nil {
		return nil, err
	}

	return &BleveIndex{index: index}, nil
}

func (index *BleveIndex) SearchSprecher(q string) ([]*document.Sprecher, error) {
	matchQuery := bleve.NewWildcardQuery(q)
	publishedQuery := bleve.NewBoolFieldQuery(true)
	publishedQuery.SetField("Published")

	query := bleve.NewBooleanQuery()
	query.AddMust(publishedQuery)
	query.AddMust(matchQuery)

	search := bleve.NewSearchRequest(query)
	search.Size = 200
	search.Fields = (&document.Sprecher{}).StoredFields()

	startSearch := time.Now()
	sr, err := index.index.Search(search)
	if err != nil {
		return nil, err
	}

	sprecherResult := []*document.Sprecher{}
	for _, hit := range sr.Hits {
		doc := util.FromIndexMap(hit.Fields)
		if doc != nil {
			sprecherResult = append(sprecherResult, doc)
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("Search finished in %s", time.Since(startSearch))

	return sprecherResult, nil
}

func (index *BleveIndex) IndexSprecher(sprecher ...*document.Sprecher) error {
	for _, s := range sprecher {
		err := index.index.Index(s.ID, s)
		if err != nil {
			return err
		}
	}

	return nil
}
