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
	"encoding/json" // DEBUG
	"flag"
	"fmt"
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
	usage = `USAGE: %s [OPTIONS] DOC_ROOT BLEVE_INDEXES`

	description = `
SYNOPSIS

	web service for search dataset collections

%s which support a web search of a dataset collection.

CONFIGURATION

%s can be configurated through environment settings. The following are
supported.

+ DATASET_URL  - (optional) sets the URL to listen on (e.g. http://localhost:8011)
+ DATASET_SSL_KEY - (optional) the path to the SSL key if using https
+ DATASET_SSL_CERT - (optional) the path to the SSL cert if using https
+ DATASET_TEMPLATE - (optional) path to search results template
`

	examples = `
EXAMPLES

Run web server using the content in the current directory
(assumes the environment variables DATASET_DOCROOT are not defined).

   %s

Run web service using "index.bleve" index, results templates in 
"templates/search.tmpl" and a "htdocs" directory for static files.

   %s -template templates/search.tmpl htdocs index.bleve
`

	// Standard options
	showHelp    bool
	showVersion bool
	showLicense bool

	// local app options
	uri         string
	sslKey      string
	sslCert     string
	searchTName string

	// Provided as an ordered command line arg
	docRoot    string
	indexNames []string
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
	defaultURL := "http://localhost:8011"

	// Standard Options
	flag.BoolVar(&showHelp, "h", false, "Display this help message")
	flag.BoolVar(&showHelp, "help", false, "Display this help message")
	flag.BoolVar(&showVersion, "v", false, "Should version info")
	flag.BoolVar(&showVersion, "version", false, "Should version info")
	flag.BoolVar(&showLicense, "l", false, "Should license info")
	flag.BoolVar(&showLicense, "license", false, "Should license info")

	// App Options
	flag.StringVar(&uri, "u", defaultURL, "The protocal and hostname listen for as a URL")
	flag.StringVar(&uri, "url", defaultURL, "The protocal and hostname listen for as a URL")
	flag.StringVar(&sslKey, "k", "", "Set the path for the SSL Key")
	flag.StringVar(&sslKey, "key", "", "Set the path for the SSL Key")
	flag.StringVar(&sslCert, "c", "", "Set the path for the SSL Cert")
	flag.StringVar(&sslCert, "cert", "", "Set the path for the SSL Cert")
	flag.StringVar(&searchTName, "template", "", "the path to the search result template")
	flag.StringVar(&searchTName, "t", "", "the path to the search result template")
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
	if len(args) > 1 {
		indexNames = args[1:]
	} else {
		fmt.Fprintf(os.Stderr, cfg.UsageText)
		fmt.Fprintf(os.Stderr, "error: one or more Bleve index is required\n")
		os.Exit(1)
	}

	log.Printf("DocRoot %s", docRoot)
	if len(indexNames) == 1 {
		log.Printf("Index %s", strings.Join(indexNames, ", "))
	} else {
		log.Printf("Indexes %s", strings.Join(indexNames, ", "))
	}
	searchTName = cfg.CheckOption(searchTName, cfg.MergeEnv("template", searchTName), false)
	if searchTName != "" {
		log.Printf("Search result template is %q", searchTName)
	} else {
		log.Printf("Using default search result template")
	}

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
	idxAlias, err := dataset.OpenIndexes(indexNames)
	if err != nil {
		log.Fatalf("Can't open indexes, %s", err)
	}
	defer idxAlias.Close()

	var searchTmpl *template.Template

	// Load and validate the templates for using in the searchHandler
	if searchTName != "" {
		searchTmpl, err = tmplfn.Assemble(tmplfn.Join(tmplfn.TimeMap, tmplfn.PageMap), searchTName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s error, %s\n", searchTName, err)
			os.Exit(1)
		}
	} else {
		if tSrc, ok := dataset.SiteDefaults["/templates/search.tmpl"]; ok == true {
			searchTmpl, err = template.New("search.tmpl").Funcs(tmplfn.Join(tmplfn.TimeMap, tmplfn.PageMap)).Parse(string(tSrc))
			if err != nil {
				fmt.Fprintf(os.Stderr, "default search template error, %s\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "can't process default template\n")
			os.Exit(1)
		}
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
		default:
			src, _ := json.MarshalIndent(results, " ", " ")
			log.Printf("DEBUG src:\n%s\n", src)
			if err := dataset.HTMLFormatter(w, results, searchTmpl); err != nil {
				http.Error(w, fmt.Sprintf("%s", err), 500)
			}
		}
	}

	// Define our search API prefix path
	http.HandleFunc("/api", searchHandler)

	// FIXME: If DocRoot defined we want to use it, otherwise we need to fallback to what is supplied in dataset.SiteDefaults

	http.Handle("/", http.FileServer(http.Dir(docRoot)))
	if u.Scheme == "https" {
		err := http.ListenAndServeTLS(u.Host, sslCert, sslKey, mkpage.RequestLogger(mkpage.StaticRouter(http.DefaultServeMux)))
		if err != nil {
			log.Fatalf("%s", err)
		}
	} else {
		err := http.ListenAndServe(u.Host, mkpage.RequestLogger(mkpage.StaticRouter(http.DefaultServeMux)))
		if err != nil {
			log.Fatalf("%s", err)
		}
	}
}
