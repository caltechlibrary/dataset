//
// Package dataset includes the operations needed for processing collections of JSON documents and their attachments.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2018, Caltech
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
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"testing"

	// 3rd Party Packages
	"github.com/blevesearch/bleve"
)

var (
	csvtable = `tind,issn,oclc,title,date
545035,,,Symposium on the electronic collection and dissemination of federal government information,
644902,,,California counts,1999
601242,,,"Manual of rules of the Committee on Ways and Means for the ... Congress, adopted January ...",
600949,0741-2665,,Annual report,1981
600622,0083-1565,,Annual report of the Librarian of Congress for the fiscal year ended ..,1939
598403,0161-4274,,National food review,1978
585896,,,Fiscal Year 1994 budget resolution : pros & cons,1993
545036,,,Symposium on the U.S. Office of Technology Assessment report : informing the nation : federal information dissemination in an electronic age,1989
544926,,,Smokestacks and silicon : regaining American's edge,1984
544391,,,High technology : public policies for the 1980's,1983
471221,,,General revenue sharing payment,
419899,0098-0404,,Air carrier traffic statistics,
418925,,,National food situation,
415750,0007-6597,,BCD : business conditions digest,
`
	cName  = path.Join("testdata", "search-test.ds")
	mbName = path.Join("testdata", "search-test.bmap")
	mName  = path.Join("testdata", "search-test.json")
	iName  = path.Join("testdata", "search-test.bleve")

	layouts = []int{
		PAIRTREE_LAYOUT,
		BUCKETS_LAYOUT,
	}
)

func TestBleveMapIndexingSearch(t *testing.T) {
	for _, cLayout := range layouts {
		// Remove stale collection and index
		os.RemoveAll(cName)
		os.RemoveAll(iName)

		// create the collection
		c, err := InitCollection(cName, cLayout)
		if err != nil {
			t.Errorf("%s", err)
			t.FailNow()
		}

		lines, err := c.ImportCSV(strings.NewReader(csvtable), -1, true, true, false)
		if err != nil {
			t.Errorf("Error import csvtable, %s", err)
			t.FailNow()
		}
		if lines != 16 {
			t.Errorf("Expected to import 16 rows, got %d", lines)
			t.FailNow()
		}
		if err := c.Indexer(iName, mbName, []string{}, 100); err != nil {
			t.Errorf("Can't create index %q, %s", iName, err)
			t.FailNow()
		}
		if err := c.Close(); err != nil {
			t.Errorf("Can't close index, %s", err)
			t.FailNow()
		}

		c, err = Open(cName)
		if err != nil {
			t.Errorf("%s", err)
			t.FailNow()
		}
		defer c.Close()

		// Run queries and test results
		opts := map[string]string{
			"result_fields": "*",
		}

		idxList, _, err := OpenIndexes([]string{iName})
		if err != nil {
			t.Errorf("Can't open index %s, %s", iName, err)
			t.FailNow()
		}

		results, err := Find(idxList.Alias, "600622", opts)
		if err != nil {
			t.Errorf("Find returned an error, %s", err)
			t.FailNow()
		}
		err = idxList.Close()
		if err != nil {
			t.Errorf("Failed to close index %s, %s", iName, err)
		}
		src, _ := json.Marshal(results)
		if len(results.Hits) != 1 {
			t.Errorf("unexpected results -> %s", src)
			t.FailNow()
		}
	}
}

