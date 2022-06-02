//
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
//
package dataset

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

//
// Clone initializes a new collection based on the list of keys provided.
// If the keys list is empty all the objects are copied from one collection
// to the other. The collections do not need to be the same storage type.
//
// NOTE: The cloned copy is not open after cloning is complete.
//
// ```
//   newName, dsnURI :=
//      "new-collection.ds", "sqlite://new-collection.ds/collection.db"
//   c, err := dataset.Open("old-collection.ds")
//   if err != nil {
//       ... // handle error
//   }
//   defer c.Close()
//   nc, err := c.Clone(newName, dsnURI, []string{}, false)
//   if err != nil {
//       ... // handle error
//   }
//   defer nc.Close()
// ```
//
func (c *Collection) Clone(cloneName string, cloneDsnURI string, keys []string, verbose bool) error {
	nc, err := Init(cloneName, cloneDsnURI)
	if err != nil {
		return fmt.Errorf("initializing clone %q, %s", cloneName, err)
	}
	defer nc.Close()
	if len(keys) == 0 {
		keys, err = c.Keys()
		if err != nil {
			return err
		}
	}
	rptSize := 2500
	errCnt := 0
	tot := len(keys)
	for i, key := range keys {
		obj := map[string]interface{}{}
		err := c.Read(key, obj)
		if err != nil {
			log.Printf("(%d/%d) failed to read %q from %q, %s", i, tot, key, c.Name, err)
			errCnt += 1
			continue
		}
		if err := nc.Create(key, obj); err != nil {
			log.Printf("(%d/%d) failed to write %q to %q, %s", i, tot, key, cloneName, err)
			errCnt += 1
		}
		if verbose && ((i % rptSize) == 0) {
			log.Printf("%d/%d objects cloned into %q", i, tot, cloneName)
		}
	}
	if verbose {
		log.Printf("%d/%d objects cloned into %q", tot, tot, cloneName)
	}
	if errCnt > 0 {
		return fmt.Errorf("%d errors cloning %q to %q", errCnt, c.Name, cloneName)
	}
	return nil
}

//
// CloneSample initializes two new collections based on a training and test // sampling of the keys in the original collection.  If the keys list is
// empty all the objects are used for creating the taining and test
// sample collections.  The collections do not need to be the same
// storage type.
//
// NOTE: The cloned copy is not open after cloning is complete.
//
// ```
//   trainingSetSize := 10000
//   trainingName, trainingDsnURI :=
//      "training.ds", "sqlite://training.ds/collection.db"
//   testName, testDsnURI := "test.ds", "sqlite://test.ds/collection.db"
//   c, err := dataset.Open("old-collection")
//   if err != nil {
//       ... // handle error
//   }
//   defer c.Close()
//   nc, err := c.CloneSample(trainingName, trainingDsnURI,
//                            testName, testDsnURI, []string{},
//                            trainingSetSize, false)
//   if err != nil {
//       ... // handle error
//   }
//   defer nc.Close()
// ```
//
func (c *Collection) CloneSample(trainingName string, trainingDsnURI string, testName string, testDsnURI string, keys []string, sampleSize int, verbose bool) error {
	if sampleSize < 1 {
		return fmt.Errorf("sample size should be greater than zero")
	}
	// Copy the keys provided
	var err error
	workKeys := keys[:]
	if len(keys) == 0 {
		workKeys, err = c.Keys()
		if err != nil {
			return fmt.Errorf("Can't read keys from %q, %s", c.Name, err)
		}
	}
	// so a random sort on the work key list
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	random.Shuffle(len(workKeys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	// Split randomly sorted into two key lists
	trainingKeys := workKeys[0:sampleSize]
	testKeys := workKeys[sampleSize:]

	// Clone into respective collections
	if err := c.Clone(trainingName, trainingDsnURI, trainingKeys, verbose); err != nil {
		return err
	}
	if testName != "" {
		if err := c.Clone(testName, testDsnURI, testKeys, verbose); err != nil {
			return err
		}
	}
	return nil
}
