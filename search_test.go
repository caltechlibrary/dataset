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

func TestIndexing(t *testing.T) {
	// Remove stale collection and index
	os.RemoveAll(cName)
	os.RemoveAll(iName)

	// Create the collection
	c, err := Create(cName, GenerateBucketNames("ab", 2))
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
	if err := c.Indexer(iName, mName); err != nil {
		t.Errorf("Can't create index %q, %s", iName, err)
		t.FailNow()
	}
	if err := c.Close(); err != nil {
		t.Errorf("Can't close index, %s", err)
		t.FailNow()
	}
}

func TestSearch(t *testing.T) {
	// Run queries and test results
	opts := map[string]string{
		"result_fields": "*",
	}

	results, err := Find(os.Stderr, []string{iName}, []string{"600622"}, opts)
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