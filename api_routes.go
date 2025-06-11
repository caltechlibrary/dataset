package dataset

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"

	// Caltech Library packages
	"github.com/caltechlibrary/models"

	// 3rd Party packages
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
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

// statusOKText returns a text indicating the status of a request
// is OK.
func statusIsOKText(w http.ResponseWriter, r *http.Request, statusCode int, cName string, key string, action string, target string, successRedirect string) {
	log.Printf("Success request URL: %s %q  statusCode: %d, cName: %q, key: %q, action %q, target: %q, successRedirect: %q", r.Method, r.URL, statusCode, cName, key, action, target, successRedirect)
	if successRedirect != "" {
		// Redirecting the to the success page, using HTTP status found
		http.Redirect(w, r, successRedirect, http.StatusFound)
		return
	}
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%s for %s\n", http.StatusText(statusCode), cName)
	if key != "" {
		fmt.Fprintf(w, "  %s\n", key)
	}
	if action != "" {
		fmt.Fprintf(w, "  %s\n", action)
	}
	if target != "" {
		fmt.Fprintf(w, "  %s\n", target)
	}
}

// statusIsError provides a consistent way to express an error condition in the API request.
// It is not intended to be used where `http.NotFound(w, r)` is more appropriate. It was created
// as a means of handling redirects specified in the YAML for for datasetd.
func statusIsError(w http.ResponseWriter, r *http.Request, statusText string, statusCode int, errorRedirect string) {
	// Do we have a redirect?
	if errorRedirect != "" {
		// Redirecting using http status NotModified
		log.Printf("ERROR request URL: %s %q  statusText: %q, statusCode: %d, errorRedirect: %q", r.Method, r.URL, statusText, statusCode, errorRedirect)
		// Redirecting the to the error page.
		http.Redirect(w, r, errorRedirect, http.StatusNotModified)
		return
	}
	// Fallback to the default error handler
	http.Error(w, statusText, statusCode)
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
	fmt.Fprintf(w, "%s %s", filepath.Base(api.AppName), api.Version)
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
			statusIsError(w, r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, "")
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
			http.NotFound(w, r)
			return
		}
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "%s", src)
		return
	}
	http.NotFound(w, r)
	return
}

// getAttrNames take a request URL and parameter value as attribute names (column headings in CSV).
// If not found then it returns a empty list of attribute names.
func getAttrNames(q url.Values, key string) []string {
	if q != nil {
		if q != nil {
			s := q.Get(key)
			if s != "" {
				return strings.Split(s, ",")
			}
		}
	}
	return []string{}
}

