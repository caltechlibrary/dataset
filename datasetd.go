package dataset

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

const (
	// timestamp holds the Format of a MySQL time field
	timestamp = "2006-01-02 15:04:05"
	datestamp = "2006-01-02"

	// jsonSizeLimit is the maximum size JSON object we'll accept via
	// our service. Current 1 MB (2^20)
	jsonSizeLimit = 1048576

	// attachmentSizeLimit is the maximum size of Attachments we'll
	// accept via our service. Current 250 MiB
	attachmentSizeLimit = (jsonSizeLimit * 250)
)

var (
	config *Config
)

//
// response handlers
//

func packageDocument(w http.ResponseWriter, src string) (int, error) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, src)
	return 200, nil
}

func packageJSON(w http.ResponseWriter, collectionID string, src []byte, err error) (int, error) {
	if err != nil {
		log.Printf("ERROR: (%s) package JSON error, %s", collectionID, err)
		return 500, fmt.Errorf("Internal Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", src)
	return 200, nil
}

//
// Expose the collections available
//

func collectionsEndPoint(w http.ResponseWriter, r *http.Request) (int, error) {
	collections := []string{}
	for collectionID, _ := range config.Collections {
		collections = append(collections, collectionID)
	}
	src, err := json.MarshalIndent(collections, "", "  ")
	if err != nil {
		return 500, fmt.Errorf("Internal Server Error, %s", err)
	}
	return packageJSON(w, "", src, err)
}

//
// End Point handlers (route as defined `/<COLLECTION_ID>/<END-POINT>`,
// or `/<COLLECTION_ID/<END-POINT>/<KEY>`)
//

func keysEndPoint(w http.ResponseWriter, r *http.Request, collectionID string, key string) (int, error) {
	contentType := r.Header.Get("Content-Type")
	if r.Method != "GET" {
		return 405, fmt.Errorf(`Method Not Allowed
%s %s`, r.Method, contentType)
	}
	_, ok := config.Collections[collectionID]
	if ok == false {
		return 400, fmt.Errorf("Bad Request")
	}
	ds := config.Collections[collectionID].DS
	if ds == nil {
		return 500, fmt.Errorf("Internal Server Error")
	}
	keys := ds.Keys()
	src, err := json.MarshalIndent(keys, "", "    ")
	if err != nil {
		return 500, fmt.Errorf("Internal Server Error")
	}
	return packageJSON(w, collectionID, src, err)
}

func createEndPoint(w http.ResponseWriter, r *http.Request, collectionID string, key string) (int, error) {
	if key == "" {
		return packageDocument(w, createDocument(collectionID))
	}
	contentType := r.Header.Get("Content-Type")
	if r.Method != "POST" {
		return 405, fmt.Errorf(`Method Not Allowed
%s %s`, r.Method, contentType)
	}
	if contentType != "application/json" {
		return 415, fmt.Errorf(`Unsupported Media Type
%s %s`, r.Method, contentType)
	}
	_, ok := config.Collections[collectionID]
	if ok == false {
		return 400, fmt.Errorf(`Bad Request
%s %s`, r.Method, contentType)
	}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, jsonSizeLimit))
	if err != nil {
		return 400, fmt.Errorf(`Bad Request
cannot read request body for %s

%s
`, key, err)
	}
	data := map[string]interface{}{}
	if err := json.Unmarshal(body, &data); err != nil {
		return 400, fmt.Errorf(`Bad Request
Invalid JSON Object %s

%s
`, key, err)
	}
	ds := config.Collections[collectionID].DS
	if err := ds.Create(key, data); err != nil {
		return 507, fmt.Errorf(`Insufficient Storage
cannot create %s

%s
`, key, err)
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "OK, created %s", key)
	return 201, nil
}

func readEndPoint(w http.ResponseWriter, r *http.Request, collectionID string, key string) (int, error) {
	if key == "" {
		return packageDocument(w, readDocument(collectionID))
	}
	contentType := r.Header.Get("Content-Type")
	if r.Method != "GET" {
		return 405, fmt.Errorf(`Method Not Allowed
%s %s`, r.Method, contentType)
	}
	_, ok := config.Collections[collectionID]
	if ok == false {
		return 400, fmt.Errorf("Bad Request")
	}
	ds := config.Collections[collectionID].DS
	if ds == nil {
		return 500, fmt.Errorf("Internal Server Error")
	}
	data := map[string]interface{}{}
	if err := ds.Read(key, data, false); err != nil {
		return 404, fmt.Errorf(`Not Found
%s
`, err)
	}
	src, err := json.MarshalIndent(data, "", "   ")
	return packageJSON(w, collectionID, src, err)
}

