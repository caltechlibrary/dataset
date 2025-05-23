// Package dataset includes the operations needed for processing collections of JSON documents and their attachments.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2022, Caltech
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
package dataset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	// Caltech Library packages
	"github.com/caltechlibrary/dsv1"
	"github.com/caltechlibrary/semver"
)

// sniffVersionNumber tries to get the dataset version
// string from collection.json file. Returns a semver
// or nil (on failure)
func sniffVersionNumber(cName string) *semver.Semver {
	collection := path.Join(cName, "collection.json")
	src, err := ioutil.ReadFile(collection)
	if err != nil {
		return nil
	}
	o := map[string]interface{}{}
	err = json.Unmarshal(src, &o)
	if err != nil {
		return nil
	}
	version, ok := o["dataset"]
	if ok {
		s := version.(string)
		sv, err := semver.ParseString(s)
		if err == nil && sv != nil {
			return sv
		}
		return nil
	}
	return nil
}

// Analyzer checks the collection version and analyzes current
// state of collection reporting on errors.
//
// NOTE: the collection MUST BE CLOSED when Analyzer is called otherwise
// the results will not be accurate.
func Analyzer(cName string, verbose bool) error {
	var (
		eCnt int
		wCnt int
		data interface{}
		c    *Collection
		err  error
	)

	collectionPath := cName
	// Make sure collection exists
	_, err = os.Stat(collectionPath)
	if err != nil {
		return err
	}

	// Check for collections.json file.
	collection := path.Join(collectionPath, "collection.json")
	_, err = os.Stat(collection)
	if err != nil {
		return err
	}

	// Sniff the version number of the collection
	v2 := semver.NewSemver(2, 0, 0, "")
	currentSV := sniffVersionNumber(cName)
	if currentSV != nil && semver.Less(currentSV, v2) {
		repairLog(verbose, "WARNING: %q is a version 1 dataset collection", cName)
		return dsv1.Analyzer(cName, verbose)
	}

	// Make sure the JSON documents in the collectionPath can be
	// parsed.
	files, err := os.ReadDir(collectionPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		filename := file.Name()
		isDir := file.IsDir()
		if !isDir && strings.HasSuffix(filename, ".json") {
			// Make sure we can JSON parse the file
			docPath := path.Join(collectionPath, filename)
			if src, err := os.ReadFile(docPath); err == nil {
				if err := json.Unmarshal(src, &data); err == nil {
					// release the memory
					data = nil
				} else {
					return fmt.Errorf("error parsing %s, %s", docPath, err)
				}
			} else {
				return fmt.Errorf("error opening %s, %s", docPath, err)
			}
		}
	}

	// NOTE: Check to see if we have a codemeta.json file
	codemeta := path.Join(collectionPath, "codemeta.json")
	_, err = os.Stat(codemeta)
	if err != nil {
		repairLog(verbose, "WARNING: Missing codemeta.json\n")
		wCnt++
	}

	// Now try to open the collection ...
	c, err = Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()

	if c.StoreType == SQLSTORE {
		_, err := c.SQLStore.Keys()
		if err != nil {
			return fmt.Errorf("WARNING: The collection.json's .name and .dsn_uri to not match the database connection and expected table name.")
		}
		return nil
	}

	if c.StoreType != SQLSTORE {
		return fmt.Errorf("analyzer only SQL storage")
	}
	// FIXME: Need to run table check on primary and history tables, look for orphaned records that are not delete in primary table but exist in history.
	// If history is enabled then check if the number of unique keys in history is greater than or equal the number of keys in primary table.
	// Make sure we can join the primary and history tables.

	// FIXME: need to check for attachments and make sure they are recorded OK

	if eCnt > 0 || wCnt > 0 {
		return fmt.Errorf("%d errors, %d warnings detected", eCnt, wCnt)
	}
	return nil
}