func TestIndexingSearch(t *testing.T) {
	for _, cLayout := range layouts {
		// Remove stale collection and index
		os.RemoveAll(cName)
		os.RemoveAll(iName)

		// create the collection
		c, err := InitCollection(cName, cLayout)
		if err != nil {
			t.Errorf("%s", err)
			t.FailNow()
		}

		lines, err := c.ImportCSV(strings.NewReader(csvtable), -1, true, true, false)
		if err != nil {
			t.Errorf("Error import csvtable, %s", err)
			t.FailNow()
		}
		if lines != 16 {
			t.Errorf("Expected to import 16 rows, got %d", lines)
			t.FailNow()
		}
		if err = c.Indexer(iName, mName, []string{}, 100); err != nil {
			t.Errorf("Can't create index %q, %s", iName, err)
			t.FailNow()
		}
		if err = c.Close(); err != nil {
			t.Errorf("Can't close index, %s", err)
			t.FailNow()
		}

		c, err = Open(cName)
		if err != nil {
			t.Errorf("%s", err)
			t.FailNow()
		}
		defer c.Close()

		// Run queries and test results
		opts := map[string]string{
			"result_fields": "*",
		}

		idxList, _, err := OpenIndexes([]string{iName})
		if err != nil {
			t.Errorf("Can't open index %s, %s", iName, err)
			t.FailNow()
		}

		results, err := Find(idxList.Alias, "600622", opts)
		if err != nil {
			t.Errorf("Find returned an error, %s", err)
			t.FailNow()
		}
		err = idxList.Close()
		if err != nil {
			t.Errorf("failed to close index %s, %s", iName, err)
		}
		src, _ := json.Marshal(results)
		if len(results.Hits) != 1 {
			t.Errorf("unexpected results -> %s", src)
			t.FailNow()
		}
	}
}

func TestIndexerDeindexer(t *testing.T) {
	for _, cLayout := range layouts {
		cName := path.Join("testdata", "test_index.ds")
		indexName := path.Join("testdata", "test_index.bleve")
		indexMapName := path.Join("testdata", "test_index_map.json")
		os.RemoveAll(cName)
		os.RemoveAll(indexName)
		os.RemoveAll(indexMapName)

		testRecords := map[string]map[string]interface{}{}
		src := []byte(`{
    "gutenberg:21489": {"title": "The Secret of the Island", "formats": ["epub","kindle", "plain text", "html"], "authors": [{"given": "Jules", "family": "Verne"}], "url": "http://www.gutenberg.org/ebooks/21489", "categories": "fiction, novel"},
    "gutenberg:2488": { "title": "Twenty Thousand Leagues Under the Seas: An Underwater Tour of the World", "formats": ["epub","kindle","plain text"], "authors": [{ "given": "Jules", "family": "Verne" }], "url": "https://www.gutenberg.org/ebooks/2488", "categories": "fiction, novel"},
    "gutenberg:21839": { "title": "Sense and Sensibility", "formats": ["epub", "kindle", "plain text"], "authors": [{"given": "Jane", "family": "Austin"}], "url": "http://www.gutenberg.org/ebooks/21839", "categories": "fiction, novel" },
    "gutenberg:3186": {"title": "The Mysterious Stranger, and Other Stories", "formats": ["epub","kindle", "plain text", "html"], "authors": [{ "given": "Mark", "family": "Twain"}], "url": "http://www.gutenberg.org/ebooks/3186", "categories": "fiction, short story"},
    "hathi:uc1321060001561131": { "title": "A year of American travel - Narrative of personal experience", "formats": ["pdf"], "authors": [{"given": "Jessie Benton", "family": "Fremont"}], "url": "https://babel.hathitrust.org/cgi/pt?id=uc1.32106000561131;view=1up;seq=9", "categories": "non-fiction, memoir" }
}`)
		err := json.Unmarshal(src, &testRecords)
		if err != nil {
			t.Errorf("Can't unmarshal test records, %s", err)
			t.FailNow()
		}
		testRecordCount := len(testRecords)
		if testRecordCount != 5 {
			t.Errorf("Something went wrong with unmarshalling test records, expected 5 got %d", testRecordCount)
			t.FailNow()
		}

		src = []byte(`{
	"title": {
		"object_path": ".title"
	},
	"family": {
		"object_path": ".authors[:].family"
	},
	"categories": {
		"object_path": ".categories"
	}
}`)

		// Make sure our test definition is valid JSON!
		indexMap := map[string]interface{}{}
		err = json.Unmarshal(src, &indexMap)
		if err != nil {
			t.Errorf("Can't unmarshal test map, %s", err)
			t.FailNow()
		}

		err = ioutil.WriteFile(indexMapName, src, 0664)
		if err != nil {
			t.Errorf("Can't write test index map file, %s", err)
			t.FailNow()
		}
		// create the collection
		c, err := InitCollection(cName, cLayout)
		if err != nil {
			t.Errorf("%s", err)
			t.FailNow()
		}
		defer c.Close()

		// Populate our test collection
		for k, v := range testRecords {
			err = c.Create(k, v)
			if err != nil {
				t.Errorf("Can't create %s in %s, %s", k, cName, err)
				t.FailNow()
			}
		}

		err = c.Indexer(indexName, indexMapName, []string{}, 2)
		if err != nil {
			t.Errorf("%s", err)
			t.FailNow()
		}

		idxList, _, err := OpenIndexes([]string{indexName})
		if err != nil {
			t.Errorf("Can't open index %s, %s", indexName, err)
			t.FailNow()
		}

		queryString := `+family:"Verne"`
		results, err := Find(idxList.Alias, queryString, map[string]string{})
		if err != nil {
			t.Errorf("Can't find %q in index %s, %s", queryString, indexName, err)
			t.FailNow()
		}
		err = idxList.Close()
		if err != nil {
			t.Errorf("Failed to close index %s, %s", indexName, err)
		}
		src, err = json.Marshal(results.Hits)
		if err != nil {
			t.Errorf("Can't marshal results")
		}
		if results.Total != 2 {
			t.Errorf("Expected two results got %s", src)
			t.FailNow()
		}
		delKeys := []string{}
		for _, hit := range results.Hits {
			delKeys = append(delKeys, hit.ID)
		}
		err = c.Deindexer(indexName, delKeys, 0)
		if err != nil {
			t.Errorf("deindexer failed for %s, %s", indexName, err)
			t.FailNow()
		}
	}
}