func updateEndPoint(w http.ResponseWriter, r *http.Request, collectionID string, key string) (int, error) {
	if key == "" {
		return packageDocument(w, createDocument(collectionID))
	}
	contentType := r.Header.Get("Content-Type")
	if r.Method != "POST" {
		return 405, fmt.Errorf(`Method Not Allowed
%s %s
`, r.Method, contentType)
	}
	if contentType != "application/json" {
		return 415, fmt.Errorf(`Unsupported Media Type
%s %s
`, r.Method, contentType)
	}
	_, ok := config.Collections[collectionID]
	if ok == false {
		return 400, fmt.Errorf(`Bad Request
%s %s
`, r.Method, contentType)
	}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, jsonSizeLimit))
	if err != nil {
		return 400, fmt.Errorf(`Bad Request
cannot read request body for %s

%s
`, key, err)
	}
	data := map[string]interface{}{}
	if err := json.Unmarshal(body, &data); err != nil {
		return 400, fmt.Errorf(`Bad Request
Invalid JSON Object %s

%s
`, key, err)
	}
	ds := config.Collections[collectionID].DS
	if err := ds.Update(key, data); err != nil {
		return 507, fmt.Errorf(`Insufficient Storage
cannot update %s

%s
`, key, err)
	}
	return packageDocument(w, fmt.Sprintf("OK, updated %s", key))
}

func deleteEndPoint(w http.ResponseWriter, r *http.Request, collectionID string, key string) (int, error) {
	if key == "" {
		return packageDocument(w, createDocument(collectionID))
	}
	contentType := r.Header.Get("Content-Type")
	if r.Method != "GET" {
		return 405, fmt.Errorf(`Method Not Allowed
%s %s
`, r.Method, contentType)
	}
	_, ok := config.Collections[collectionID]
	if (r.Method != "GET") || (ok == false) {
		return 400, fmt.Errorf("Bad Request")
	}
	ds := config.Collections[collectionID].DS
	if err := ds.Delete(key); err != nil {
		return 500, fmt.Errorf(`Internal Server Error
cannot delete %s

%s
`, key, err)
	}
	return packageDocument(w, fmt.Sprintf("OK, deleted %s", key))
}

func attachEndPoint(w http.ResponseWriter, r *http.Request, collectionID string, key string) (int, error) {
	if key == "" {
		return packageDocument(w, attachDocument(collectionID))
	}
	contentType := r.Header.Get("Content-Type")
	if r.Method != "POST" {
		return 405, fmt.Errorf(`Method Not Allowed
%s %s
`, r.Method, contentType)
	}
	if contentType != `multipart/form-data` {
		return 415, fmt.Errorf(`Unsupported Media Type
%s %s
`, r.Method, contentType)
	}
	_, ok := config.Collections[collectionID]
	if ok == false {
		return 400, fmt.Errorf(`Bad Request
%s %s
`, r.Method, contentType)
	}
	log.Printf("attachEndPoint() not implemented")
	return 501, fmt.Errorf("Not Implemented")
}

func attachmentsEndPoint(w http.ResponseWriter, r *http.Request, collectionID string, key string) (int, error) {
	log.Printf("attachmentsEndPoint() not implemented")
	return 501, fmt.Errorf("Not Implemented")
}

func retrieveEndPoint(w http.ResponseWriter, r *http.Request, collectionID string, key string) (int, error) {
	log.Printf("retrieveEndPoint() not implemented")
	return 501, fmt.Errorf("Not Implemented")
}

func pruneEndPoint(w http.ResponseWriter, r *http.Request, collectionID string, key string) (int, error) {
	log.Printf("pruneEndPoint() not implemented")
	return 501, fmt.Errorf("Not Implemented")
}

//
// The following define the API as a service handling errors,
// routes and logging.
//

func logRequest(r *http.Request, status int, err error) {
	q := r.URL.Query()
	errStr := "OK"
	if err != nil {
		errStr = fmt.Sprintf("%s", err)
	}
	if len(q) > 0 {
		log.Printf("%s %s RemoteAddr: %s UserAgent: %s Query: %+v Response: %d %s", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), q, status, errStr)
	} else {
		log.Printf("%s %s RemoteAddr: %s UserAgent: %s Response: %d %s", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), status, errStr)
	}
}

func handleError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, `ERROR: %d %s`, statusCode, err)
}

