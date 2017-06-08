//
// Read index table and each for last names with "Oh", "On", "Of", "The"
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
)

var (
	src = []byte(`
[
	{"email":"ho.li@scientists.example.org","first_name":"Ho","last_name":"Li"},
	{"email":"ho.lee@scientists.example.org","first_name":"Ho","last_name":"Lee"},
	{"email":"mark.oh@scientists.example.org","first_name":"Mark","last_name":"Oh"},
	{"email":"ann.oh@scientists.example.org","first_name":"Ann","last_name":"Oh"},
	{"email":"andor@digital-circus.example.org","first_name":"and","last_name":"or"},
	{"email":"andan@digital-circus.example.org","first_name":"and","last_name":"an"},
	{"email":"andoi@digital-circus.example.org","first_name":"And","last_name":"Oi"},
	{"email":"andon@digital-circus.example.org","first_name":"And","last_name":"On"},
	{"email":"andof@digital-circus.example.org","first_name":"And","last_name":"Of"},
	{"email":"offthe@digital-circus.example.org","first_name":"off","last_name":"the"},
	{"email":"t.j.turu@zbs.example.org","first_name":"t.j.","last_name":"turu"}
]`)
)

func main() {
	// Define my index
	indexMapping := bleve.NewIndexMapping()

	lastNameMapping := bleve.NewTextFieldMapping()
	lastNameMapping.Analyzer = simple.Name

	firstNameMapping := bleve.NewTextFieldMapping()
	firstNameMapping.Analyzer = simple.Name

	emailMapping := bleve.NewTextFieldMapping()
	emailMapping.Analyzer = simple.Name

	indexMapping.DefaultMapping.AddFieldMappingsAt("last_name", lastNameMapping)
	indexMapping.DefaultMapping.AddFieldMappingsAt("first_name", firstNameMapping)
	indexMapping.DefaultMapping.AddFieldMappingsAt("email", emailMapping)
	indexMapping.DefaultAnalyzer = simple.Name

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

	for _, s := range []string{"*", "Or", "the", "Li"} {
		qry := bleve.NewQueryStringQuery(s)
		search := bleve.NewSearchRequestOptions(qry, 20, 0, false)

		// Set up the sort order
		search.SortBy([]string{"last_name"})

		// Display field
		search.Fields = []string{"last_name", "first_name", "email"}

		// Search for last_name = "Or"
		results, err := idx.Search(search)
		if err != nil {
			log.Fatalf("search error, %s", err)
		}
		fmt.Printf("search for %q\n%s\n", s, results)
	}
}
