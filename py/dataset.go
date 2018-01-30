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

func main() {}
