package main

import (
	"C"
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
	err = json.Marshal(v, &m)
	if err != nil {
		log.Printf("Can't marshal %s, %s", k, err)
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
}

//export update_record
func update_record(name, key, src *C.char) C.int {
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
