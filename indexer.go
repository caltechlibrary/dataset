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
	"io/ioutil"
	"log"
	"os"

	// Caltech Library packages
	"github.com/caltechlibrary/datatools/dotpath"

	// 3rd Party packages
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
)

// readIndexDefinition reads in a JSON document and converts it into a record map and a Bleve index mapping.
func readIndexDefinition(mapName string) (map[string]string, *mapping.IndexMappingImpl, error) {
	var (
		src []byte
		err error
	)

	if src, err = ioutil.ReadFile(mapName); err != nil {
		return nil, nil, err
	}

	cfg := map[string]string{}
	if err := json.Unmarshal(src, &cfg); err != nil {
		return nil, nil, err
	}

	//FIXME: convert definition into an appropriate index map

	return cfg, bleve.NewIndexMapping(), nil
}

// recordMapToIndexRecord takes the definition map, Unmarshals the JSON record and
// renders a new map[string]string that is ready to be indexed.
func recordMapToIndexRecord(recordMap map[string]string, src []byte) (map[string]interface{}, error) {
	raw := map[string]interface{}{}
	idxMap := map[string]interface{}{}
	err := json.Unmarshal(src, &raw)
	if err != nil {
		return nil, err
	}
	// Copy the dot path elements to new smaller map
	for pName, dPath := range recordMap {
		if val, err := dotpath.Eval(dPath, raw); err == nil {
			idxMap[pName] = val
		}
	}
	return idxMap, nil
}

// Indexeri ingests all the records of a collection
func (c *Collection) Indexer(idxName string, idxMapName string) error {
	var (
		idx bleve.Index
		err error
	)
	recordMap, idxMap, err := readIndexDefinition(idxMapName)
	if err != nil {
		return err
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
