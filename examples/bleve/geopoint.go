//
// index a tabla with a geopoint column and place
//
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	// Bleve imports
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/standard"
	"github.com/blevesearch/bleve/geo"
)

// geoPointParse takes a string in latitude,longitude decimal format (e.g. 34.1358302,-118.127694 for Caltech, Pasadena, California)
// and converts them to a Morton hash used by Bleve
func geoPointParse(s string) (uint64, error) {
	if strings.Contains(s, ",") == false {
		return 0, fmt.Errorf("Missing coordinate pair in %q", s)
	}
	latlon := strings.Split(s, ",")
	if len(latlon) != 2 {
		return 0, fmt.Errorf("Wrong number for pair in %q", s)
	}
	lat, err := strconv.ParseFloat(latlon[0], 64)
	if err != nil {
		return 0, fmt.Errorf("Can't parse latitude in %q, %s", s, err)
	}
	lon, err := strconv.ParseFloat(latlon[1], 64)
	if err != nil {
		return 0, fmt.Errorf("Can't parse longitude in %q, %s", s, err)
	}
	return geo.MortonHash(lon, lat), nil
}

func main() {
	var (
		src = []byte(`
[
	{
		"place":"Colonia, Yap, Micronesia",
		"loc":"9.5165712,138.1190079"
	},
	{
		"place":"Ulithi, Yap, Micronesia",
		"loc":"10.028822,139.764123"
	},
	{
		"place":"Pechilemu,Chile",
		"loc":"-34.3972709,-72.0263149"
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
	for key, val := range data {
		place, _ := val["place"].(string)
		loc, _ := val["loc"].(string)
		geoHash, err := geoPointParse(loc)
		if err != nil {
			log.Fatal(err)
		}

		idx.Index(fmt.Sprintf("%d", key), &map[string]interface{}{
			"place": place,
			"loc":   geoHash,
		})
	}
	idx.Close()

	// Open index for reading
	idx, err = bleve.Open("demo.bleve")
	if err != nil {
		log.Fatalf("can't open demo.bleve, %s", err)
	}

	for _, s := range []string{"place:Ulithi", "place:Colonia", "place:Pechilemu"} {
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
