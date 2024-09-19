// api is a part of dataset
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
package dataset

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
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
	Settings *Settings
	// CMap is a map to the collections supported by the web service.
	CMap map[string]*Collection

	// Routes holds a double map of prefix path and HTTP method that
	// points to the function that will be dispatched if found.
	//
	// The the first level map identifies the prefix path for the route
	// e.g. "api/version".  No leading slash is expected.
	// The second level map is organized by HTTP method, e.g. "GET",
	// "POST". The second map points to the function to call when
	// the route and method matches.
	Routes map[string]map[string]func(http.ResponseWriter, *http.Request, *API, string, string, []string)

	// Debug if set true will cause more verbose output.
	Debug bool

	// Process ID
	Pid int
}

var (
	settings *Settings
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
		// See if we need to set a header of JavaScript or TypeScript files.
		if strings.HasSuffix(r.URL.Path, ".js") || strings.HasSuffix(r.URL.Path, ".mjs") {
			w.Header().Add("Content-Type", "application/javascript; charset=utf-8")
		}
		if strings.HasSuffix(r.URL.Path, ".ts") {
			w.Header().Add("Content-Type", "application/typescript; charset=utf-8")
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
		default:
			cName, verb, options = parts[0], parts[1], parts[2:]
		}
		if api.Debug {
			log.Printf("DEBUG cName %q, verb: %q, options: %+v\n", cName, verb, options)
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
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
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
	//api.Collections = map[string]*Collection{}
	log.Printf(`Shutdown completed %s pid: %d exit code: %d `, appName, pid, exitCode)
	return exitCode
}

// Reload performs a Shutdown and an init after re-reading
// in the settings.yaml file.
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
//
//	func Version(w http.ResponseWriter, r *http.Reqest, api *API, verb string, options []string) {
//	   ...
//	}
//
//	...
//
//	err := api.RegistereRoute("version", http.MethodGet, Version)
//	if err != nil {
//	   ...
//	}
//
// ```
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

	settings, err := ConfigOpen(settingsFile)
	if err != nil {
		return err
	}
	if settings == nil {
		return fmt.Errorf("ConfigOpen(%q) returned nil, nil", settingsFile)
	}
	api.Settings = settings

	// Setup an empty request router
	api.Routes = make(map[string]map[string]func(http.ResponseWriter, *http.Request, *API, string, string, []string))

	// We always should have a version route!
	err = api.RegisterRoute("version", http.MethodGet, ApiVersion)
	if err != nil {
		return err
	}
	err = api.RegisterRoute("collections", http.MethodGet, Collections)
	if err != nil {
		return err
	}
	err = api.RegisterRoute("codemeta", http.MethodGet, Codemeta)
	if err != nil {
		return err
	}

	// Get out cName from config
	for _, cfg := range api.Settings.Collections {
		if api.CMap == nil {
			api.CMap = make(map[string]*Collection)
		}
		// NOTE: cName is the name used in our CMap as well as in building
		// paths for service.
		cName := path.Base(cfg.CName)
		c, err := Open(cfg.CName)
		if err != nil {
			log.Printf("WARNING: failed to open %q, %s", cfg.CName, err)
		} else {
			api.CMap[cName] = c
		}
		// NOTE: Need to review the permissions in cfg and then
		// add the appropriate routes to api.routes.
		if cfg.Keys {
			prefix := path.Join(cName, "keys")
			if err = api.RegisterRoute(prefix, http.MethodGet, Keys); err != nil {
				return err
			}
		}
		if cfg.Create {
			prefix := path.Join(cName, "object")
			if err = api.RegisterRoute(prefix, http.MethodPost, Create); err != nil {
				return err
			}
		}
		if cfg.Read {
			prefix := path.Join(cName, "object")
			if err = api.RegisterRoute(prefix, http.MethodGet, Read); err != nil {
				return err
			}
		}
		if cfg.Update {
			prefix := path.Join(cName, "object")
			if err = api.RegisterRoute(prefix, http.MethodPut, Update); err != nil {
				return err
			}
		}
		if cfg.Delete {
			prefix := path.Join(cName, "object")
			if err = api.RegisterRoute(prefix, http.MethodDelete, Delete); err != nil {
				return err
			}
		}
		if cfg.QueryFn != nil && len(cfg.QueryFn) > 0 {
			prefix := path.Join(cName, "query")
			if err = api.RegisterRoute(prefix, http.MethodPost, Query); err != nil {
				return err
			}
		}
		if cfg.Attachments {
			prefix := path.Join(cName, "attachments")
			if err = api.RegisterRoute(prefix, http.MethodGet, Attachments); err != nil {
				return err
			}
		}
		if cfg.Retrieve {
			prefix := path.Join(cName, "attachment")
			if err = api.RegisterRoute(prefix, http.MethodGet, Retrieve); err != nil {
				return err
			}
		}
		if cfg.Attach {
			prefix := path.Join(cName, "attachment")
			if err = api.RegisterRoute(prefix, http.MethodPost, Attach); err != nil {
				return err
			}
		}
		if cfg.Prune {
			prefix := path.Join(cName, "attachment")
			if err = api.RegisterRoute(prefix, http.MethodDelete, Prune); err != nil {
				return err
			}
		}
		if cfg.FrameRead {
			prefix := path.Join(cName, "frames")
			if err = api.RegisterRoute(prefix, http.MethodGet, Frames); err != nil {
				return err
			}
			prefix = path.Join(cName, "has-frame")
			if err = api.RegisterRoute(prefix, http.MethodGet, HasFrame); err != nil {
				return err
			}
			prefix = path.Join(cName, "frame-objects")
			if err = api.RegisterRoute(prefix, http.MethodGet, FrameObjects); err != nil {
				return err
			}
			prefix = path.Join(cName, "frame")
			if err = api.RegisterRoute(prefix, http.MethodGet, FrameDef); err != nil {
				return err
			}
			prefix = path.Join(cName, "frame-keys")
			if err = api.RegisterRoute(prefix, http.MethodGet, FrameKeys); err != nil {
				return err
			}
		}
		if cfg.FrameWrite {
			prefix := path.Join(cName, "frame")
			if err = api.RegisterRoute(prefix, http.MethodPost, FrameCreate); err != nil {
				return err
			}
			if err = api.RegisterRoute(prefix, http.MethodPut, FrameUpdate); err != nil {
				return err
			}
			if err = api.RegisterRoute(prefix, http.MethodDelete, FrameDelete); err != nil {
				return err
			}
		}

		//NOTE: Need to apply versioned routes if a collection uses
		// versioning and permission is granted to see versions. I am
		// not planning on implementing version support in the web
		// service until it is needed, e.g. at some future 2.x release.
		if cfg.Versions && c != nil && c.Versioning != "" {
			prefix := path.Join(cName, "object-versions")
			if err = api.RegisterRoute(prefix, http.MethodGet, ObjectVersions); err != nil {
				return err
			}
			if cfg.Read {
				prefix := path.Join(cName, "object-version")
				if err = api.RegisterRoute(prefix, http.MethodGet, ReadVersion); err != nil {
					return err
				}
			}
			if cfg.Delete {
				prefix := path.Join(cName, "object-version")
				if err = api.RegisterRoute(prefix, http.MethodDelete, DeleteVersion); err != nil {
					return err
				}
			}

			if cfg.Attachments {
				prefix := path.Join(cName, "attachment-versions")
				if err = api.RegisterRoute(prefix, http.MethodGet, AttachmentVersions); err != nil {
					return err
				}
			}
			if cfg.Retrieve {
				prefix := path.Join(cName, "attachment-version")
				if err = api.RegisterRoute(prefix, http.MethodGet, RetrieveVersion); err != nil {
					return err
				}
			}
			if cfg.Prune {
				prefix := path.Join(cName, "attachment-version")
				if err = api.RegisterRoute(prefix, http.MethodDelete, PruneVersion); err != nil {
					return err
				}
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
//
//	appName := path.Base(sys.Argv[0])
//	settingsFile := "settings.yaml"
//	if err := api.RunAPI(appName, settingsFile); err != nil {
//	   ...
//	}
//
// ```
func RunAPI(appName string, settingsFile string, debug bool) error {
	api := new(API)
	// Open collection
	if err := api.Init(appName, settingsFile); err != nil {
		return err
	}
	api.Debug = debug

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