func TestSearchSort(t *testing.T) {
	for _, cLayout := range layouts {
		cName := path.Join("testdata", "test_search_sort.ds")
		iName := path.Join("testdata", "test_search_sort.bleve")
		iMapName := path.Join("testdata", "test_search_sort_map.json")
		os.RemoveAll(cName)
		os.RemoveAll(iName)
		os.RemoveAll(iMapName)

		src := []byte(`{
	"title": {
		"object_path": ".title",
		"field_mapping": "text"
	},
	"created": {
		"object_path": ".created",
		"field_mapping": "datetime"
	}
}`)
		err := ioutil.WriteFile(iMapName, src, 0775)
		if err != nil {
			t.Errorf("Can't write %s, %s", iMapName, err)
			t.FailNow()
		}

		// create the collection
		c, err := InitCollection(cName, cLayout)
		if err != nil {
			t.Errorf("%s", err)
			t.FailNow()
		}
		defer c.Close()

		src = []byte(`{
		"one": {
			"title":"G one",
			"created": "2018-03-04"
		},
		"two": {
			"title": "F two",
			"created": "2018-04-04"
		},
		"three": {
			"title": "E three",
			"created": "2018-05-04"
		},
		"four": {
			"title": "D four",
			"created": "2018-06-04"
		},
		"five": {
			"title": "C five",
			"created": "2018-07-04"
		},
		"six": {
			"title": "B six",
			"created": "2018-08-04"
		},
		"seven": {
			"title": "A seven",
			"created": "2018-09-04"
		}
	}`)

		records := map[string]map[string]interface{}{}
		err = json.Unmarshal(src, &records)
		if err != nil {
			t.Errorf("Can't unmarshal src, %s", err)
		}

		for key, record := range records {
			err := c.Create(key, record)
			if err != nil {
				t.Errorf("Can't add %s to %s, %s", key, cName, err)
			}
		}
		keys := c.Keys()

		// Now we are ready to index our collection
		err = c.Indexer(iName, iMapName, keys, 100)
		if err != nil {
			t.Errorf("Can't create index %q, %s", iName, err)
			t.FailNow()
		}
		//c.CloseIndexes(iName)
		idxLists, _, err := OpenIndexes([]string{iName})
		if err != nil {
			t.Errorf("Can't open index %q, %s", iName, err)
			t.FailNow()
		}
		defer idxLists.Close()

		// Now we can test our sorting with a query of '*'
		options := map[string]string{
			"fields": ".created,.title",
			"sort":   "-.created,+.title",
		}
		results, err := Find(idxLists.Alias, "*", options)
		if err != nil {
			t.Errorf("Can't find '*' with sort option, %s", err)
			t.FailNow()
		}
		if results.Hits == nil || len(results.Hits) == 0 {
			src, _ := json.MarshalIndent(results, "", "    ")
			t.Errorf("Expected hits in result, %s", src)
		}
		expected := []string{
			"seven",
			"six",
			"five",
			"four",
			"three",
			"two",
			"one",
		}
		for i, hit := range results.Hits {
			src, err = json.MarshalIndent(hit, "", "    ")
			if err != nil {
				t.Errorf("Can't marshal search results (%d), %s", i, err)
				t.FailNow()
			}
			if expected[i] != hit.ID {
				t.Errorf("expected (%d) %q, got %q", i, expected[i], hit.ID)
			}
		}
	}
}

