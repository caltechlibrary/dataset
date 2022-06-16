package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//
// NOTE: Examples routes using curl assume the "host" has been
// configured for the default "localhost:8485".
//

// ApiVersion returns the version of the web service running.
// This will normally be the same version of dataset you installed.
//
// ```shell
//    curl -X GET http://localhost:8485/api/version
// ```
//
func ApiVersion(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
	fmt.Fprintf(w, "%s %s", api.AppName, api.Version)
}

// Collections returns a list of dataset collections supported
// by the running web service.
//
// ```shell
//    curl -X GET http://localhost:8485/api/collections
// ```
//
func Collections(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
	collections := []string{}
	w.Header().Add("Content-Type", "application/json")
	if len(api.CMap) > 0 {
		for k := range api.CMap {
			collections = append(collections, k)
		}
		src, err := json.MarshalIndent(collections, "", "     ")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "%s", src)
		return
	}
	fmt.Fprintf(w, "[]")
}

// Collection returns the codemeta JSON for a specific collection.
// Example collection name "journals.ds"
//
// ```shell
//    curl -X GET http://localhost:8485/api/collection/journals.ds
// ```
//
func Codemeta(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

// ApiKeys returns the available keys in a collection as a JSON array.
// Example collection name "journals.ds"
//
// ```shell
//    curl -X GET http://localhost:8485/api/journals.ds/keys
// ```
//
func Keys(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
	if c, ok := api.CMap[cName]; ok {
		keys, err := c.Keys()
		if err != nil {
			log.Printf("c.Keys() returned error %s", err)
			http.NotFound(w, r)
			return
		}
		src, err := json.MarshalIndent(keys, "", "    ")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// Set header to application/json
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", src)
		return
	}
	fmt.Fprintf(w, "ApiKeys(w, r, api, %q, %q, %s) not implemented", cName, verb, strings.Join(options, " "))
	return
}

// Create deposit a JSON object in the collection for a given key.
//
// In this example the json document is in the working directory called
// "record-123.json" and the environment variable KEY holds the document
// key which is the string "123".
//
// ```shell
//    KEY="123"
//    curl -X POST http://localhost:8585/api/journals.ds/object/$KEY
//         -H "Content-Type: application/json" \
//          --data-binary "@./record-123.json"
// ```
//
func Create(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
	defer r.Body.Close()
	if len(options) != 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	key := options[0]

	if c, ok := api.CMap[cName]; ok {
		src, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		o := map[string]interface{}{}
		err = json.Unmarshal(src, &o)
		if err != nil {
			log.Printf("unmarshal error %+v, %s", o, err)
			http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
			return
		}
		if err := c.Create(key, o); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// Set header to application/json
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok", "key": %q, "action": "created"}`, key)
		return
	}
	http.NotFound(w, r)
	return
}

// Read retrieves a JSON object from the collection for a given key.
//
// In this example the json retrieved will be called "record-123.json"
// and the environment variable KEY holds the document key
// as a string "123".
//
// ```shell
//    KEY="123"
//    curl -o "record-123.json" -X GET \
//         http://localhost:8585/api/journals.ds/object/$KEY
// ```
//
func Read(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
	if len(options) != 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	key := options[0]

	if c, ok := api.CMap[cName]; ok {
		o := map[string]interface{}{}
		err := c.Read(key, o)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		src, err := json.MarshalIndent(o, "", "    ")
		if err != nil {
			log.Printf("marshal error %+v, %s", o, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// Set header to application/json
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `%s`, src)
		return
	}
	http.NotFound(w, r)
	return
}

// Update replaces a JSON object in the collection for a given key.
//
// In this example the json document is in the working directory called
// "record-123.json" and the environment variable KEY holds the document
// key which is the string "123".
//
// ```shell
//    KEY="123"
//    curl -X PUT http://localhost:8585/api/journals.ds/object/$KEY
//         -H "Content-Type: application/json" \
//          --data-binary "@./record-123.json"
// ```
//
func Update(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
	defer r.Body.Close()
	if len(options) != 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	key := options[0]

	if c, ok := api.CMap[cName]; ok {
		src, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		o := map[string]interface{}{}
		err = json.Unmarshal(src, &o)
		if err != nil {
			log.Printf("unmarshal error %+v, %s", o, err)
			http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
			return
		}
		if err := c.Update(key, o); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// Set header to application/json
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok", "key": %q, "action": "updated"}`, key)
		return
	}
	http.NotFound(w, r)
	return
}

// Delete removes a JSON object from the collection for a given key.
//
// In this example the environment variable KEY holds the document
// key which is the string "123".
//
// ```shell
//    KEY="123"
//    curl -X DELETE http://localhost:8585/api/journals.ds/object/$KEY
// ```
//
func Delete(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
	defer r.Body.Close()
	if len(options) != 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	key := options[0]

	if c, ok := api.CMap[cName]; ok {
		if err := c.Delete(key); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// Set header to application/json
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok", "key": %q, "action": "deleted"}`, key)
		return
	}
	http.NotFound(w, r)
	return
}

//
// The following routes handle attachments
//

func Attachments(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func Attach(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func Retrieve(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func Prune(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

//
// The following routes handle frames
//

// HasFrame checks a collection for a frame by its name
//
//```shell
//    FRM_NAME="name"
//    curl -X GET http://localhost:8585/api/journals.ds/has-frame/$FRM_NAME
//```
//
func HasFrame(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	// Get Frame name
	frameName := ""
	if len(options) > 0 {
		frameName = options[0]
	}
	// Get collection
	c, ok := api.CMap[cName]
	if ok {
		if c.HasFrame(frameName) {
			fmt.Fprint(w, "true")
			return
		}
	}
	// Check if frame is in collection
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func FrameCreate(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	// Get Frame name
	frameName := ""
	if len(options) > 0 {
		frameName = options[0]
	}
	// Process post
	src, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("FrameCreate, Bad Request %s %q %s", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	m := map[string][]string{}
	if err := json.Unmarshal(src, &m); err != nil {
		log.Printf("FrameCreate, Bad Request %s %q %s", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	keys := []string{}
	dotPaths := []string{}
	labels := []string{}
	if data, ok := m["dot_paths"]; ok {
		dotPaths = data[:]
	}
	if data, ok := m["labels"]; ok {
		labels = data[:]
	}
	if data, ok := m["keys"]; ok {
		keys = data[:]
	}
	// Get collection
	c, ok := api.CMap[cName]
	if ok {
		if _, err := c.FrameCreate(frameName, keys, dotPaths, labels, false); err != nil {
			log.Printf("FrameCreate, Bad Request %s %q %s", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		fmt.Fprintf(w, "OK")
		return
	}
	// Check if frame is in collection
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	return
}

func Frames(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func FrameKeys(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func FrameDef(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func FrameObjects(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func FrameRefresh(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func FrameReframe(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func FrameDelete(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

//***************************************************
// The following routes handle JSON object versions
//***************************************************

func ObjectVersions(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func ReadVersion(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func DeleteVersion(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

//**************************************************
// The following routes handle attachment versions
//**************************************************

func AttachmentVersions(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func RetrieveVersion(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func PruneVersion(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