// Query returns the results from a SQL function stored in MySQL or Postgres.
// The query takes a query name followed by a path part that maps the order of
// the fields. This is needed because the SQL prepared statments use paramter
// order is mostly common to SQL dialects.
//
// In this example we're runing the SQL statement with the name of "journal_search"
// with title mapped to `$1` and journal mapped to `$2`.
//
// ```shell
//
//	curl -X POST http://localhost:8485/api/journals.ds/query/journal_search/title/journal \
//	     --data "title=Princess+Bride" \
//	     --data "journal=Movies+and+Popculture"
//
// ```
//
// NOTE: the SQL query must conform to the same constraints as dsquery SQL constraints.
func Query(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
	// NOTE: Need to determine content type requested via the content type heeader
	//
	// CSV -> text/csv
	// YAML -> application/yaml
	// JSON -> application/json
	//
	// if contentType is "" then it will default to application/json
	contentType := "application/json"
	if r.Header != nil {
		contentType = r.Header.Get("content-type")
	}
	// Make content type align to "csv" or "yaml" query parameter when passed.
	urlQuery := r.URL.Query()
	if urlQuery.Get("csv") != "" {
		contentType = "text/csv"
	}
	if urlQuery.Get("yaml") != "" {
		contentType = "application/yaml"
	}
	if api.Debug {
		log.Printf("DEBUG Query got a query, cName: %q, verb: %q, content type: %q, options: %+v\n", cName, verb, contentType, options)
	}
	if len(options) == 0 {
		log.Printf("Query, Bad Request %s %q, missing query name", r.Method, r.URL.Path)
		statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
		return
	}
	qName := options[0]
	cfg, err := api.Settings.GetCfg(cName)
	if err != nil {
		log.Printf("Query, Bad Request %s %q %s, not found", r.Method, r.URL.Path, qName)
		statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
		return
	}
	if api.Debug {
		log.Printf("DEBUG cfg attribute: c -> %+v\n", cfg)
	}
	// Make sure we have queries defined in the configuration
	if cfg.QueryFn == nil || len(cfg.QueryFn) == 0 {
		log.Printf("Query, Bad Request %s %q, undefined query", r.Method, r.URL.Path)
		statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
		return
	}
	qStmt, ok := cfg.QueryFn[qName]
	if !ok {
		log.Printf("Query, Bad Request %s %q, undefined query %q", r.Method, r.URL.Path, qName)
		statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
		return
	}
	if api.Debug {
		log.Printf("DEBUG qStmt: %q\n", qStmt)
	}

	if c, ok := api.CMap[cName]; ok {
		if api.Debug {
			log.Printf("DEBUG c : c -> %+v\n", c)
		}
		src, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Query, Bad Request %s %q %s", r.Method, r.URL.Path, err)
			statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
			return
		}
		defer r.Body.Close()
		o := map[string]interface{}{}
		if len(src) > 0 {
			err = json.Unmarshal(src, &o)
			if err != nil {
				log.Printf("Query, unmarshal error %+v, %s", o, err)
				statusIsError(w, r, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable, "")
				return
			}
		}

		if api.Debug {
			log.Printf("DEBUG verb %q, options %+v\n", verb, options)
		}
		// FIXME: how to map form names to a list of parameters?
		var rows *sql.Rows
		qParams := []interface{}{}
		if len(options) > 0 && len(o) > 0 {
			for _, key := range options {
				if val, ok := o[key]; ok {
					qParams = append(qParams, val)
				} else if api.Debug {
					log.Printf("DEBUG option %q, failed to map to query object %+v\n", key, o)
				}
			}
			rows, err = c.SQLStore.db.Query(qStmt, qParams...)
		} else {
			if api.Debug {
				log.Printf("DEBUG qStmt %+v", qStmt)
			}
			rows, err = c.SQLStore.db.Query(qStmt)
		}
		if err != nil {
			log.Printf("Query, failed stmt: %q, %s, %+v, %s", qName, qStmt, qParams, err)
			statusIsError(w, r, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable, "")
			return
		}
		src = []byte(`[`)
		i := 0
		for rows.Next() {
			// Get our row values
			obj := []byte{}
			if err := rows.Scan(&obj); err != nil {
				log.Printf("Query, failed row.Scan(obj): %q, %s, (%d) %+v, %s", qName, qStmt, i, obj, err)
				statusIsError(w, r, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable, "")
				return
			}
			if i > 0 {
				src = append(src, ',')
			}
			src = append(src, obj...)
			i++
		}
		src = append(src, ']')
		err = rows.Err()
		if err != nil {
			log.Printf("Query, row.Err(): %q, %s, (%d) %+v, %s", qName, qStmt, i, src, err)
			statusIsError(w, r, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable, "")
			return
		}
		// NOTE: I need to handle the content type requested -> CSV, YAML or JSON (default)
		switch contentType {
		case "text/csv":
			src, err = MakeCSV(src, getAttrNames(urlQuery, "csv"))
			if err != nil {
				log.Printf("Failed to convert %q to %q, %s", r.URL.Path, contentType, err)
				statusIsError(w, r, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable, "")
				return
			}
		case "application/yaml":
			src, err = MakeCSV(src, getAttrNames(urlQuery, "yaml"))
			if err != nil {
				log.Printf("Failed to convert %q to %q, %s", r.URL.Path, contentType, err)
				statusIsError(w, r, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable, "")
				return
			}
		}
		w.Header().Add("Content-Type", contentType)
		fmt.Fprintf(w, "%s", src)
		return
	}
	http.NotFound(w, r)
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
			statusIsError(w, r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, "")
			return
		}
		// Set header to application/json
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", src)
		return
	}
	http.NotFound(w, r)
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
	models.SetDebug(api.Debug)
	defer r.Body.Close()
	var key string
	if len(options) > 0 {
		key = options[0]
		if api.Debug {
			log.Printf("DEBUG using key from URL.path %q", key)
		}
	}
	//NOTE: Handle the cases for submissions encoded as JSON or urlencoded data.
	contentType := r.Header.Get("content-type")
	idName := ""
	successRedirect := ""
	errorRedirect := ""
	if (api.Settings != nil) {
		if cfg, err := api.Settings.GetCfg(cName); err == nil {
			successRedirect, errorRedirect = cfg.CreateSuccess, cfg.CreateError
			if api.Debug {
				log.Printf("DEBUG successRedirect: %q, errorRedirect: %q", successRedirect, errorRedirect)
			}
		}
	}
	o := map[string]interface{}{}
	if c, ok := api.CMap[cName]; ok {
		if api.Debug {
			log.Printf("DEBUG contentType -> %q", contentType)
		}
		switch contentType {
		case "application/json":
			src, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Printf("Create, Bad Request %s %q %s", r.Method, r.URL.Path, err)
				statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, errorRedirect)
				return
			}
			err = json.Unmarshal(src, &o)
			if err != nil {
				log.Printf("Create, json unmarshal error %+v, %s", o, err)
				statusIsError(w, r, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable, errorRedirect)
				return
			}
		default:
			//NOTE: Need to know the form field names, this is in .Model
			r.ParseForm()
			if c.Model == nil {
				if api.Debug {
					log.Printf("DEBUG c.Model is nil, accept all fields without validation")
				}
				for key, _ := range r.Form {
					o[key] = r.Form.Get(key)
				}
			} else {
				// NOTE: We only want to grab the fields defined in the model!!
				idName = c.Model.GetPrimaryId()
				if api.Debug {
					log.Printf("DEBUG c.Model is populated, validating model's fields")
				}
				for _, key := range c.Model.GetElementIds() {
					//FIXME: Need to validate the field value
					val := r.Form.Get(key)
					o[key] = val
				}
			}
			if api.Debug {
				log.Printf("DEBUG creating form object -> %+v", o)
			}
		}
		if key == "" {
			if idName == "" {
				log.Printf("Missing primary id, bad request %s %q, model %+v, data %+v", r.Method, r.URL.Path, c.Model, o)
				statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, errorRedirect)
				return
			}
			if id, ok := o[idName]; ok {
				if api.Debug {
					log.Printf("DEBUG key set to record 'id' value, %+v", id)
				}
				key = id.(string)
			}
		}
		// NOTE: We need to handle case that browsers don't support "PUT" method without JavaScript
		if c.HasKey(key) {
			if c.Model != nil {
				// NOTE: Handle generated types on update (e.g. don't overwrite one_time_timestamp or uuid)
				generatedTypes := c.Model.GetGeneratedTypes()
				if len(generatedTypes) > 0 {
					// Get existing record so we can handle generated field updates properly.
					oldO := map[string]interface{}{}
					if err := c.Read(key, oldO); err != nil {
						log.Printf("Failed to retrieve record before update failed %+v, %s", o, err)
						statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, errorRedirect)
						return
					}
					for k, genType := range generatedTypes {
						// Handle write once types, created_timestamp, uuid
						if _, ok := o[k]; ok {
							switch genType {
							case "uuid":
								// Only generate the UUID if field is missing.
								if val, ok := oldO[k]; (ok == false) || (val == "") {
									uid, err := uuid.NewV7()
									if err == nil {
										o[k] = uid.String()
									}
								} else {
									// Preserve the old UUID assigned.
									o[k] = oldO[k]
								}
								// If idName is found then force key to UUID set in o map
								if k == idName {
									key = o[k].(string)
								}
							case "created_timestamp":
								if val, ok := oldO[k]; (ok == false) || (val == "") {
									o[k] = time.Now().Format(time.RFC3339)
								}
							}
						}
						// Handle overwriting types (i.e. timestamp, current_timestamp)
						switch genType {
						case "timestamp":
							o[k] = time.Now().Format(time.RFC3339)
						case "current_timestamp":
							o[k] = time.Now().Format(time.RFC3339)
						}
					}
				}
				// Now we need to validate the form data against our model.
				if ok := c.Model.ValidateMapInterface(o); !ok {
					log.Printf("Failed to validate create form, bad request %s %q -> %+v", r.Method, r.URL.Path, o)
					statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, errorRedirect)
					return
				}
			}
			if api.Debug {
				txt, _ := json.MarshalIndent(o, "", "  ")
				log.Printf("DEBUGING form data:\n%s\n\n", txt)
			}
			// Now if we have formData populated it needs to get validated after generated types appliced
			if err := c.Update(key, o); err != nil {
				log.Printf("Update failed %+v, %s", o, err)
				statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, errorRedirect)
				return
			}
			//FIXME: If urlencoded data then redirect for the form to some URL. What is the way to indicate this?
			// Early web stuff used a hidden form field for this, seems really clunky.
			if contentType == "application/json" {
				statusIsOK(w, http.StatusOK, cName, key, "updated", "")
			} else {
				statusIsOKText(w, r, http.StatusOK, cName, key, "updated", "", fmt.Sprintf("%s?key=%s", successRedirect, key))
			}
			return
		} else {
			if c.Model != nil {
				// NOTE: Handle generated types on create
				for k, genType := range c.Model.GetGeneratedTypes() {
					switch genType {
					case "uuid":
						uid, err := uuid.NewV7()
						if err == nil {
							o[k] = uid.String()
							if k == idName {
								key = uid.String()
							}
						}
					case "timestamp":
						o[k] = time.Now().Format(time.RFC3339)
					case "current_timestamp":
						o[k] = time.Now().Format(time.RFC3339)
					case "created_timestamp":
						o[k] = time.Now().Format(time.RFC3339)
					}
				}
			}
			if api.Debug {
				txt, _ := json.MarshalIndent(o, "", "  ")
				log.Printf("DEBUG form data:\n%s\n\n", txt)
			}
			// Now we need to validate the form data.
			if c.Model != nil {
				if ok := c.Model.ValidateMapInterface(o); !ok {
					log.Printf("Failed to validate create form, bad request %s %q -> %+v", r.Method, r.URL.Path, o)
					statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, errorRedirect)
					return
				}
			}
			if err := c.Create(key, o); err != nil {
				log.Printf("Create failed %+v, %s", o, err)
				statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, errorRedirect)
				return
			}
			//FIXME: If urlencoded data then redirect for the form (is this in the referrer URL?)
			if contentType == "application/json" {
				statusIsOK(w, http.StatusCreated, cName, key, "created", "")
			} else {
				statusIsOKText(w, r, http.StatusOK, cName, key, "updated", "", fmt.Sprintf("%s?key=%s", successRedirect, key))
			}
			return
		}
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
		statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
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
		//FIXME: Is the request for JSON or YAML data?
		var src []byte
		contentType := r.Header.Get("content-type")
		if api.Debug {
			log.Printf("DEBUG contenType: %q", contentType)
		}
		switch contentType {
		case "application/yaml":
			src, err = yaml.Marshal(o)
			if err != nil {
				log.Printf("Read, yaml marshal error %+v, %s", o, err)
				statusIsError(w, r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, "")
				return
			}
		case "application/json":
			src, err = JSONMarshalIndent(o, "", "    ")
			if err != nil {
				log.Printf("Read, json marshal error %+v, %s", o, err)
				statusIsError(w, r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, "")
				return
			}
		default:
			src, err = JSONMarshalIndent(o, "", "    ")
			if err != nil {
				log.Printf("Read, json marshal error %+v, %s", o, err)
				statusIsError(w, r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, "")
				return
			}
		}
		// Set header to application/json
		w.Header().Add("Content-Type", contentType)
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
		statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
		return
	}
	key := options[0]

	idName := ""
	formData := map[string]string{}
	if c, ok := api.CMap[cName]; ok {
		o := map[string]interface{}{}
		//FIXME: Are we being sent JSON, YAML or urlencoded data?
		contentType := r.Header.Get("content-type")
		if api.Debug {
			log.Printf("DEBUG contentType: %q", contentType)
		}
		switch contentType {
		case "application/json":
			src, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Printf("Update, Bad Request %s %q %s", r.Method, r.URL.Path, err)
				statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
				return
			}
			err = json.Unmarshal(src, &o)
			if err != nil {
				log.Printf("Update, json unmarshal error %+v, %s", o, err)
				statusIsError(w, r, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable, "")
				return
			}
		default:
			//NOTE: Need to know the form field names, this is in .Model
			r.ParseForm()
			idName = c.Model.GetPrimaryId()
			if c.Model == nil {
				if api.Debug {
					log.Printf("DEBUG c.Model is nil, accept all fields without validation")
				}
				for key, _ := range r.Form {
					o[key] = r.Form.Get(key)
				}
			} else {
				// NOTE: We only want to grab the fields defined in the model!!
				if api.Debug {
					log.Printf("DEBUG c.Model is populated, validating model's fields")
				}
				for _, key := range c.Model.GetElementIds() {
					//FIXME: Need to validate the field value
					val := r.Form.Get(key)
					o[key] = val
					formData[key] = val
				}
				if api.Debug {
					txt, _ := json.MarshalIndent(formData, "", "  ")
					log.Printf("DEBUG form data:\n%s\n\n", txt)
				}
				// Now we need to validate the form data.
				if ok := c.Model.Validate(formData); !ok {
					log.Printf("Failed to validate create form, bad request %s %q -> %+v", r.Method, r.URL.Path, formData)
					//statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
					//return
				}
			}
			if api.Debug {
				log.Printf("DEBUG creating form object -> %+v", o)
			}
		}
		if key == "" {
			if idName == "" {
				log.Printf("Missing primary id, bad request %s %q, model %+v, data %+v", r.Method, r.URL.Path, c.Model, formData)
				statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
				return
			}
			if id, ok := o[idName]; ok {
				if api.Debug {
					log.Printf("DEBUG key set to record 'id' value, %+v", id)
				}
				key = id.(string)
			}
		}
		if c.Model != nil {
			generatedTypes := c.Model.GetGeneratedTypes()
			// NOTE: Handle generated types on update (e.g. don't overwrite one_time_timestamp or uuid)
			if len(generatedTypes) > 0 {
				// Get existing record so we can handle generated field updates properly.
				oldO := map[string]interface{}{}
				if err := c.Read(key, oldO); err != nil {
					log.Printf("Failed to retrieve record before update failed %+v, %s", o, err)
					statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
					return
				}
				for k, genType := range generatedTypes {
					// Handle write once types, created_timestamp, uuid
					if _, ok := o[k]; ok {
						switch genType {
						case "uuid":
							// Only generate the UUID if field is missing.
							if val, ok := oldO[k]; (ok == false) || (val == "") {
								uid, err := uuid.NewV7()
								if err == nil {
									o[k] = uid.String()
								}
							} else {
								// Preserve the old UUID assigned.
								o[k] = oldO[k]
							}
							// If idName is found then force key to UUID set in o map
							if k == idName {
								key = o[k].(string)
							}
						case "created_timestamp":
							if val, ok := oldO[k]; (ok == false) || (val == "") {
								o[k] = time.Now().Format(time.RFC3339)
							}
						}
					}
					// Handle overwriting types (i.e. timestamp, current_timestamp)
					switch genType {
					case "timestamp":
						o[k] = time.Now().Format(time.RFC3339)
					case "current_timestamp":
						o[k] = time.Now().Format(time.RFC3339)
					}
				}
			}
		}
		if err := c.Update(key, o); err != nil {
			statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
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
		statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
		return
	}
	key := options[0]

	if c, ok := api.CMap[cName]; ok {
		if err := c.Delete(key); err != nil {
			statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
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
		statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
		return
	}
	key := options[0]

	if c, ok := api.CMap[cName]; ok {
		fNames, err := c.Attachments(key)
		if err != nil {
			statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
			return
		}
		src, err := JSONMarshalIndent(fNames, "", "    ")
		if err != nil {
			log.Printf("Attachments, unmarshal error %+v, %s", fNames, err)
			statusIsError(w, r, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable, "")
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
		statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
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
				statusIsError(w, r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, "")
				return

			}
			defer file.Close()

			// Now we can attach the file to the record.
			if err := c.AttachStream(key, fName, file); err != nil {
				log.Printf("Failed to attach %q to %q, %s", fName, key, err)
				statusIsError(w, r, http.StatusText(http.StatusInternalServerError)+" "+err.Error(), http.StatusInternalServerError, "")
				return
			}
			statusIsOK(w, http.StatusCreated, cName, key, "attach", fName)
			return
		} else {
			// Assume raw bytes and read them.
			if err := c.AttachStream(key, fName, r.Body); err != nil {
				log.Printf("Failed to attach stream %q to %q, %s", fName, key, err)
				statusIsError(w, r, http.StatusText(http.StatusInternalServerError)+" "+err.Error(), http.StatusInternalServerError, "")
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
		statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
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
		statusIsError(w, r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, "")
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
		statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
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
	http.NotFound(w, r)
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
		statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
		return
	}
	m := map[string][]string{}
	if err := json.Unmarshal(src, &m); err != nil {
		log.Printf("FrameCreate, unmarshal error %+v, %s", m, err)
		statusIsError(w, r, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable, "")
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
			statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
		}
		statusIsOK(w, http.StatusCreated, cName, "", "frame-create", frameName)
		return
	}
	// Check if frame is in collection
	http.NotFound(w, r)
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
			statusIsError(w, r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, "")
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(w, "%s", src)
		return
	}
	// Collection not found
	http.NotFound(w, r)
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
		statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
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
			statusIsError(w, r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, "")
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(w, "%s", src)
		return
	}
	// Collection not found
	http.NotFound(w, r)
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
		statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
		return
	}
	frameName := options[0]
	// Get collection
	c, ok := api.CMap[cName]
	if ok {
		def, err := c.FrameDef(frameName)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		src, err := JSONMarshalIndent(def, "", "    ")
		if err != nil {
			log.Printf("marshal error %+v, %s", def, err)
			statusIsError(w, r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, "")
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(w, "%s", src)
		return
	}
	// Collection not found
	http.NotFound(w, r)
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
		statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
		return
	}
	frameName := options[0]
	// Get collection
	c, ok := api.CMap[cName]
	if ok {
		objects, err := c.FrameObjects(frameName)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		src, err := JSONMarshalIndent(objects, "", "    ")
		if err != nil {
			log.Printf("marshal error %+v, %s", objects, err)
			statusIsError(w, r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, "")
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(w, "%s", src)
		return
	}
	// Collection not found
	http.NotFound(w, r)
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
		statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
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
				statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
				return
			}
			if err := c.FrameReframe(frameName, keys, verbose); err != nil {
				statusIsError(w, r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, "")
				return
			}
			statusIsOK(w, http.StatusOK, cName, "", "reframe", frameName)
			return
		}
		err = c.FrameRefresh(frameName, verbose)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		statusIsOK(w, http.StatusOK, cName, "", "refresh", frameName)
		return
	}
	// Collection not found
	http.NotFound(w, r)
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
		statusIsError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, "")
		return
	}
	frameName := options[0]
	// Get collection
	c, ok := api.CMap[cName]
	if ok {
		err := c.FrameDelete(frameName)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		statusIsOK(w, http.StatusOK, cName, "", "frame-delete", frameName)
		return
	}
	// Collection not found
	http.NotFound(w, r)
	return
}

//**************************************************************
// The following routes handle JSON object versions. These are
// placeholders for some future release of 2.X.
//**************************************************************

func ObjectVersions(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	statusIsError(w, r, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented, "")
}

func ReadVersion(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	statusIsError(w, r, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented, "")
}

func DeleteVersion(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	statusIsError(w, r, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented, "")
}

//**************************************************************
// The following routes handle attachment versions. These are
// placeholders for some future release of 2.X.
//**************************************************************

func AttachmentVersions(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	statusIsError(w, r, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented, "")
}

func RetrieveVersion(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	statusIsError(w, r, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented, "")
}

func PruneVersion(w http.ResponseWriter, r *http.Request, api *API, cName, verb string, options []string) {
	statusIsError(w, r, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented, "")
}
