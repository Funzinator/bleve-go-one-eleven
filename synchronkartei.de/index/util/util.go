package util

import (
	"fmt"

	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/search/query"
	"synchronkartei.de/index/document"
)

func NewAnalyzedWildcardQuery(wildcard string, filter []analysis.CharFilter) *query.WildcardQuery {
	if filter == nil {
		return query.NewWildcardQuery(wildcard)
	}

	input := []byte(wildcard)

	for _, cf := range filter {
		input = cf.Filter(input)
	}

	return query.NewWildcardQuery(string(input))
}

func FromIndexMap(fields map[string]interface{}) (s *document.Sprecher) {
	if fields == nil {
		return nil
	}

	s = &document.Sprecher{}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Conversion failed", r)
			s = nil
		}
	}()

	for field, value := range fields {
		switch field {
		case "Vorname":
			s.Vorname = value.(string)
		case "Nachname":
			s.Nachname = value.(string)
		case "Anrede":
			s.Anrede = value.(string)
		case "Zusatz":
			s.Zusatz = value.(string)
		case "PseudoVorname":
			s.PseudoVorname = value.(string)
		case "PseudoNachname":
			s.PseudoNachname = value.(string)
		case "Published":
			s.Published = value.(bool)
		case "AnzahlRollen":
			s.AnzahlRollen = uint64(value.(float64))
		}

	}
	return s
}