// FixMissingCollectionJson will scan the collection directory
// and environment making an educated guess to type of
// collection collection type
func FixMissingCollectionJson(cName string) error {
	collectionJson := path.Join(cName, "collection.json")
	//FIXME: Need to check to see if we should default to the old pairtree or SQLite3 database
	dsnURI := os.Getenv("DATASET_DSN_URI")
	sqlitePath := path.Join(cName, "collection.db")
	storeType := ""
	history := false
	if _, err := os.Stat(sqlitePath); err == nil {
		storeType = SQLSTORE
		dsnURI = "sqlite://collection.db"
	}
	if dsnURI != "" {
		storeType = SQLSTORE
		history = true
	}
	if storeType == "" {
		return fmt.Errorf("unable to determine storage type for %q", cName)
	}
	c := &Collection{}
	c.Name = path.Base(cName)
	c.DatasetVersion = Version
	c.StoreType = storeType
	c.DsnURI = dsnURI
	c.History = history
	src, err := JSONMarshalIndent(c, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to encode %q, %s", collectionJson, err)
	}
	return ioutil.WriteFile(collectionJson, src, 0664)
}

// Repair a SQLite3 base collection.
func repairSqlite3(c *Collection) error {
	// Fixme see if SQLite3 is installed and in the path
	cmdPath, err := exec.LookPath("sqlite3")
	if err != nil {
		return fmt.Errorf("sqlite3 needs to be installed to repair %q, %s", c.Name, err)
	}
	dbName := path.Join(c.workPath, strings.TrimPrefix(c.DsnURI, "sqlite://"))
	cmdDump := exec.Command(cmdPath, dbName, ".dump")
	src, err := cmdDump.Output()
	if err != nil {
		return fmt.Errorf("failed to execute %q, %s", dbName + " .dump", err)
	}
	os.Rename(dbName, dbName + "-broken")
	cmdRestore := exec.Command(cmdPath, dbName)
	buffer := bytes.Buffer{}
	buffer.Write(src)
	cmdRestore.Stdin = &buffer
	cmdRestore.Stdout = os.Stdout
	cmdRestore.Stderr = os.Stderr
	if err := cmdRestore.Run(); err != nil {
		return fmt.Errorf("failed to retore %q, %s", dbName, err)
	}
	return nil
}

// Repair takes a collection name and calls
// walks the pairtree and repairs collection.json as appropriate.
//
// NOTE: the collection MUST BE CLOSED when repair is called otherwise
// the repaired collection may revert.
func Repair(cName string, verbose bool) error {
	var (
		c   *Collection
		err error
	)
	// Sniff the version number of the collection and delegate
	// if needed.
	v3 := semver.NewSemver(3, 0, 0, "")
	currentSV := sniffVersionNumber(cName)
	if currentSV != nil && semver.Less(currentSV, v3) {
		return fmt.Errorf("cannot repair %q dataset collections", currentSV.String())
	}

	collectionJson := path.Join(cName, "collection.json")
	// Check to see if we find a collection.json, if not see if we
	// can make a educated guess
	if _, err := os.Stat(collectionJson); err != nil {
		err := FixMissingCollectionJson(cName)
		if err != nil {
			return err
		}
	}
	// See if we can open a collection, if not then create an empty struct
	c, err = Open(cName)
	if err != nil {
		repairLog(verbose, "Open %s error, %s, attempting to re-create collection.json", cName, err)
		err = os.WriteFile(path.Join(cName, "collection.json"), []byte("{}"), 0664)
		if err != nil {
			repairLog(verbose, "Can't re-initilize %s, %s", cName, err)
			return err
		}
		repairLog(verbose, "Attempting to re-open %s", cName)

		c, err = Open(cName)
		if err != nil {
			repairLog(verbose, "Failed to re-open %s, %s", cName, err)
			return err
		}
	}
	defer c.Close()

	c.DatasetVersion = Version


	repairLog(verbose, "Saving metadata for %s", cName)
	// Save the collections' operational metadata
	c.Repaired = time.Now().Format("2006-01-02")
	src, err := JSONMarshalIndent(c, "", "    ")
	filename := path.Join(c.workPath, "collection.json")
	err = ioutil.WriteFile(filename, src, 664)
	if err != nil {
		return err
	}
	return nil
}

//
// Helper functions
//

func repairLog(verbose bool, rest ...interface{}) {
	if verbose == true {
		s := rest[0].(string)
		log.Printf(s, rest[1:]...)
	}
}
