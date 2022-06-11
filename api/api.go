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
)

const (
	// jsonSizeLimit is the maximum size JSON object we'll accept via
	// our service. Current 1 MB (2^20)
	jsonSizeLimit = 1048576

	// attachmentSizeLimit is the maximum size of Attachments we'll
	// accept via our service. Current 250 MiB
	attachmentSizeLimit = (jsonSizeLimit * 250)
)

type API struct {
	AppName  string
	Settings string
	Version  string
	Cfg      *Config
	CName    string
	c        *ds.Collection

	// Routes holds a double map of prefix path and HTTP method that
	// points to the function that will be dispatched if found.
	//
	// The the first level map identifies the prefix path for the route
	// e.g. "api/version".  No leading slash is expected.
	// The second level map is organized by HTTP method, e.g. "GET",
	// "POST". The second map points to the function to call when
	// the route and method matches.
	Routes map[string]map[string]func(http.ResponseWriter, *http.Request, *API, string, []string)

	// Process ID
	Pid int
}

var (
	config *Config
)

// IsDotPath checks to see if a path is requested with a dot file (e.g. docs/.git/* or docs/.htaccess)
func IsDotPath(p string) bool {
	for _, part := range strings.Split(path.Clean(p), "/") {
		if strings.HasPrefix(part, "..") == false && strings.HasPrefix(part, ".") == true && len(part) > 1 {
			return true
		}
	}
	return false
}

// RequestLogger logs the request based on the request object passed into it.
func RequestLogger(next http.Handler) http.Handler {
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

// ResponseLogger logs the response based on a request, status and error message
func ResponseLogger(r *http.Request, status int, err error) {
	q := r.URL.Query()
	if len(q) > 0 {
		log.Printf("Response: %s Path: %s RemoteAddr: %s UserAgent: %s Query: %+v Status: %d, %s %q\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), q, status, http.StatusText(status), err)
	} else {
		log.Printf("Response: %s Path: %s RemoteAddr: %s UserAgent: %s Status: %d, %s %q\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), status, http.StatusText(status), err)
	}
}

// StaticRouter scans the request object to either add a .html extension or prevent serving a dot file path
func StaticRouter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If given a dot file path, send forbidden
		if IsDotPath(r.URL.Path) == true {
			http.Error(w, "Forbidden", 403)
			ResponseLogger(r, 403, fmt.Errorf("Forbidden, requested a dot path"))
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
			verb    string
			options []string
		)
		switch len(parts) {
		case 1:
			verb, options = parts[0], []string{}
		case 2:
			verb, options = parts[0], parts[1:]
		default:
			http.NotFound(w, r)
			return
		}
		if route, ok := api.Routes[verb]; ok {
			if fn, ok := route[r.Method]; ok {
				fn(w, r, api, verb, options)
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
	host := api.Cfg.Host
	// Define Routes here.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		api.Router(w, r)
	})
	log.Printf("%s start, listening on %s", api.AppName, api.Cfg.Host)
	return http.ListenAndServe(host, RequestLogger(mux))
}

// Shutdown attemtps a graceful shutdown of the service.
// returns an exit code.
func (api *API) Shutdown(sigName string) int {
	appName := api.AppName
	exitCode := 0
	pid := os.Getpid()
	log.Printf(`Received signal %s`, sigName)
	log.Printf(`Closing dataset connection %s pid: %d`, appName, pid)
	if err := api.c.Close(); err != nil {
		exitCode = 1
	}
	api.c = nil
	log.Printf(`Shutdown completed %s pid: %d exit code: %d `, appName, pid, exitCode)
	return exitCode
}

// Reload performs a Shutdown and an init after re-reading
// in the settings.json file.
func (api *API) Reload(sigName string) error {
	appName := api.AppName
	settings := api.Settings
	exitCode := api.Shutdown(sigName)
	if exitCode != 0 {
		return fmt.Errorf("Reload failed, could not shutdown the current processes")
	}
	// Reload the configuration
	cfg, err := LoadConfig(settings)
	if err != nil {
		log.Fatal(err)
	}
	api.Cfg = cfg
	_, err = api.Init(appName, cfg)
	if err != nil {
		return err
	}
	defer func() {
		if api.c != nil {
			api.c.Close()
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
func (api *API) RegisterRoute(prefix string, method string, fn func(http.ResponseWriter, *http.Request, *API, string, []string)) error {
	if _, ok := api.Routes[prefix]; ok {
		if _, ok := api.Routes[prefix][method]; ok {
			return fmt.Errorf("%q %s already registered", prefix, method)
		}
		api.Routes[prefix][method] = fn
		return nil
	}
	api.Routes[prefix] = make(map[string]func(http.ResponseWriter, *http.Request, *API, string, []string))
	api.Routes[prefix][method] = fn
	return nil
}

// Init setups and the API to run.
func (api *API) Init(appName string, cfg *Config) (*ds.Collection, error) {
	var err error
	api.AppName = path.Base(appName)
	api.Version = Version

	// Get out cName from config
	api.CName = cfg.CName
	api.Pid = os.Getpid()
	api.c, err = ds.Open(api.CName)
	if err != nil {
		return nil, err
	}
	// Setup an empty request router
	api.Routes = make(map[string]map[string]func(http.ResponseWriter, *http.Request, *API, string, []string))

	// NOTE: Need to review the permissions in cfg and then
	// add the appropriate routes to api.routes.
	err = api.RegisterRoute("version", http.MethodGet, ApiVersion)
	if err != nil {
		return nil, err
	}
	if cfg.Keys {
		if err = api.RegisterRoute("keys", http.MethodGet, ApiKeys); err != nil {
			return nil, err
		}
	}
	if cfg.Create {
		if err = api.RegisterRoute("object", http.MethodPost, ApiCreate); err != nil {
			return nil, err
		}
	}
	if cfg.Read {
		if err = api.RegisterRoute("object", http.MethodGet, ApiRead); err != nil {
			return nil, err
		}
	}
	if cfg.Update {
		if err = api.RegisterRoute("object", http.MethodPut, ApiUpdate); err != nil {
			return nil, err
		}
	}
	if cfg.Delete {
		if err = api.RegisterRoute("object", http.MethodDelete, ApiDelete); err != nil {
			return nil, err
		}
	}
	return api.c, nil
}

// RunAPI takes a JSON configuration file and opens
// all the collections to be used by web service.
//
// ```
//   appName := path.Base(sys.Argv[0])
//   configFile := "settings.json"
//   if err := api.RunAPI(appName, settings); err != nil {
//      ...
//   }
// ```
//
func RunAPI(appName string, settings string) error {
	var err error

	api := new(API)

	cfg, err := LoadConfig(settings)
	if err != nil {
		return err
	}
	api.Cfg = cfg

	// Open collection
	c, err := api.Init(appName, cfg)
	if err != nil {
		return err
	}
	defer func() {
		if c != nil {
			c.Close()
		}
	}()

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
