package dataset

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestMain(m *testing.M) {
	appName := path.Base(os.Args[0])
	settings := []byte(`
{
    "host": "localhost:8485",
    "collections": {
        "t1": {
            "dataset": "t1.ds",
            "keys": true,
            "create": true,
            "read": true,
            "updated": true,
            "delete": true
        }
    }
}
`)
	dName := "testdata"
	fName := path.Join(dName, "test-settings.json")
	if _, err := os.Stat(dName); os.IsNotExist(err) {
		os.MkdirAll("testdata", 0777)
	}
	cName := path.Join(dName, "t1.ds")
	if _, err := os.Stat(cName); os.IsNotExist(err) {
		if _, err := InitCollection(cName); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create %s, %s\n", cName, err)
			fmt.Fprintf(os.Stderr, "Aborting %s\n", appName)
			os.Exit(1)
		}
	}
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		if err := ioutil.WriteFile(fName, settings, 0666); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create %s\n", fName)
			fmt.Fprintf(os.Stderr, "Aborting %s\n", appName)
			os.Exit(1)
		}
	}
	fmt.Fprintf(os.Stdout, "Test PID: %d\n", os.Getpid())
	os.Exit(m.Run())
}