func setupSearchTests(m *testing.M) {
	var (
		err          error
		defn1, defn2 []byte
	)
	err = os.MkdirAll(path.Dir(mName), 0775)
	if err != nil {
		log.Fatalf("Can't create dir %s, %s", path.Dir(mName), err)
	}
	err = os.MkdirAll(path.Dir(mbName), 0775)
	if err != nil {
		log.Fatalf("Can't create dir %s, %s", path.Dir(mbName), err)
	}

	// Remove stale collection and index
	os.RemoveAll(mName)
	os.RemoveAll(mbName)

	// Build an Bleve index map to test with
	idxMap := bleve.NewIndexMapping()
	document := bleve.NewDocumentMapping()

	fieldMapping := bleve.NewTextFieldMapping()
	fieldMapping.Store = true
	fieldMapping.Index = true

	fieldMapping.Name = "tind"
	fieldMapping.Analyzer = "simple"
	document.AddFieldMapping(fieldMapping)

	fieldMapping.Name = "issn"
	fieldMapping.Analyzer = "keyword"
	document.AddFieldMapping(fieldMapping)

	fieldMapping.Name = "oclc"
	fieldMapping.Analyzer = "keyword"
	document.AddFieldMapping(fieldMapping)

	fieldMapping.Name = "title"
	fieldMapping.Analyzer = "en"
	document.AddFieldMapping(fieldMapping)

	fieldMapping.Name = "title"
	fieldMapping.Analyzer = "standard"
	document.AddFieldMapping(fieldMapping)

	datetimeFieldMapping := bleve.NewDateTimeFieldMapping()
	datetimeFieldMapping.Name = "year"
	document.AddFieldMapping(datetimeFieldMapping)

	idxMap.AddDocumentMapping("default", document)

	defn1, err = json.MarshalIndent(idxMap, "", "    ")
	if err != nil {
		log.Fatal("Can't marshal mapping for tests")
	}
	err = ioutil.WriteFile(mbName, defn1, 0666)
	if err != nil {
		log.Fatalf("Can't write %q, %s", mbName, err)
	}

	// Build an Simple index map to test with
	defn2 = []byte(`{
	"tind_id": {
		"object_path": ".tind",
		"field_mapping": "numeric",
		"store": true
	},
	"issn": {
		"object_path": ".issn",
		"field_mapping": "keyword",
		"store": true
	},
	"oclc": {
		"object_path": ".oclc",
		"field_mapping": "keyword",
		"store": true
	},
	"title": {
		"object_path": ".title",
		"field_mapping": "text",
		"analyzer": "lang",
		"lang":"en",
		"store": true
	},
	"title_simple": {
		"object_path": ".title",
		"field_mapping": "text",
		"analyzer": "simple",
		"store": true
	},
	"title_standard": {
		"object_path": ".title",
		"field_mapping": "text",
		"analyzer": "standard",
		"store": true
	},
	"title_keyword": {
		"object_path": ".title",
		"field_mapping": "text",
		"analyzer": "keyword",
		"store": true
	},
	"title_web": {
		"object_path": ".title",
		"field_mapping": "text",
		"analyzer": "web",
		"store": true
	},
	"year": {
		"object_path": ".date",
		"field_mapping": "numeric",
		"store": true
	}
}`)

	err = ioutil.WriteFile(mName, defn2, 0666)
	if err != nil {
		log.Fatalf("Can't write %q, %s", mName, err)
	}
}
