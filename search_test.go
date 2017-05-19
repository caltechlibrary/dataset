package dataset

import (
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
	iName = "testdata/search-test.bleve"
)

func TestSearch(t *testing.T) {
	// Remove stale collection and index
	os.RemoveAll(cName)
	os.RemoveAll(iName)

	// Create the collection
	c, err := Create(cName, GenerateBucketNames("ab", 2))
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	defer c.Close()

	lines, err := c.ImportCSV(strings.NewReader(csvtable), false, 0, false, false)
	if err != nil {
		t.Errorf("Error import csvtable, %s", err)
		t.FailNow()
	}
	if lines != 16 {
		t.Errorf("Expected to import 16 rows, got %d", lines)
		t.FailNow()
	}
	// Build an index to test with
	t.Errorf("building test index not implemented")
	// Run queries and test results
	t.Errorf("run tests of queries and results, not implemented")
}
