//
// Package dataset includes the operations needed for processing collections of JSON documents and their attachments.
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
	"io/ioutil"
	"os"
	"strings"
	"testing"
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
	cName = "testdata/search-test"
	mName = "testdata/search-test.json"
	iName = "testdata/search-test.bleve"

	defn1 = []byte(`{
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
)

func TestIndexingSearch(t *testing.T) {
	// Remove stale collection and index
	os.RemoveAll(cName)
	os.RemoveAll(iName)

	// create the collection
	c, err := create(cName, generateBucketNames("ab", 2))
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}

	lines, err := c.ImportCSV(strings.NewReader(csvtable), true, -1, true, false)
	if err != nil {
		t.Errorf("Error import csvtable, %s", err)
		t.FailNow()
	}
	if lines != 16 {
		t.Errorf("Expected to import 16 rows, got %d", lines)
		t.FailNow()
	}
	// Build an index to test with
	if err := ioutil.WriteFile(mName, defn1, 0666); err != nil {
		t.Errorf("Can't write %q, %s", mName, err)
		t.FailNow()
	}
	if err := c.Indexer(iName, mName, 100, []string{}); err != nil {
		t.Errorf("Can't create index %q, %s", iName, err)
		t.FailNow()
	}
	if err := c.Close(); err != nil {
		t.Errorf("Can't close index, %s", err)
		t.FailNow()
	}
}

func TestSearch(t *testing.T) {
	c, err := Open(cName)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	defer c.Close()

	// Run queries and test results
	opts := map[string]string{
		"result_fields": "*",
	}

	idx, _, err := OpenIndexes([]string{iName})
	if err != nil {
		t.Errorf("Can't open index %s, %s", iName, err)
		t.FailNow()
	}

	results, err := Find(os.Stderr, idx, []string{"600622"}, opts)
	if err != nil {
		t.Errorf("Find returned an error, %s", err)
		t.FailNow()
	}
	src, _ := json.Marshal(results)
	if len(results.Hits) != 1 {
		t.Errorf("unexpected results -> %s", src)
		t.FailNow()
	}
}
