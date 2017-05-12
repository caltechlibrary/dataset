//
// index a tabla with a geopoint column and place
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

func main() {
	var (
		src = []byte(`
[
	{
		"place":"Caltech, Pasadena, California",
		"loc": {
			"lat": 34.1377,
			"lng": -118.1253
		}
	},
	{
		"place":"Colonia, Yap, Micronesia",
		"loc":{
			"lng": -138.1190079,
			"lat": 9.5165712
		}
	},
	{
		"place":"Ulithi, Yap, Micronesia",
		"loc": {
			"lng": -139.764123,
			"lat": 10.028822
		}
	},
	{
		"place":"Pechilemu,Chile",
		"loc":{
			"lat": -34.3972709,
			"lng": -72.0263149
		}
	}
]
`)
	)

	// Convert my JSON to array of maps
	data := []map[string]interface{}{}
	err := json.Unmarshal(src, &data)
	if err != nil {
		log.Fatalf("can't unpack demo data, %s", err)
	}

	// Define my index
	indexMapping := bleve.NewIndexMapping()

	placeMapping := bleve.NewTextFieldMapping()
	placeMapping.Store = true
	placeMapping.Analyzer = standard.Name

	locMapping := bleve.NewGeoPointFieldMapping()
	locMapping.Store = true

	indexMapping.DefaultMapping.AddFieldMappingsAt("place", placeMapping)
	indexMapping.DefaultMapping.AddFieldMappingsAt("loc", locMapping)

	indexMapping.DefaultAnalyzer = standard.Name

	// Build index
	os.RemoveAll("demo.bleve")
	idx, err := bleve.New("demo.bleve", indexMapping)
	if err != nil {
		log.Fatalf("can't create demo.bleve, %s", err)
	}

	// Index out data
	var (
		loc   interface{}
		place string
	)
	for key, val := range data {
		// Note: We want place to be a string for indexing purposes
		place, _ = val["place"].(string)
		// Note: we accept the map[string]interface{} for location (with keys lat, lng)  as is.
		loc, _ = val["loc"]
		idx.Index(fmt.Sprintf("%d", key), &map[string]interface{}{
			"place": place,
			"loc":   loc,
		})
	}
	idx.Close()

	// Open index for reading
	idx, err = bleve.Open("demo.bleve")
	if err != nil {
		log.Fatalf("can't open demo.bleve, %s", err)
	}

	for _, s := range []string{"place:Ulithi", "place:Colonia", "place:Pechilemu", "*"} {
		qry := bleve.NewQueryStringQuery(s)
		search := bleve.NewSearchRequestOptions(qry, 20, 0, false)

		// Set up the sort order
		search.SortBy([]string{"loc"})

		// Display field
		search.Fields = []string{"place", "loc"}

		// Search for place by location
		results, err := idx.Search(search)
		if err != nil {
			log.Fatalf("search error, %s", err)
		}
		fmt.Printf("search for %q\n%s\n", s, results)
	}
}
