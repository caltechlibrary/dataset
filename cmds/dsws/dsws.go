//
// dsws.go - A web server/service for hosting dataset search and related static pages.
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2017, Caltech
// All rights not granted herein are expressly reserved by Caltech
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
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"text/template"

	// Caltech Library packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/dataset"
	"github.com/caltechlibrary/mkpage"
	"github.com/caltechlibrary/tmplfn"
)

// Flag options
var (
	usage = `USAGE: %s [OPTIONS] [DOCROOT]`

	description = `
SYNOPSIS

	web service for search dataset collections

%s which support a web search of a dataset collection.

CONFIGURATION

%s can be configurated through environment settings. The following are
supported.

+ DATASET_URL  - sets the URL to listen on (e.g. http://localhost:8000)
+ DATASET_DOCROOT - sets the document path to use
+ DATASET_SSL_KEY - the path to the SSL key if using https
+ DATASET_SSL_CERT - the path to the SSL cert if using https
+ DATASET_INDEXES - a list of Bleve indexes available to query
+ DATASET_TEMPLATES - directory holding the templates
`

	examples = `
EXAMPLES

Run web server using the content in the current directory
(assumes the environment variables DATASET_DOCROOT are not defined).

   %s

Run web service using "index.bleve" index, results templates in 
"templates" direcotry and a "htdocs" directory for static files.

   %s -indexes index.bleve -templates templates htdocs
`

	// Standard options
	showHelp    bool
	showVersion bool
	showLicense bool

	// local app options
	uri         string
	docRoot     string
	sslKey      string
	sslCert     string
	indexNames  string
	templateDir string
)

func logRequest(r *http.Request) {
	log.Printf("Request: %s Path: %s RemoteAddr: %s UserAgent: %s\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		next.ServeHTTP(w, r)
	})
}

func init() {
	defaultDocRoot := "."
	defaultURL := "http://localhost:8000"

	flag.BoolVar(&showHelp, "h", false, "Display this help message")
	flag.BoolVar(&showHelp, "help", false, "Display this help message")
	flag.BoolVar(&showVersion, "v", false, "Should version info")
	flag.BoolVar(&showVersion, "version", false, "Should version info")
	flag.BoolVar(&showLicense, "l", false, "Should license info")
	flag.BoolVar(&showLicense, "license", false, "Should license info")
	flag.StringVar(&docRoot, "d", defaultDocRoot, "Set the htdocs path")
	flag.StringVar(&docRoot, "docs", defaultDocRoot, "Set the htdocs path")
	flag.StringVar(&uri, "u", defaultURL, "The protocal and hostname listen for as a URL")
	flag.StringVar(&uri, "url", defaultURL, "The protocal and hostname listen for as a URL")
	flag.StringVar(&sslKey, "k", "", "Set the path for the SSL Key")
	flag.StringVar(&sslKey, "key", "", "Set the path for the SSL Key")
	flag.StringVar(&sslCert, "c", "", "Set the path for the SSL Cert")
	flag.StringVar(&sslCert, "cert", "", "Set the path for the SSL Cert")

	flag.StringVar(&indexNames, "indexes", "", "A colon delimited list of Bleve indexes")
	flag.StringVar(&indexNames, "i", "", "A colon delimited list of Bleve indexes")

	flag.StringVar(&templateDir, "templates", "", "the path to the templates directory")
	flag.StringVar(&templateDir, "t", "", "the path to the templates directory")
}

