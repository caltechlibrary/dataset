//
// Package dataset is a go package for managing JSON documents stored on disc
//
// Author R. S. Doiel, <rsdoiel@library.caltech.edu>
//
// Copyright (c) 2017, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package dataset

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/datatools/dotpath"

	// 3rd Party packages
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/analysis/analyzer/simple"
	"github.com/blevesearch/bleve/analysis/analyzer/standard"
	"github.com/blevesearch/bleve/analysis/analyzer/web"
	"github.com/blevesearch/bleve/analysis/lang/ar"
	"github.com/blevesearch/bleve/analysis/lang/ca"
	"github.com/blevesearch/bleve/analysis/lang/cjk"
	"github.com/blevesearch/bleve/analysis/lang/ckb"
	"github.com/blevesearch/bleve/analysis/lang/de"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/analysis/lang/es"
	"github.com/blevesearch/bleve/analysis/lang/fa"
	"github.com/blevesearch/bleve/analysis/lang/fr"
	"github.com/blevesearch/bleve/analysis/lang/hi"
	"github.com/blevesearch/bleve/analysis/lang/it"
	"github.com/blevesearch/bleve/analysis/lang/pt"
	"github.com/blevesearch/bleve/geo"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/search/highlight/highlighter/ansi"
	"github.com/blevesearch/bleve/search/highlight/highlighter/html"
)

var (
	// languagesSupported by Analyzer
	languagesSupported = map[string]string{
		"ar":  ar.AnalyzerName,
		"ca":  ca.ArticlesName,
		"cjk": cjk.AnalyzerName,
		"ckb": ckb.AnalyzerName,
		"de":  de.AnalyzerName,
		"en":  en.AnalyzerName,
		"es":  es.AnalyzerName,
		"fa":  fa.AnalyzerName,
		"fr":  fr.AnalyzerName,
		"hi":  hi.AnalyzerName,
		"it":  it.AnalyzerName,
		"pt":  pt.AnalyzerName,
	}

	// supportedNamedTimeFormats for named Golang time strings (e.g. RFC3339) plus
	// ones that are for convienence (e.g. mysqldate, mysqldatetime)
	supportedNamedTimeFormats = map[string]string{
		"ansic":    "Mon Jan _2 15:04:05 2006",
		"unixdate": "Mon Jan _2 15:04:05 MST 2006",
		"rubydate": "Mon Jan 02 15:04:05 -0700 2006",
		"rfc822":   "02 Jan 06 15:04 MST",
		// RFC822 with numeric zone
		"rfc822z": "02 Jan 06 15:04 -0700",
		"rfc850":  "Monday, 02-Jan-06 15:04:05 MST",
		"rfc1123": "Mon, 02 Jan 2006 15:04:05 MST",
		// RFC1123 with numeric zone
		"rfc1123z":    "Mon, 02 Jan 2006 15:04:05 -0700",
		"rfc3339":     "2006-01-02T15:04:05Z07:00",
		"rfc3339nano": "2006-01-02T15:04:05.999999999Z07:00",
		"kitchen":     "3:04PM",
		// Handy time stamps.
		"stamp":         "Jan _2 15:04:05",
		"stampmilli":    "Jan _2 15:04:05.000",
		"stampmicro":    "Jan _2 15:04:05.000000",
		"stampnano":     "Jan _2 15:04:05.000000000",
		"mysqldate":     "2006-01-02",
		"mysqldatetime": "2006-01-02 14:05:05",
	}
)

// isTrueString normlize string values to true if they are "true", "t", "1" case insensitive
// otherwise it returns false
func isTrueString(s string) bool {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return v
}

