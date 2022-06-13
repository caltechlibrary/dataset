//
// api is a submodule of dataset
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2022, Caltech
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
package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	// Caltech Library packages
	ds "github.com/caltechlibrary/dataset"
	"github.com/caltechlibrary/dataset/config"
)

const (
	// jsonSizeLimit is the maximum size JSON object we'll accept via
	// our service. Current 1 MB (2^20)
	jsonSizeLimit = 1048576

	// attachmentSizeLimit is the maximum size of Attachments we'll
	// accept via our service. Current 250 MiB
	attachmentSizeLimit = (jsonSizeLimit * 250)
)

// API this structure holds the information for running an
// web service instance. One web service may host many collections.
type API struct {
	// AppName is the name of the running application. E.g. os.Args[0]
	AppName string
	// SettingsFile is the path to the settings file.
	SettingsFile string
	// Version is the version of the API running
	Version string
	// Settings is the configuration reading from SettingsFile
	Settings *config.Settings
	// CMap is a map to the collections supported by the web service.
	CMap map[string]*ds.Collection

	// Routes holds a double map of prefix path and HTTP method that
	// points to the function that will be dispatched if found.
	//
	// The the first level map identifies the prefix path for the route
	// e.g. "api/version".  No leading slash is expected.
	// The second level map is organized by HTTP method, e.g. "GET",
	// "POST". The second map points to the function to call when
	// the route and method matches.
	Routes map[string]map[string]func(http.ResponseWriter, *http.Request, *API, string, string, []string)

	// Process ID
	Pid int
}

var (
	settings *config.Settings
)

// hasDotPath checks to see if a path is requested with a dot file (e.g. docs/.git/* or docs/.htaccess)
func hasDotPath(p string) bool {
	for _, part := range strings.Split(path.Clean(p), "/") {
		if strings.HasPrefix(part, "..") == false && strings.HasPrefix(part, ".") == true && len(part) > 1 {
			return true
		}
	}
	return false
}

// requestLogger logs the request based on the request object passed into it.
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if len(q) > 0 {
			log.Printf("Request: %s Path: %s RemoteAddr: %s UserAgent: %s Query: %+v\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), q)
		} else {
			log.Printf("Request: %s Path: %s RemoteAddr: %s UserAgent: %s\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
		}
		next.ServeHTTP(w, r)
	})
}

// responseLogger logs the response based on a request, status and error message
func responseLogger(r *http.Request, status int, err error) {
	q := r.URL.Query()
	if len(q) > 0 {
		log.Printf("Response: %s Path: %s RemoteAddr: %s UserAgent: %s Query: %+v Status: %d, %s %q\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), q, status, http.StatusText(status), err)
	} else {
		log.Printf("Response: %s Path: %s RemoteAddr: %s UserAgent: %s Status: %d, %s %q\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), status, http.StatusText(status), err)
	}
}

// staticRouter scans the request object to either add a .html extension or prevent serving a dot file path
func staticRouter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If given a dot file path, send forbidden
		if hasDotPath(r.URL.Path) == true {
			http.Error(w, "Forbidden", 403)
			responseLogger(r, 403, fmt.Errorf("Forbidden, requested a dot path"))
			return
		}
		// If we make it this far, fall back to the default handler
		next.ServeHTTP(w, r)
	})
}

func (api *API) Router(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		relPath := strings.TrimPrefix(r.URL.Path, "/api/")
		parts := strings.Split(relPath, "/")
		if len(parts) < 1 {
			http.NotFound(w, r)
			return
		}
		var (
			cName   string
			verb    string
			options []string
		)
		switch len(parts) {
		case 1:
			cName, verb, options = "", parts[0], []string{}
		case 2:
			cName, verb, options = parts[0], parts[1], []string{}
		case 3:
			cName, verb, options = parts[0], parts[1], parts[2:]
		default:
			http.NotFound(w, r)
			return
		}
		prefix := path.Join(cName, verb)
		if route, ok := api.Routes[prefix]; ok {
			if fn, ok := route[r.Method]; ok {
				fn(w, r, api, cName, verb, options)
				return
			}
			code := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(code), code)
			return
		}
	}
	http.NotFound(w, r)
}

// WebService this starts and runs a web server implementation
// of dataset.
func (api *API) WebService() error {
	mux := http.NewServeMux()
	host := api.Settings.Host
	// Define Routes here.
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		api.Router(w, r)
	})
	if api.Settings.Htdocs != "" {
		mux.Handle("/", staticRouter(http.FileServer(http.Dir(api.Settings.Htdocs))))
	}
	log.Printf("%s start, listening on %s", api.AppName, api.Settings.Host)
	return http.ListenAndServe(host, requestLogger(mux))
}

// Shutdown attemtps a graceful shutdown of the service.
// returns an exit code.
func (api *API) Shutdown(sigName string) int {
	appName := api.AppName
	exitCode := 0
	pid := os.Getpid()
	log.Printf(`Received signal %s`, sigName)
	log.Printf(`Closing dataset connection %s pid: %d`, appName, pid)

	for cName, c := range api.CMap {
		if c != nil {
			if err := c.Close(); err != nil {
				log.Printf("WARNING: error closing %q, %s", cName, err)
				exitCode = 1
			}
		}
	}
	//api.Collections = map[string]*ds.Collection{}
	log.Printf(`Shutdown completed %s pid: %d exit code: %d `, appName, pid, exitCode)
	return exitCode
}