func loadTemplates(templateDir string) (*template.Template, error) {
	templateNames := []string{}

	tMaps := tmplfn.Join(tmplfn.TimeMap, tmplfn.PageMap)

	files, err := ioutil.ReadDir(templateDir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		fName := file.Name()
		ext := path.Ext(fName)
		if ext == ".tmpl" {
			templateNames = append(templateNames, fName)
		}
	}
	return tmplfn.Assemble(tMaps, templateNames...)
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()
	args := flag.Args()

	// Configuration and command line interation
	cfg := cli.New(appName, "DATASET", fmt.Sprintf(dataset.License, appName, dataset.Version), dataset.Version)
	cfg.UsageText = fmt.Sprintf(usage, appName)
	cfg.DescriptionText = fmt.Sprintf(description, appName, appName)
	cfg.ExampleText = fmt.Sprintf(examples, appName, appName)

	// Process flags and update the environment as needed.
	if showHelp == true {
		fmt.Println(cfg.Usage())
		os.Exit(0)
	}
	if showLicense == true {
		fmt.Println(cfg.License())
		os.Exit(0)
	}
	if showVersion == true {
		fmt.Println(cfg.Version())
		os.Exit(0)
	}

	// setup from command line
	if len(args) > 0 {
		docRoot = args[0]
	}

	docRoot = cfg.CheckOption(docRoot, cfg.MergeEnv("docroot", docRoot), true)
	log.Printf("DocRoot %s", docRoot)
	indexNames = cfg.CheckOption(indexNames, cfg.MergeEnv("indexes", indexNames), true)
	log.Printf("Indexes %s", docRoot)
	templateDir = cfg.CheckOption(templateDir, cfg.MergeEnv("templates", templateDir), true)
	log.Printf("Templates %s", docRoot)

	uri = cfg.CheckOption(uri, cfg.MergeEnv("url", uri), true)
	u, err := url.Parse(uri)
	if err != nil {
		log.Fatalf("Can't parse %q, %s", uri, err)
	}

	log.Printf("Listening for %s", uri)
	if u.Scheme == "https" {
		sslKey = cfg.CheckOption(sslKey, cfg.MergeEnv("ssl_key", sslKey), true)
		sslCert = cfg.CheckOption(sslCert, cfg.MergeEnv("ssl_cert", sslCert), true)
		log.Printf("SSL Key %s", sslKey)
		log.Printf("SSL Cert %s", sslCert)
	}

	// Open the indexes for reading
	idxAlias, err := dataset.OpenIndexes(strings.Split(indexNames, ":"))
	if err != nil {
		log.Fatalf("Can't open indexes, %s", err)
	}
	defer idxAlias.Close()

	// Load and validate the templates for using in the searchHandler
	tmpl, err := loadTemplates(templateDir)
	if err != nil {
		log.Fatalf("Can't load templates, %s", err)
	}

	// Construct our handler
	searchHandler := func(w http.ResponseWriter, r *http.Request) {
		opts := map[string]string{}
		values := r.URL.Query()
		format := values.Get("fmt")
		qString := values.Get("q")
		// Get the options understood by dataset.Find()
		for _, ky := range []string{"size", "from", "ids", "sort", "explain", "fields", "highlight"} {
			if v := values.Get(ky); v != "" {
				opts[ky] = v
			}
		}
		buf := bytes.NewBufferString("")
		results, err := dataset.Find(buf, idxAlias, []string{qString}, opts)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), 500)
		}

		// Based on the request info, format the results appropriately
		switch strings.ToLower(format) {
		case "csv":
			fields := strings.Split(values.Get("fields"), ":")
			if len(fields) == 0 {
				fields = []string{"*"}
			}
			if err := dataset.CSVFormatter(w, results, fields); err != nil {
				http.Error(w, fmt.Sprintf("%s", err), 500)
			}
		case "json":
			if err := dataset.JSONFormatter(w, results); err != nil {
				http.Error(w, fmt.Sprintf("%s", err), 500)
			}
		case "html":
			if err := dataset.HTMLFormatter(w, results, tmpl); err != nil {
				http.Error(w, fmt.Sprintf("%s", err), 500)
			}
		case "include":
			if err := dataset.IncludeFormatter(w, results, tmpl); err != nil {
				http.Error(w, fmt.Sprintf("%s", err), 500)
			}
		default:
			// Assume plain text results
			fmt.Fprintf(w, "%s\n", results)
		}
	}

	// Define our search API prefix path
	wsapi := new(mkpage.WSAPI)

	if err := wsapi.AddRoute("/q", searchHandler); err != nil {
		log.Fatal("can't add search route, %s", err)
	}

	http.Handle("/", http.FileServer(http.Dir(docRoot)))
	if u.Scheme == "https" {
		err := http.ListenAndServeTLS(u.Host, sslCert, sslKey, mkpage.RequestLogger(wsapi.Router(mkpage.StaticRouter(http.DefaultServeMux))))
		if err != nil {
			log.Fatalf("%s", err)
		}
	} else {
		err := http.ListenAndServe(u.Host, mkpage.RequestLogger(wsapi.Router(mkpage.StaticRouter(http.DefaultServeMux))))
		if err != nil {
			log.Fatalf("%s", err)
		}
	}
}
