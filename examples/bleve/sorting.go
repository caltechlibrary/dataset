//
// Demo sorting titles using different analyzers
//
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	// Bleve imports
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/simple"
	"github.com/blevesearch/bleve/analysis/analyzer/standard"
)

var (
	src = []byte(`
[
	{"title":"The story of a cat"},
	{"title":"Once upon a time"},
	{"title":"Once upon a time, a commentary"},
	{"title":"Upon an appointed time, a commentary"},
	{"title":"They A, B, C of time book"},
	{"title":"Story of hats"},
	{"title":"Summers of misconduct"},
	{"title":"Stories of misconduct"},
	{"title":"Story misconduct"},
	{"title":"A story of time"}
]`)
)

func main() {
	// Define my index
	indexMapping := bleve.NewIndexMapping()

	titleMapping := bleve.NewTextFieldMapping()
	titleMapping.Analyzer = simple.Name

	indexMapping.DefaultMapping.AddFieldMappingsAt("title", titleMapping)
	indexMapping.DefaultAnalyzer = standard.Name

	// Convert my JSON to array of maps
	data := []map[string]string{}
	err := json.Unmarshal(src, &data)
	if err != nil {
		log.Fatalf("can't unpack demo data, %s", err)
	}

	// Build index
	os.RemoveAll("demo.bleve")
	idx, err := bleve.New("demo.bleve", indexMapping)
	if err != nil {
		log.Fatalf("can't create demo.bleve, %s", err)
	}
	for key, rec := range data {
		idx.Index(fmt.Sprintf("%d", key), rec)
	}
	idx.Close()

	// Open index for reading
	idx, err = bleve.Open("demo.bleve")
	if err != nil {
		log.Fatalf("can't open demo.bleve, %s", err)
	}

	for _, s := range []string{"story", `title:"story"`, `title:"story of"`} {
		qry := bleve.NewQueryStringQuery(s)
		fmt.Printf("Relevence order")
		search := bleve.NewSearchRequestOptions(qry, 20, 0, false)

		// Display field
		search.Fields = []string{"title"}

		results, err := idx.Search(search)
		if err != nil {
			log.Fatalf("search error, %s", err)
		}

		fmt.Printf("search for %q\n%s\n", s, results)
		fmt.Printf("Title sort order")
		search = bleve.NewSearchRequestOptions(qry, 20, 0, false)
		// Set up the sort order
		search.SortBy([]string{"title"})

		// Display field
		search.Fields = []string{"title"}

		results, err = idx.Search(search)
		if err != nil {
			log.Fatalf("search error, %s", err)
		}
		fmt.Printf("search for %q\n%s\n", s, results)
	}
}
