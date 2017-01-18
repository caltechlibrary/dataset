//
// Package dataset is a go package for managing JSON documents stored on disc
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
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
// Examples:
//
//   // Create a collection
//   collection, err := dataset.Create("mystuff", "dataset", GenerateBucketNames("abc", 3))
//   if err != nil {
//       log.Fatalf("%s", err)
//	 }
//   defer collection.Close()
//   // Add a record
//   record := map[string]string{"name":"freda","email":"freda@inverness.example.org"}
//   if err := collection.Create("freda", record); err != nil {
//       log.Fatalf("%s", err)
//   }
//   // Read a record
//   if err := collection.Read("freda", record); err != nil {
//       log.Fatalf("%s", err)
//   }
//   // Update a record
//   record["email"] = "freda@zbs.example.org"
//   if err := collection.Update("freda", record); err != nil {
//       log.Fatalf("%s", err)
//   }
//   // Delete a record
//   if err := collection.Delete("freda"); err != nil {
//       log.Fatalf("%s", err)
//   }
//
package dataset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// countToBucketID turns a count assigns it to a letter sequence (e.g. 0-999 is aa, 1000 - 1999 is ab, etc)
func countToBucketID(i int, bucketNames []string) string {
	bucketsize := len(bucketNames)
	// Calculate bucket number
	bucketIndex := i % bucketsize
	return bucketNames[bucketIndex]
}

// generateCombinations from an alphabet and length
//
// this function is based on example at https://play.golang.org/p/0bWDCibSUJ
func generateCombinations(alphabet string, length int) <-chan string {
	c := make(chan string)

	// Starting a separate goroutine that will create all the combinations,
	// feeding them to the channel c
	go func(c chan string) {
		defer close(c) // Once the iteration function is finished, we close the channel

		addLetter(c, "", alphabet, length) // We start by feeding it an empty string
	}(c)

	return c // Return the channel to the calling function
}

// addLetter adds a letter to the combination to create a new combination.
// This new combination is passed on to the channel before we call AddLetter once again
// to add yet another letter to the new combination in case length allows it
//
// this function is based on gist at https://play.golang.org/p/0bWDCibSUJ
func addLetter(c chan string, combo string, alphabet string, length int) {
	// Check if we reached the length limit
	// If so, we just return without adding anything
	if length <= 0 {
		return
	}

	var newCombo string
	for _, ch := range alphabet {
		newCombo = combo + string(ch)
		c <- newCombo
		addLetter(c, newCombo, alphabet, length-1)
	}
}

// intToBucketName converts an integer to a bucket (using round robin pick via modulo)
func intToBucketName(val int, buckets []string) string {
	size := len(buckets)
	return buckets[(val % size)]
}

// GenerateBucketNames provides a list of permutations of requested length to use as bucket names
func GenerateBucketNames(alphabet string, length int) []string {
	l := []string{}
	for combo := range generateCombinations(alphabet, length) {
		if len(combo) == length {
			l = append(l, combo)
		}
	}
	return l
}

type Collection struct {
	Dataset string            `json:"dataset"`
	Buckets []string          `json:"buckets"`
	Name    string            `json:"name"`
	Keys    map[string]string `json:"keys"`
	dirty   bool
}

// CreateCollection - create a new collection structure on disc
// name should be filesystem friendly
func Create(name string, bucketNames []string) (*Collection, error) {
	c := new(Collection)
	c.Name = path.Base(name)
	c.Dataset = path.Dir(name)
	c.Buckets = bucketNames
	c.Keys = map[string]string{}
	// Make the collection directory
	if err := os.MkdirAll(path.Join(c.Dataset, c.Name), 0770); err != nil {
		return nil, err
	}
	// Add the JSON metadata for collection
	src, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(path.Join(c.Dataset, c.Name, "collection.json"), src, 0664); err != nil {
		return nil, err
	}
	return c, nil
}

// Open reads in a collection's metadata and returns and new collection structure and err
func Open(name string) (*Collection, error) {
	dataPath := path.Dir(name)
	fname := path.Base(name)
	src, err := ioutil.ReadFile(path.Join(dataPath, fname, "collection.json"))
	if err != nil {
		return nil, err
	}
	c := new(Collection)
	if err := json.Unmarshal(src, &c); err == nil {
		return c, err
	}
	return c, nil
}

// Delete an entire collection
func Delete(name string) error {
	if err := os.RemoveAll(name); err != nil {
		return err
	}
	return nil
}

// Close closes a collection, writing the updated keys to disc
func (c *Collection) Close() error {
	if c.dirty == true {
		src, err := json.Marshal(c)
		if err != nil {
			return fmt.Errorf("Can't marshal %s, %s", c.Name, err)
		}
		if err := ioutil.WriteFile(path.Join(c.Dataset, c.Name, "collection.json"), src, 0664); err != nil {
			return err
		}
	}
	// Cleanup c so it can't accidentally get reused
	c.Dataset = ""
	c.Buckets = []string{}
	c.Name = ""
	c.Keys = map[string]string{}
	c.dirty = false
	return nil
}

// Create adds a JSON blob to a collection, if problem return an error
// name must be unique
func (c *Collection) Create(name string, data interface{}) error {
	_, keyExists := c.Keys[name]
	if keyExists == true {
		return fmt.Errorf("%q already exists", name)
	}
	if len(c.Buckets) == 0 {
		return fmt.Errorf("collection is not valid, zero buckets")
	}
	bucketName := intToBucketName(len(c.Buckets), c.Buckets)
	p := path.Join(c.Dataset, c.Name, bucketName)
	err := os.MkdirAll(p, 0770)
	if err != nil {
		return fmt.Errorf("WriteJSON() mkdir %s", p, err)
	}
	src, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("WriteJSON() JSON encode %s/%s, %s", p, name, err)
	}
	// We've almost made it, save the key's bucket name and write the blob to bucket
	c.Keys[name] = path.Join(bucketName)
	return ioutil.WriteFile(path.Join(p, name), src, 0664)
}

// Read finds the record in a collection, updates data and returns error if a problem is encountered
// name must exist or an error is returned
func (c *Collection) Read(name string, data interface{}) error {
	bucketName, ok := c.Keys[name]
	if ok != true {
		return fmt.Errorf("%q does not exist", name)
	}
	p := path.Join(c.Dataset, c.Name, bucketName)
	src, err := ioutil.ReadFile(path.Join(p, name))
	if err != nil {
		return err
	}
	err = json.Unmarshal(src, &data)
	if err != nil {
		return fmt.Errorf("json.Unmarshal() failed for %s, %s", name, err)
	}
	return nil
}

// Update a JSON blob in a collection (note: Record must exist or returns an error )
func (c *Collection) Update(name string, data interface{}) error {
	bucketName, ok := c.Keys[name]
	if ok != true {
		return fmt.Errorf("%q does not exist", name)
	}
	p := path.Join(c.Dataset, c.Name, bucketName)
	err := os.MkdirAll(p, 0770)
	if err != nil {
		return fmt.Errorf("WriteJSON() mkdir %s", p, err)
	}
	src, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("WriteJSON() JSON encode %s/%s, %s", p, name, err)
	}
	return ioutil.WriteFile(path.Join(p, name), src, 0664)
}

// Delete removes a JSON blob from a collection
func (c *Collection) Delete(name string) error {
	bucketName, ok := c.Keys[name]
	if ok != true {
		return fmt.Errorf("%q key not found", name)
	}
	p := path.Join(c.Dataset, c.Name, bucketName, name)
	if err := os.Remove(p); err != nil {
		return fmt.Errorf("Error removing %q, %s", p, err)
	}
	return nil
}
