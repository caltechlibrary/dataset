//
// index a tabla with a numeric column and title
//
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	// Bleve imports
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/standard"
)

var (
	src = []byte(`
[
	{
		"title":"The Fourth Tower of Inverness",
		"year":1972
	},
	{
		"title":"Moon Over Morocco",
		"year":1974
	},
	{
		"title":"The Ah-Ha Phenomenon",
		"year":1977
	},
	{
		"title":"The Incredible Adventures of Jack Flanders",
		"year":1978
	}
]
`)
)

func main() {
	// Define my index
	indexMapping := bleve.NewIndexMapping()

	titleMapping := bleve.NewTextFieldMapping()
	titleMapping.Analyzer = standard.Name

	yearMapping := bleve.NewNumericFieldMapping()

	indexMapping.DefaultMapping.AddFieldMappingsAt("title", titleMapping)
	indexMapping.DefaultMapping.AddFieldMappingsAt("year", yearMapping)

	indexMapping.DefaultAnalyzer = standard.Name

	// Convert my JSON to array of maps
	data := []map[string]interface{}{}
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

	for _, s := range []string{"1972", "year:>=1974", "year:<=1978", "+year:>=1972 +year:<1975"} {
		qry := bleve.NewQueryStringQuery(s)
		search := bleve.NewSearchRequestOptions(qry, 20, 0, false)

		// Set up the sort order
		search.SortBy([]string{"title"})

		// Display field
		search.Fields = []string{"title", "year"}

		// Search for title = "Or"
		results, err := idx.Search(search)
		if err != nil {
			log.Fatalf("search error, %s", err)
		}
		fmt.Printf("search for %q\n%s\n", s, results)
	}
}
