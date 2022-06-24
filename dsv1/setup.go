package dsv1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
	"time"

	// Caltech Library packages
	"github.com/caltechlibrary/dataset/v2/pairtree"
)

func SetupV1TestCollection(cName string, records map[string]map[string]interface{}) error {
	// Create collection.json using v1 structures
	if len(cName) == 0 {
		return fmt.Errorf("missing a collection name")
	}
	// Make root directory to hold collection.
	os.MkdirAll(cName, 0775)
	// Generate a v1 collection
	c := new(Collection)
	// Save the date and time
	dt := time.Now()
	// date and time is in RFC3339 format
	c.Created = dt.Format(time.RFC3339)
	// When is a date in YYYY-MM-DD format (can be approximate)
	// e.g. 2021, 2021-01, 2021-01-02
	c.When = dt.Format("2006-01-02")
	c.DatasetVersion = "1.1.1"
	c.Name = path.Base(cName)
	c.Version = "v0.0.0"
	userinfo, err := user.Current()
	if err == nil {
		if userinfo.Name != "" {
			c.Who = []string{userinfo.Name}
		} else {
			c.Who = []string{userinfo.Username}
		}
	}
	if len(c.Who) > 0 {
		c.What = fmt.Sprintf("A dataset (%s) collection initilized on %s by %s.", Version, dt.Format("Monday, January 2, 2006 at 3:04pm MST."), strings.Join(c.Who, ", "))
	} else {
		c.What = fmt.Sprintf("A dataset %s collection initilized on %s", Version, dt.Format("Monday, January 2, 2006 at 3:04pm MST.."))
	}
	c.workPath = cName
	c.KeyMap = make(map[string]string)

	if err := c.saveMetadata(); err != nil {
		return err
	}
	// Now populate with some test records records.
	for key, obj := range records {
		src, err := json.MarshalIndent(obj, "", "    ")
		pair := pairtree.Encode(key)
		pPath := path.Join("pairtree", pair)
		c.KeyMap[key] = pPath
		filename := path.Join(cName, pPath, key+".json")
		if err := os.MkdirAll(path.Dir(filename), 0775); err != nil {
			return err
		}
		if err := ioutil.WriteFile(filename, src, 0664); err != nil {
			return err
		}
		// Update collection.json to have keys for added records
		if err = c.saveMetadata(); err != nil {
			return err
		}
	}
	return nil
}