// Reload performs a Shutdown and an init after re-reading
// in the settings.json file.
func (api *API) Reload(sigName string) error {
	appName := api.AppName
	settingsFile := api.SettingsFile
	exitCode := api.Shutdown(sigName)
	if exitCode != 0 {
		return fmt.Errorf("Reload failed, could not shutdown the current processes")
	}
	// Reload the configuration
	if err := api.Init(appName, settingsFile); err != nil {
		return err
	}
	defer func() {
		for cName, c := range api.CMap {
			if c != nil {
				if err := c.Close(); err != nil {
					log.Printf("WARNING: error closing %q, %s", cName, err)
				}
			}
		}
	}()
	return api.WebService()
}

// RegisterRoute resigns a prefix path to a route handler.
//
// prefix is the url path prefix minus the leading slash that
// is targetted by this handler.
//
// method is the HTTP method the func will process
// fn is the function that handles this route.
//
// ```
//     func Version(w http.ResponseWriter, r *http.Reqest, api *API, verb string, options []string) {
//        ...
//     }
//
//     ...
//
//     err := api.RegistereRoute("version", http.MethodGet, Version)
//     if err != nil {
//        ...
//     }
// ```
//
func (api *API) RegisterRoute(prefix string, method string, fn func(http.ResponseWriter, *http.Request, *API, string, string, []string)) error {
	if _, ok := api.Routes[prefix]; ok {
		if _, ok := api.Routes[prefix][method]; ok {
			return fmt.Errorf("%q %s already registered", prefix, method)
		}
		api.Routes[prefix][method] = fn
		return nil
	}
	api.Routes[prefix] = make(map[string]func(http.ResponseWriter, *http.Request, *API, string, string, []string))
	api.Routes[prefix][method] = fn
	return nil
}

// Init setups and the API to run.
func (api *API) Init(appName string, settingsFile string) error {
	var err error
	api.AppName = path.Base(appName)
	api.Version = Version
	api.Pid = os.Getpid()
	api.SettingsFile = settingsFile

	settings, err := config.Open(settingsFile)
	if err != nil {
		return err
	}
	if settings == nil {
		return fmt.Errorf("Open(%q) returned nil, nil", settingsFile)
	}
	api.Settings = settings

	// Setup an empty request router
	api.Routes = make(map[string]map[string]func(http.ResponseWriter, *http.Request, *API, string, string, []string))

	// We always should have a version route!
	err = api.RegisterRoute("version", http.MethodGet, ApiVersion)
	if err != nil {
		return err
	}
	err = api.RegisterRoute("collections", http.MethodGet, ApiCollections)
	if err != nil {
		return err
	}

	// Get out cName from config
	for _, cfg := range api.Settings.Collections {
		if api.CMap == nil {
			api.CMap = make(map[string]*ds.Collection)
		}
		// NOTE: cName is the name used in our CMap as well as in building
		// paths for service.
		cName := path.Base(cfg.CName)
		c, err := ds.Open(cfg.CName)
		if err != nil {
			log.Printf("WARNING: failed to open %q, %s", cfg.CName, err)
		} else {
			api.CMap[cName] = c
		}
		// NOTE: Need to review the permissions in cfg and then
		// add the appropriate routes to api.routes.
		if cfg.Keys {
			prefix := path.Join(cName, "keys")
			if err = api.RegisterRoute(prefix, http.MethodGet, ApiKeys); err != nil {
				return err
			}
		}
		if cfg.Create {
			prefix := path.Join(cName, "object")
			if err = api.RegisterRoute(prefix, http.MethodPost, ApiCreate); err != nil {
				return err
			}
		}
		if cfg.Read {
			prefix := path.Join(cName, "object")
			if err = api.RegisterRoute(prefix, http.MethodGet, ApiRead); err != nil {
				return err
			}
		}
		if cfg.Update {
			prefix := path.Join(cName, "object")
			if err = api.RegisterRoute(prefix, http.MethodPut, ApiUpdate); err != nil {
				return err
			}
		}
		if cfg.Delete {
			prefix := path.Join(cName, "object")
			if err = api.RegisterRoute(prefix, http.MethodDelete, ApiDelete); err != nil {
				return err
			}
		}
	}
	if len(api.CMap) == 0 {
		return fmt.Errorf("failed to open any collections")
	}
	return nil
}

// RunAPI takes a JSON configuration file and opens
// all the collections to be used by web service.
//
// ```
//   appName := path.Base(sys.Argv[0])
//   settingsFile := "settings.json"
//   if err := api.RunAPI(appName, settingsFile); err != nil {
//      ...
//   }
// ```
//
func RunAPI(appName string, settingsFile string) error {
	api := new(API)
	// Open collection
	if err := api.Init(appName, settingsFile); err != nil {
		return err
	}

	// Listen for Ctr-C
	processControl := make(chan os.Signal, 1)
	signal.Notify(processControl, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	go func() {
		sig := <-processControl
		switch sig {
		case syscall.SIGINT:
			os.Exit(api.Shutdown(sig.String()))
		case syscall.SIGTERM:
			os.Exit(api.Shutdown(sig.String()))
		case syscall.SIGHUP:
			if err := api.Reload(sig.String()); err != nil {
				log.Println(err)
				os.Exit(1)
			}
		default:
			os.Exit(api.Shutdown(sig.String()))
		}
	}()

	// Run the service
	return api.WebService()
}