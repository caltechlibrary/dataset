package dataset

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	// timestamp holds the Format of a MySQL time field
	timestamp = "2006-01-02 15:04:05"
	datestamp = "2006-01-02"
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
// End point documentation
//
func readmeDocument() string {
	return fmt.Sprintf(`
Datasetd
========

Overview
--------

__datasetd__ is a minimal web service typically run on localhost port 8485
that exposes a dataset collection as a web service. It features a subset of functionality available with the dataset command line program. __datasetd__ does support multi-process/asynchronous update to a dataset collection. 

__datasetd__ is notable in what it does not provide. It does not provide user/role access restrictions to a collection. It is not intended to be a stand alone web service on the public internet or local area network. It does not provide support for search or complex querying. If you need these features I suggest looking at existing mature NoSQL style solutions like Couchbase, MongoDB, MySQL (which now supports JSON objects) or Postgres (which also support JSON objects). __datasetd__ is a simple, miminal service.

NOTE: You could run __datasetd__ with access control based on a set of set of URL paths by running __datasetd__ behind a full feature web server like Apache 2 or NginX but that is beyond the skope of this project.

Configuration
-------------

__datasetd__ can make one or more dataset collections visible over HTTP/HTTPS. The dataset collections hosted need to be avialable on the same file system as where __datasetd__ is running. __datasetd__ is configured by reading a "settings.json" file in either the local directory where it is launch or by a specified directory on the command line.  

The "settings.json" file has the following structure

    {
        "host": "localhost:8483",
        "collections": {
            "<COLLECTION_ID>": {
                "dataset": "PATH_TO_DATASET_COLLECTION",
                "keys": true,
                "create": true,
                "read": true,
                "update": true,
                "delete": false
            }
        }
    }

In the "collections" object the "<COLLECTION_ID>" is a string which will be used as the start of the path in the URL. The "dataset" attribute sets the path to the dataset collection made available at "<COLLECTION_ID>". For each
collection you can allow the following sub-paths, "create", "read", "update", "delete" and "keys". These sub-paths correspond to their counter parts in the dataset command line tool. In this way would can have a
dataset collection function as a drop box, a read only list or a simple JSON
object storage service.

Running datasetd
----------------

__datasetd__ runs as a HTTP/HTTPS service and as such can be exploit as other network based services can be.  It is recommend you only run __datasetd__ on localhost on a trusted machine. If the machine is a multi-user machine all users can have access to the collections exposed by __datasetd__ regardless of the file permissions they may in their account.
E.g. If all dataset collections are in a directory only allowed access to be the "web-data" user but another user on the system can run cURL then they can access the dataset collections based on the rights of the "web-data" user.  This is a typical situation for most web services and you need to be aware of it if you choose to run __datasetd__.

Supported Features
------------------

__datasetd__ provide a limitted subset of actions support by the standard datset command line tool. It only supports the following verbs

1. keys (return a list of all keys in the collection)
2. create (create a new JSON document in the collection)
3. read (read a JSON document from a collection)
4. update (update a JSON document in the collection)
5. delete (delete a JSON document in the collection)

Each of theses "actions" can be restricted in the configuration (
i.e. "settings.json" file) by setting the value to "false". If the
attribute for the action is not specified in the JSON settings file
then it is assumed to be "false".

Example
-------

E.g. if I have a settings file for "recipes" based on the collection
"recipes.ds" and want to make it read only I would make the attribute
"read" set to true and if I want the option of listing the keys in the collection I would set that true also.

    {
        "host": "localhost:8485",
        "collections": {
            "recipes": {
                "dataset": "recipes.ds",
                "keys": true,
                "read": true
            }
        }
    }

I would start __datasetd__ with the following command line.

    datasetd settings.json

This would display the start up message and log output of the service.

In another shell session I could then use cURL to list the keys and read
a record. In this example I assume that "waffles" is a JSON document
in dataset collection "recipes.ds".

    curl http://localhost:8485/recipies/read/waffles

This would return the "waffles" JSON document or a 404 error if the 
document was not found.

Listing the keys for "recipes.ds" could be done with this cURL command.

    curl http://localhost:8485/recipies/keys

This would return a list of keys, one per line. You could show
all JSON documents in the collection be retrieving a list of keys
and iterating over them using cURL. Here's a simple example in Bash.

    for KEY in $(curl http://localhost:8485/recipes/keys); do
       curl "http://localhost/8485/recipe/read/${KEY}"
    done

Documentation
-------------

__datasetd__ provide documentation as plain text output via request
to the service end points without parameters. Continuing with our
"recipes" example. Try the following URLs with cURL.

    curl http://localhost:8485
    curl http://localhost:8485/collections
    curl http://localhost:8485/recipes
    curl http://localhost:8485/recipes/read

`)
}

func keysDocument(collectionID string) string {
	return fmt.Sprintf(`"/%s/keys" accepts a "GET" returns a list of available keys as JSON array or HTTP error if creation fails.`, collectionID)
}

func createDocument(collectionID string) string {
	return fmt.Sprintf(`"/%s/created/<KEY>" accepts a "POST", creates a JSON document for the <KEY> provided and HTTP 200 OK or HTTP error if creation fails. The "POST" needs to be JSON encoded.`, collectionID)
}

func readDocument(collectionID string) string {
	return fmt.Sprintf(`"/%s/read/<KEY>" accepts a "GET", returns the JSON document for given <KEY> or a HTTP error if not found.`, collectionID)
}

func updateDocument(collectionID string) string {
	return fmt.Sprintf(`"/%s/updated/<KEY>" accepts a "POST", updates the JSON document the <KEY> provided and returns HTTP 200 OK or HTTP error if update fails. The "POST" needs to be JSON encoded.`, collectionID)
}

func deleteDocument(collectionID string) string {
	return fmt.Sprintf(`"/%s/deleted/<KEY>" accepts a "GET", returns HTTP 200 OK on successful deletion or an HTTP error otherwise.`, collectionID)
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
// or '`/<COLLECTION_ID/<END-POINT>/<KEY>`)
//

func keysEndPoint(w http.ResponseWriter, r *http.Request, collectionID string, key string) (int, error) {
	return 501, fmt.Errorf("Not Implemented")
}

func createEndPoint(w http.ResponseWriter, r *http.Request, collectionID string, key string) (int, error) {
	if key == "" {
		return packageDocument(w, createDocument(collectionID))
	}
	return 501, fmt.Errorf("Not Implemented")
}

func readEndPoint(w http.ResponseWriter, r *http.Request, collectionID string, key string) (int, error) {
	if key == "" {
		return packageDocument(w, readDocument(collectionID))
	}
	return 501, fmt.Errorf("Not Implemented")
}

func updateEndPoint(w http.ResponseWriter, r *http.Request, collectionID string, key string) (int, error) {
	if key == "" {
		return packageDocument(w, updateDocument(collectionID))
	}
	return 501, fmt.Errorf("Not Implemented")
}

func deleteEndPoint(w http.ResponseWriter, r *http.Request, collectionID string, key string) (int, error) {
	if key == "" {
		return packageDocument(w, deleteDocument(collectionID))
	}
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
		collectionID, endPoint, key = args[0], args[1], key
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

func RunDatasetAPI(appName string) error {
	/* Setup web server */
	log.Printf(`
%s %s

Listening on http://%s

Press ctl-c to terminate.
`, appName, Version, config.Hostname)
	http.HandleFunc("/", api)
	return http.ListenAndServe(config.Hostname, nil)
}