// readIndexDefinition reads in a JSON document and converts it into a record map and a Bleve index mapping.
func readIndexDefinition(mapName string) (map[string]string, *mapping.IndexMappingImpl, error) {
	var (
		src []byte
		err error
	)

	if src, err = ioutil.ReadFile(mapName); err != nil {
		return nil, nil, err
	}

	definitions := map[string]map[string]string{}
	if err := json.Unmarshal(src, &definitions); err != nil {
		return nil, nil, fmt.Errorf("error unpacking definition: %s", err)
	}

	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultAnalyzer = simple.Name

	//NOTE: convert definition into an appropriate index mappings and record data paths
	cfg := map[string]string{}
	var fieldMap *mapping.FieldMapping

	for fieldName, defn := range definitions {
		if dPath, ok := defn["object_path"]; ok == true {
			cfg[fieldName] = dPath
		}
		if fieldType, ok := defn["field_mapping"]; ok == true {
			switch fieldType {
			case "numeric":
				fieldMap = bleve.NewNumericFieldMapping()
			case "datetime":
				fieldMap = bleve.NewDateTimeFieldMapping()
			case "boolean":
				fieldMap = bleve.NewBooleanFieldMapping()
			case "geopoint":
				fieldMap = bleve.NewGeoPointFieldMapping()
			default:
				fieldMap = bleve.NewTextFieldMapping()
			}
		} else {
			fieldMap = bleve.NewTextFieldMapping()
		}
		if sVal, ok := defn["store"]; ok == true {
			if isTrueString(sVal) == true {
				fieldMap.Store = true
			} else {
				fieldMap.Store = false
			}
		}
		if analyzerType, ok := defn["analyzer"]; ok == true {
			switch analyzerType {
			case "keyword":
				fieldMap.Analyzer = keyword.Name
			case "simple":
				fieldMap.Analyzer = simple.Name
			case "standard":
				fieldMap.Analyzer = standard.Name
			case "web":
				fieldMap.Analyzer = web.Name
			case "lang":
				if langCode, ok := defn["lang"]; ok == true {
					if langAnalyzer, ok := languagesSupported[langCode]; ok == true {
						fieldMap.Analyzer = langAnalyzer
					}
				}
			}
		}
		if sVal, ok := defn["include_in_all"]; ok == true {
			if isTrueString(sVal) == true {
				fieldMap.IncludeInAll = true
			} else {
				fieldMap.IncludeInAll = false
			}
		}
		if sVal, ok := defn["include_term_vectors"]; ok == true {
			if isTrueString(sVal) == true {
				fieldMap.IncludeTermVectors = true
			} else {
				fieldMap.IncludeTermVectors = false
			}
		}
		if sVal, ok := defn["date_format"]; ok == true {
			if fmt, ok := supportedNamedTimeFormats[strings.ToLower(strings.TrimSpace(sVal))]; ok == true {
				fieldMap.DateFormat = fmt
			} else {
				fieldMap.DateFormat = strings.TrimSpace(sVal)
			}
		}
		indexMapping.DefaultMapping.AddFieldMappingsAt(fieldName, fieldMap)
	}
	return cfg, indexMapping, nil
}

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

// recordMapToIndexRecord takes the definition map, Unmarshals the JSON record and
// renders a new map[string]string that is ready to be indexed.
func recordMapToIndexRecord(recordMap map[string]string, src []byte) (map[string]interface{}, error) {
	idxMap := map[string]interface{}{}

	raw, err := dotpath.JSONDecode(src)
	if err != nil {
		return nil, err
	}
	// Copy the dot path elements to new smaller map
	for pName, dPath := range recordMap {
		if val, err := dotpath.Eval(dPath, raw); err == nil {
			switch val.(type) {
			case json.Number:
				if i, err := (val.(json.Number)).Int64(); err == nil {
					idxMap[pName] = i
				} else if f, err := (val.(json.Number)).Float64(); err == nil {
					idxMap[pName] = f
				} else {
					idxMap[pName] = (val.(json.Number)).String()
				}
			default:
				if pt, err := geoPointParse(val.(string)); err == nil {
					idxMap[pName] = pt
				} else {
					idxMap[pName] = val
				}
			}
		}
	}
	return idxMap, nil
}

