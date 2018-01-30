//
// py/dataset.go is a C shared library targetting support in Python for dataset
//
// @author R. S. Doiel, <rsdoiel@library.caltech.edu>

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
package main

import (
	"C"
	"encoding/json"
	"log"

	// Caltech Library Packages
	"github.com/caltechlibrary/dataset"
)

//export init_collection
func init_collection(name *C.char) C.int {
	collectionName := C.GoString(name)
	log.Printf("creating %s\n", collectionName)
	_, err := dataset.InitCollection(collectionName)
	if err != nil {
		log.Printf("Cannot create collection %s, %s", collectionName, err)
		return C.int(0)
	}
	log.Printf("%s initialized", collectionName)
	return C.int(1)
}

//export create_record
func create_record(name, key, src *C.char) C.int {
	collectionName := C.GoString(name)
	k := C.GoString(key)
	v := C.GoString(src)

	c, err := dataset.Open(collectionName)
	if err != nil {
		log.Printf("Cannot open collection %s, %s", collectionName, err)
		return C.int(0)
	}
	defer c.Close()

	m := map[string]interface{}{}
	err = json.Unmarshal(v, &m)
	if err != nil {
		log.Printf("Can't unmarshal %s, %s", k, err)
		return C.int(0)
	}

	err = c.Create(k, m)
	if err != nil {
		log.Printf("Create %s failed, %s", k, err)
		return C.int(0)
	}
	return C.int(1)
}

//export read_record
func read_record(name, key *C.char) *C.char {
	collectionName := C.GoString(name)
	k := C.GoString(key)

	c, err := dataset.Open(collectionName)
	if err != nil {
		log.Printf("Cannot open collection %s, %s", collectionName, err)
		return C.int(0)
	}
	defer c.Close()

	m := map[string]interface{}{}
	err = c.Read(k, m)
	if err != nil {
		log.Printf("Can't read %s, %s", k, err)
		return C.CString("")
	}

	// return a JSON string to Python
	src, err := json.Marshal(m)
	if err != nil {
		log.Printf("Can't marshal %s, %s", k, err)
		return C.CString("")
	}
	return C.CString(string(src))
}

//export update_record
func update_record(name, key, src *C.char) C.int {
	collectionName := C.GoString(name)
	k := C.GoString(key)
	v := C.GoString(src)

	c, err := dataset.Open(collectionName)
	if err != nil {
		log.Printf("Cannot open collection %s, %s", collectionName, err)
		return C.int(0)
	}
	defer c.Close()

	m := map[string]interface{}{}
	err = json.Unmarshal(v, &m)
	if err != nil {
		log.Printf("Can't unmarshal %s, %s", k, err)
		return C.int(0)
	}

	err = c.Update(k, m)
	if err != nil {
		log.Printf("Update %s failed, %s", k, err)
		return C.int(0)
	}
	return C.int(1)
}

//export delete_record
func delete_record(name, key *C.char) C.int {
}

//export keys
func keys(name, key, filter, sort_by *C.char) *C.char {
}

//export count
func count(name, key, filter *C.char) C.int {
}

func main() {}