func routeEndPoints(w http.ResponseWriter, r *http.Request) (int, error) {
	var ()
	parts := strings.Split(r.URL.Path, "/")
	// Create args by removing empty strings from path parts
	args := []string{}
	for _, arg := range parts {
		if arg != "" {
			// FIXME: URL decode that path part
			args = append(args, arg)
		}
	}
	if len(args) == 0 {
		return packageDocument(w, readmeDocument())
	}
	if len(args) == 1 {
		return packageDocument(w, strings.ReplaceAll(readmeDocument(), "<COLLECTION_ID>", args[0]))
	}
	// Expected URL structure of `/<COLLECTION_ID>/<END-POINT>/<KEY>`
	collectionID, endPoint, key := "", "", ""
	if len(args) == 2 {
		collectionID, endPoint, key = args[0], args[1], ""
	} else {
		collectionID, endPoint, key = args[0], args[1], args[2]
	}
	if routes, hasCollection := config.Routes[collectionID]; hasCollection == true {
		// Confirm we have a route
		if fn, hasRoute := routes[endPoint]; hasRoute == true {
			return fn(w, r, collectionID, key)
		}
	}
	return 404, fmt.Errorf("Not Found")
}

func api(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		statusCode int
	)
	// FIXME: the API should reject requests that are not application/json or text/plain since that is all we provide.
	if (r.Method != "GET") && (r.Method != "POST") {
		statusCode, err = 405, fmt.Errorf("Method Not Allowed")
		handleError(w, statusCode, err)
	} else {
		switch r.URL.Path {
		case "/favicon.ico":
			statusCode, err = 200, nil
			fmt.Fprintf(w, "")
			//statusCode, err = 404, fmt.Errorf("Not Found")
			//handleError(w, statusCode, err)
		case "/collections":
			statusCode, err = collectionsEndPoint(w, r)
			if err != nil {
				handleError(w, statusCode, err)
			}
		default:
			statusCode, err = routeEndPoints(w, r)
			if err != nil {
				handleError(w, statusCode, err)
			}
		}
	}
	logRequest(r, statusCode, err)
}

func InitDatasetAPI(settings string) error {
	var err error

	//NOTE: loading the settings into the global config object.
	config, err = LoadConfig(settings)
	if err != nil {
		return err
	}
	if config == nil {
		return fmt.Errorf("Failed to generate a valid configuration")
	}
	if config.Routes == nil {
		config.Routes = map[string]map[string]func(http.ResponseWriter, *http.Request, string, string) (int, error){}
	}
	// Now setup the routes for each collection.
	for collectionID, cfg := range config.Collections {
		// Initialize the map.
		if config.Routes[collectionID] == nil {
			config.Routes[collectionID] = map[string]func(http.ResponseWriter, *http.Request, string, string) (int, error){}
		}
		// Add routes (end points) for the target repository
		if cfg.Keys {
			config.Routes[collectionID]["keys"] = keysEndPoint
		}
		if cfg.Create {
			config.Routes[collectionID]["create"] = createEndPoint
		}
		if cfg.Read {
			config.Routes[collectionID]["read"] = readEndPoint
		}
		if cfg.Update {
			config.Routes[collectionID]["update"] = updateEndPoint
		}
		if cfg.Delete {
			config.Routes[collectionID]["delete"] = deleteEndPoint
		}
	}
	return nil
}

//FIXME: Need to handle a reload/restart for SIGHUP

func Shutdown(appName string) int {
	exitCode := 0
	pid := os.Getpid()
	for cName, c := range config.Collections {
		if c.DS != nil {
			log.Printf("Closing %s", cName)
			c.DS.Close()
		} else {
			log.Printf("Lost connection to %s", cName)
			exitCode = 1
		}
	}
	log.Printf(`Shutdown down %s pid: %d exit code: %d `, appName, pid, exitCode)
	return exitCode
}

func RunDatasetAPI(appName string) error {
	/* Setup web server */
	log.Printf(`
%s %s

Listening on http://%s

Press ctl-c to terminate.
`, appName, Version, config.Hostname)
	processControl := make(chan os.Signal, 1)
	signal.Notify(processControl, os.Interrupt)
	go func() {
		<-processControl
		os.Exit(Shutdown(appName))
	}()
	for cName, c := range config.Collections {
		log.Printf("Opening collection %s", cName)
		if c.DS == nil {
			ds, err := openCollection(c.CName)
			if err != nil {
				return err
			}
			c.DS = ds
		}
	}
	http.HandleFunc("/", api)
	return http.ListenAndServe(config.Hostname, nil)
}
