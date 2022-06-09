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
	// Routes holds a prefix for key pointing at a "route handler" func
	// that is built from a response writer, request, two string (
	// collection name and verb) and an array of string (options).
	// The "route handler" evalutes the collection name, verb and
	// options against the request and response via the response writer.
	//
	// The prefix is built from the first to path parts atfter "/api/"
	// in the request URL. No leading slash is expected.
	Routes map[string]func(http.ResponseWriter, *http.Request, string, string, []string)
	Pid    int
}

var (
	config *Config
)

func (api *API) Router(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/"), "/")
	if len(parts) < 2 {
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
		cName, verb, options = parts[0], "", []string{}
	case 2:
		cName, verb, options = parts[0], parts[1], []string{}
	case 3:
		cName, verb, options = parts[0], parts[1], parts[2:]
	default:
		http.NotFound(w, r)
		return
	}
	prefix := path.Join(cName, verb)
	if fn, ok := api.Routes[prefix]; ok {
		fn(w, r, cName, verb, options)
	} else {
		http.NotFound(w, r)
		return
	}
}

// WebService this starts and runs a web server implementation
// of dataset.
func (api *API) WebService() error {
	mux := http.NewServeMux()
	host := api.Cfg.Host
	// Define Routes here.
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		api.Router(w, r)
	})
	return http.ListenAndServe(host, mux)
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

// Init setups and the API to run.
func (api *API) Init(appName string, cfg *Config) (*ds.Collection, error) {
	var err error
	api.AppName = appName
	api.Version = Version

	// Get out cName from config
	api.CName = cfg.CName
	api.Pid = os.Getpid()
	api.c, err = ds.Open(api.CName)
	if err != nil {
		return nil, err
	}

	// FIXME: Need to review the permissions in cfg and then
	// add the appropriate routes to api.routes.

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
		log.Fatal(err)
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
	/* Listen for Ctr-C */
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
	return api.WebService()
}
