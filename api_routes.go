package dataset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path"
)

const (
	verbose = true // set to true for debugging, false otherwise
)

//
// NOTE: Examples routes using curl assume the "host" has been
// configured for the default "localhost:8485".
//

// statusOK returns a JSON object indicating the status of a request
// is OK.
func statusIsOK(w http.ResponseWriter, statusCode int, cName string, key string, action string, target string) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	m := map[string]string{}
	m["status"] = "OK"
	if cName != "" {
		m["collection"] = cName
	}
	if key != "" {
		m["key"] = key
	}
	if action != "" {
		m["action"] = action
	}
	if target == "" {
		m["target"] = target
	}
	src, _ := JSONMarshal(m)
	fmt.Fprintf(w, "%s", src)
}

// ApiVersion returns the version of the web service running.
// This will normally be the same version of dataset you installed.
//
// ```shell
//
//	curl -X GET http://localhost:8485/api/version
//
// ```
func ApiVersion(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "%s %s", api.AppName, api.Version)
}

// Collections returns a list of dataset collections supported
// by the running web service.
//
// ```shell
//
//	curl -X GET http://localhost:8485/api/collections
//
// ```
func Collections(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
	collections := []string{}
	w.Header().Add("Content-Type", "application/json")
	if len(api.CMap) > 0 {
		for k := range api.CMap {
			collections = append(collections, k)
		}
		src, err := JSONMarshalIndent(collections, "", "     ")
		if err != nil {
			log.Printf("marshal error %+v, %s", collections, err)
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
//
//	curl -X GET http://localhost:8485/api/collection/journals.ds
//
// ```
func Codemeta(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
	// Get collection
	c, ok := api.CMap[cName]
	if ok {
		src, err := c.Codemeta()
		if err != nil {
			log.Printf("Codemeta, not found for %s, %s", cName, err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "%s", src)
		return
	}
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	return
}

// Keys returns the available keys in a collection as a JSON array.
// Example collection name "journals.ds"
//
// ```shell
//
//	curl -X GET http://localhost:8485/api/journals.ds/keys
//
// ```
func Keys(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
	if c, ok := api.CMap[cName]; ok {
		keys, err := c.Keys()
		if err != nil {
			log.Printf("c.Keys() returned error %s", err)
			http.NotFound(w, r)
			return
		}
		src, err := JSONMarshalIndent(keys, "", "    ")
		if err != nil {
			log.Printf("marshal error %+v, %s", keys, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// Set header to application/json
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", src)
		return
	}
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	return
}

// Create deposit a JSON object in the collection for a given key.
//
// In this example the json document is in the working directory called
// "record-123.json" and the environment variable KEY holds the document
// key which is the string "123".
//
// ```shell
//
//	KEY="123"
//	curl -X POST http://localhost:8585/api/journals.ds/object/$KEY
//	     -H "Content-Type: application/json" \
//	      --data-binary "@./record-123.json"
//
// ```
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
			log.Printf("Create, Bad Request %s %q %s", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		o := map[string]interface{}{}
		err = json.Unmarshal(src, &o)
		if err != nil {
			log.Printf("Create, unmarshal error %+v, %s", o, err)
			http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
			return
		}
		if err := c.Create(key, o); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		statusIsOK(w, http.StatusCreated, cName, key, "created", "")
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
//
//	KEY="123"
//	curl -o "record-123.json" -X GET \
//	     http://localhost:8585/api/journals.ds/object/$KEY
//
// ```
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
		src, err := JSONMarshalIndent(o, "", "    ")
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
//
//	KEY="123"
//	curl -X PUT http://localhost:8585/api/journals.ds/object/$KEY
//	     -H "Content-Type: application/json" \
//	      --data-binary "@./record-123.json"
//
// ```
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
			log.Printf("Update, Bad Request %s %q %s", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		o := map[string]interface{}{}
		err = json.Unmarshal(src, &o)
		if err != nil {
			log.Printf("Update, unmarshal error %+v, %s", o, err)
			http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
			return
		}
		if err := c.Update(key, o); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		statusIsOK(w, http.StatusOK, cName, key, "updated", "")
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
//
//	KEY="123"
//	curl -X DELETE http://localhost:8585/api/journals.ds/object/$KEY
//
// ```
func Delete(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
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
		statusIsOK(w, http.StatusOK, cName, key, "delete", "")
		return
	}
	http.NotFound(w, r)
	return
}

//
// The following routes handle attachments
//

// Attachemnts lists the attachments avialable for a JSON object in the
// collection.
//
// ```shell
//
//	KEY="123"
//	curl -X GET http://localhost:8585/api/journals.ds/attachments/$KEY
//
// ```
func Attachments(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	if len(options) != 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	key := options[0]

	if c, ok := api.CMap[cName]; ok {
		fNames, err := c.Attachments(key)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		src, err := JSONMarshalIndent(fNames, "", "    ")
		if err != nil {
			log.Printf("Attachments, unmarshal error %+v, %s", fNames, err)
			http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
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

// Attach will add or replace an attachment for a JSON object in the
// collection.
//
// ```shell
//
//	KEY="123"
//	FILENAME="mystuff.zip"
//	curl -X POST \
//	   http://localhost:8585/api/journals.ds/attachment/$KEY/$FILENAME
//	     -H "Content-Type: application/zip" \
//	     --data-binary "@./mystuff.zip"
//
// ```
func Attach(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	if len(options) != 2 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	key := options[0]
	fName := options[1]
	c, ok := api.CMap[cName]
	if ok {
		// Handle multipart upload or just streamed upload
		contentType := r.Header.Get("content-type")
		if contentType == `multipart/form-data` || contentType == `multipart/mixed` {
			r.ParseMultipartForm(1024 << 20) // allow up to 1G files
			// Get file handler and name
			file, _, err := r.FormFile("file")
			if err != nil {
				log.Printf("Error in handling file upload %s", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return

			}
			defer file.Close()

			// Now we can attach the file to the record.
			if err := c.AttachStream(key, fName, file); err != nil {
				log.Printf("Failed to attach %q to %q, %s", fName, key, err)
				http.Error(w, http.StatusText(http.StatusInternalServerError)+" "+err.Error(), http.StatusInternalServerError)
				return
			}
			statusIsOK(w, http.StatusCreated, cName, key, "attach", fName)
			return
		} else {
			// Assume raw bytes and read them.
			if err := c.AttachStream(key, fName, r.Body); err != nil {
				log.Printf("Failed to attach stream %q to %q, %s", fName, key, err)
				http.Error(w, http.StatusText(http.StatusInternalServerError)+" "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer r.Body.Close()
			statusIsOK(w, http.StatusCreated, cName, key, "attach", fName)
			return
		}
	}
	// collection or key not found.
	http.NotFound(w, r)
	return
}

// Attach retrieve an attachment from a JSON object in the
// collection.
//
// ```shell
//
//	KEY="123"
//	FILENAME="mystuff.zip"
//	curl -X GET \
//	   http://localhost:8585/api/journals.ds/attachment/$KEY/$FILENAME
//
// ```
func Retrieve(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	if len(options) != 2 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	contentType := "application/octet-stream"
	key, filename := options[0], options[1]
	ext := path.Ext(filename)
	if ext != "" {
		contentType = mime.TypeByExtension(ext)
	}
	c, ok := api.CMap[cName]
	if !ok {
		log.Printf("collection %q not found", cName)
		http.NotFound(w, r)
		return
	}
	_, err := c.AttachmentPath(key, filename)
	if err != nil {
		log.Printf("attachment not found %q from %q in %q", filename, key, cName)
		http.NotFound(w, r)
		return
	}
	if contentType != "" {
		w.Header().Add("Content-Type", contentType)
	}
	err = c.RetrieveStream(key, filename, w)
	if err != nil {
		log.Printf("failed to retrieve stream %q from %q in %q, %s", filename, key, cName, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// Prune removes and attachment from a JSON object in the collection.
//
// ```shell
//
//	KEY="123"
//	FILENAME="mystuff.zip"
//	curl -X DELETE \
//	   http://localhost:8585/api/journals.ds/attachment/$KEY/$FILENAME
//
// ```
func Prune(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	if len(options) != 2 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	key, fName := options[0], options[1]
	c, ok := api.CMap[cName]
	if !ok {
		log.Printf("collection %q not found", cName)
		http.NotFound(w, r)
		return
	}
	err := c.Prune(key, fName)
	if err != nil {
		log.Printf("failed to removed %q from %q in %q", fName, key, cName)
		http.NotFound(w, r)
		return
	}
	statusIsOK(w, http.StatusOK, cName, key, "prune", fName)
}

//
// The following routes handle frames
//

// HasFrame checks a collection for a frame by its name
//
// ```shell
//
//	FRM_NAME="name"
//	curl -X GET http://localhost:8585/api/journals.ds/has-frame/$FRM_NAME
//
// ```
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
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			fmt.Fprint(w, "true")
			return
		}
	}
	// Check if frame is in collection
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// FrameCreate creates a new frame in a collection. It accepts the
// frame definition as a POST of JSON.
//
// ```shell
//
//	FRM_NAME="names"
//	cat<<EOT>frame-def.json
//	{
//	  "dot_paths": [ ".given", ".family" ],
//	  "labels": [ "Given Name", "Family Name" ],
//	  "keys": [ "Miller-A", "Stienbeck-J", "Topez-T", "Valdez-L" ]
//	}
//	EOT
//	curl -X POST http://localhost:8585/api/journals.ds/frame/$FRM_NAME
//	     -H "Content-Type: application/json" \
//	     --data-binary "@./frame-def.json"
//
// ```
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
		log.Printf("FrameCreate, unmarshal error %+v, %s", m, err)
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
		return
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
		if _, err := c.FrameCreate(frameName, keys, dotPaths, labels, verbose); err != nil {
			log.Printf("FrameCreate, Bad Request %s %q %s", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		statusIsOK(w, http.StatusCreated, cName, "", "frame-create", frameName)
		return
	}
	// Check if frame is in collection
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	return
}

// Frames retrieves a list of available frames in a collection.
//
// ```shell
//
//	curl -X GET http://localhost:8585/api/journals.ds/frames
//
// ```
func Frames(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	// Get collection
	c, ok := api.CMap[cName]
	if ok {
		frameNames := c.FrameNames()
		src, err := JSONMarshalIndent(frameNames, "", "    ")
		if err != nil {
			log.Printf("marshal error %+v, %s", frameNames, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(w, "%s", src)
		return
	}
	// Collection not found
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	return
}

// FrameKeys retrieves the list of keys associated with a frame
//
// ```shell
//
//	FRM_NAME="names"
//	curl -X GET http://localhost:8585/api/journals.ds/frame-keys/$FRM_NAME
//
// ```
func FrameKeys(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	if len(options) < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	frameName := options[0]
	// Get collection
	c, ok := api.CMap[cName]
	if ok {
		keys := c.FrameKeys(frameName)
		src, err := JSONMarshalIndent(keys, "", "    ")
		if err != nil {
			log.Printf("marshal error %+v, %s", keys, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(w, "%s", src)
		return
	}
	// Collection not found
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	return
}

// FrameDef retrieves the frame definition associated with a frame
//
// ```shell
//
//	FRM_NAME="names"
//	curl -X GET http://localhost:8585/api/journals.ds/frame-def/$FRM_NAME
//
// ```
func FrameDef(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	if len(options) < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	frameName := options[0]
	// Get collection
	c, ok := api.CMap[cName]
	if ok {
		def, err := c.FrameDef(frameName)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		src, err := JSONMarshalIndent(def, "", "    ")
		if err != nil {
			log.Printf("marshal error %+v, %s", def, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(w, "%s", src)
		return
	}
	// Collection not found
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	return
}

// FrameObjects retrieves the frame objects associated with a frame
//
// ```shell
//
//	FRM_NAME="names"
//	curl -X GET http://localhost:8585/api/journals.ds/frame-objects/$FRM_NAME
//
// ```
func FrameObjects(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	if len(options) < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	frameName := options[0]
	// Get collection
	c, ok := api.CMap[cName]
	if ok {
		objects, err := c.FrameObjects(frameName)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		src, err := JSONMarshalIndent(objects, "", "    ")
		if err != nil {
			log.Printf("marshal error %+v, %s", objects, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(w, "%s", src)
		return
	}
	// Collection not found
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	return
}

// FrameUpdate updates a frame either refreshing the current frame objects
// on the keys associated with the object or if a JSON array of keys is
// provided it reframes the objects using the new list of keys.
//
// ```shell
//
//	FRM_NAME="names"
//	curl -X PUT http://localhost:8585/api/journals.ds/frame/$FRM_NAME
//
// ```
//
// Reframing a frame providing new keys looks something like this --
//
// ```shell
//
//	FRM_NAME="names"
//	cat<<EOT>frame-keys.json
//	[ "Gentle-M", "Stienbeck-J", "Topez-T", "Valdez-L" ]
//	EOT
//	curl -X PUT http://localhost:8585/api/journals.ds/frame/$FRM_NAME \
//	     -H "Content-Type: application/json" \
//	     --data-binary "@./frame-keys.json"
//
// ```
func FrameUpdate(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	if len(options) < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	frameName := options[0]
	// Get collection
	c, ok := api.CMap[cName]
	if ok {
		// Check to see if we have a body containing a list of keys
		body, err := ioutil.ReadAll(r.Body)
		if err == nil && len(body) > 0 {
			// Handle reframe
			keys := []string{}
			if err := json.Unmarshal(body, &keys); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			if err := c.FrameReframe(frameName, keys, verbose); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			statusIsOK(w, http.StatusOK, cName, "", "reframe", frameName)
			return
		}
		err = c.FrameRefresh(frameName, verbose)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		statusIsOK(w, http.StatusOK, cName, "", "refresh", frameName)
		return
	}
	// Collection not found
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	return
}

// FrameDelete removes a frame from a collection.
//
// ```shell
//
//	FRM_NAME="names"
//	curl -X DELETE http://localhost:8585/api/journals.ds/frame/$FRM_NAME
//
// ```
func FrameDelete(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	if len(options) < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	frameName := options[0]
	// Get collection
	c, ok := api.CMap[cName]
	if ok {
		err := c.FrameDelete(frameName)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		statusIsOK(w, http.StatusOK, cName, "", "frame-delete", frameName)
		return
	}
	// Collection not found
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	return
}

//**************************************************************
// The following routes handle JSON object versions. These are
// placeholders for some future release of 2.X.
//**************************************************************

func ObjectVersions(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func ReadVersion(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func DeleteVersion(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

//**************************************************************
// The following routes handle attachment versions. These are
// placeholders for some future release of 2.X.
//**************************************************************

func AttachmentVersions(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func RetrieveVersion(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func PruneVersion(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
