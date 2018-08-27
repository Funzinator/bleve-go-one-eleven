package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"synchronkartei.de/index/document"
	"synchronkartei.de/index/indexer"
)

func main() {
	path := "/tmp/test-index"
	index, err := indexer.NewBleveIndex(path)

	needIndex := true

	if err != nil {
		log.Printf("unable to open index, trying to create a new one")
		needIndex = true
		index, err = indexer.CreateBleveIndex(path)
	}
	if err != nil {
		log.Fatalf("error while opening index: %s", err)
	}

	var wg sync.WaitGroup
	if needIndex {
		sprecherListe := createData()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			wg.Add(1)
			defer wg.Done()
			log.Printf("Start Indexing")
			startIndex := time.Now()
			i := 0
			for _, sprecher := range sprecherListe {
				err = index.IndexSprecher(sprecher)
				if err != nil {
					log.Fatalf("Error while indexing: %s", err)
				}

				i++
			}
			duration := time.Since(startIndex)
			if i > 0 {
				x := int64(duration) / int64(i)
				log.Printf("Indexing of %d elements finished in %s (%s per element)", i, duration.String(), time.Duration(x))
			} else {
				log.Printf("No elements indexed in %s", duration.String())
			}
		}()
	}

	time.Sleep(2 * time.Second)

	sprecherResult, err := index.SearchSprecher("maria")
	if err != nil {
		log.Fatalf("Search failed: %s", err)
	}
	for _, s := range sprecherResult {
		fmt.Printf("%s %s, %s %s (%d)\n", s.Nachname, s.Zusatz, s.Anrede, s.Vorname, s.AnzahlRollen)
	}

	log.Printf("waiting for indexing")
	wg.Wait()

}

func createData() []*document.Sprecher {
	r := []*document.Sprecher{}
	for i := 0; i < 9000; i++ {
		s := document.Sprecher{}
		s.Vorname = fmt.Sprintf("aaaaaaaaaaaaaaaaaaaaaaaaaa%d", i)
		s.Nachname = fmt.Sprintf("bbbbbbbbbbbbbbbbbbbbbbbbbbbbb%d", i)
		s.ID = fmt.Sprintf("t%d", i)

		r = append(r, &s)
	}

	return r
}
