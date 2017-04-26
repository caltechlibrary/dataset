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
	"log"
	//"os"

	// 3rd Party packages
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
)

// ReadIndexMapFile reads a JSON document and converts it into a Bleve index map.
func ReadIndexMapFile(mapName string) (*mapping.IndexMappingImpl, error) {
	//FIXME: translate mapName into an appropriate mapping
	return bleve.NewIndexMapping(), nil
}

// Indexeri ingests all the records of a collection
func (c *Collection) Indexer(indexName string, indexMap *mapping.IndexMappingImpl) error {
	var (
		rec map[string]interface{}
	)
	//FIXME: if indexName exists use bleve.Open()
	index, err := bleve.New(indexName, indexMap)
	if err != nil {
		return err
	}

	// Get all the keys and index each record
	keys := c.Keys()
	for _, key := range keys {
		rec = map[string]interface{}{}
		if err = c.Read(key, rec); err == nil {
			log.Printf("DEBUG key %s, rec: %+v\n", key, rec)
			index.Index(key, rec)
		} else {
			log.Printf("Can't index %s, %s\n", key, err)
		}
	}
	return nil
}