// Indexer ingests all the records of a collection applying the definition
// creating or updating a Bleve index. Returns an error.
func (c *Collection) Indexer(idxName string, idxMapName string) error {
	var (
		idx bleve.Index
		err error
	)
	recordMap, idxMap, err := readIndexDefinition(idxMapName)
	if err != nil {
		return fmt.Errorf("failed to read index definition %s, %s", idxMapName, err)
	}

	//NOTE: if indexName exists use bleve.Open() instead of bleve.New()
	if _, e := os.Stat(idxName); os.IsNotExist(e) {
		idx, err = bleve.New(idxName, idxMap)
	} else {
		idx, err = bleve.Open(idxName)
	}
	if err != nil {
		return err
	}

	// Get all the keys and index each record
	keys := c.Keys()
	cnt := 0
	for i, key := range keys {
		if src, err := c.ReadAsJSON(key); err == nil {
			if rec, err := recordMapToIndexRecord(recordMap, src); err == nil {
				idx.Index(key, rec)
				cnt++
				if (cnt % 100) == 0 {
					log.Printf("%d records indexed", cnt)
				}
			}
		} else {
			log.Printf("%d, can't index %s, %s", i, key, err)
		}
	}
	log.Printf("%d total records indexed", cnt)
	return nil
}

// Find takes a Bleve index name and query string, opens the index, and writes the
// results to the os.File provided. Function returns an error if their are problems.
func Find(out io.Writer, indexNames []string, queryStrings []string, options map[string]string) (*bleve.SearchResult, error) {
	// Opening all our indexes
	var (
		idxAlias bleve.IndexAlias
		size     int
		from     int
		explain  bool
		err      error
	)
	for i, idxName := range indexNames {
		idx, err := bleve.Open(idxName)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			idxAlias = bleve.NewIndexAlias(idx)
		} else {
			idxAlias.Add(idx)
		}
	}

	if sVal, ok := options["from"]; ok == true {
		from, err = strconv.Atoi(sVal)
		if err != nil {
			return nil, err
		}
	}
	if sVal, ok := options["size"]; ok == true {
		size, err = strconv.Atoi(sVal)
		if err != nil {
			return nil, err
		}
	} else {
		size = 12
	}
	if sVal, ok := options["explain"]; ok == true {
		explain, err = strconv.ParseBool(sVal)
		if err != nil {
			return nil, err
		}
	} else {
		explain = false
	}

	//Note: find uses the Query String Query, it'll join queryStrings with a space
	query := bleve.NewQueryStringQuery(strings.Join(queryStrings, " "))
	search := bleve.NewSearchRequestOptions(query, size, from, explain)

	// Handle various options modifying search
	if sVal, ok := options["highlight"]; ok == true {
		if isTrueString(sVal) == true {
			if sHighlighter, ok := options["highlighter"]; ok == true {
				switch strings.TrimSpace(strings.ToLower(sHighlighter)) {
				case "ansi":
					search.Highlight = bleve.NewHighlightWithStyle(ansi.Name)
				case "html":
					search.Highlight = bleve.NewHighlightWithStyle(html.Name)
				default:
					log.Printf("Unknown highlighter, %q, using defaults", sHighlighter)
					search.Highlight = bleve.NewHighlight()
				}
			} else {
				search.Highlight = bleve.NewHighlight()
			}
		}
	}

	if sVal, ok := options["result_fields"]; ok == true {
		if strings.Contains(sVal, ":") == true {
			search.Fields = strings.Split(sVal, ":")
		} else {
			search.Fields = []string{sVal}
		}
	}

	if sVal, ok := options["sort_by"]; ok == true {
		if strings.Contains(sVal, ":") == true {
			search.SortBy(strings.Split(sVal, ":"))
		} else {
			search.SortBy([]string{sVal})
		}
	}

	// Run the query and process results
	results, err := idxAlias.Search(search)
	if err != nil {
		return nil, err
	}
	return results, nil
}
