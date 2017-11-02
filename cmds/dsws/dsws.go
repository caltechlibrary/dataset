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
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/dataset"
	"github.com/caltechlibrary/mkpage"
	"github.com/caltechlibrary/tmplfn"

	// Other packages
	"golang.org/x/crypto/acme/autocert"
)

// Flag options
var (
	// Standard options
	showHelp     bool
	showVersion  bool
	showLicense  bool
	showExamples bool

	// local app options
	uri           string
	sslKey        string
	sslCert       string
	searchTName   string
	devMode       bool
	showTemplates bool
	indexList     string
	letsEncrypt   bool

	// Provided as an ordered command line arg
	docRoot    string
	indexNames []string
)

// trimmedSplit splits a string on commas and run performs a TrimSpace on the resulting array elements
func trimmedSplit(s, delimiter string) []string {
	r := strings.Split(s, delimiter)
	for i, val := range r {
		r[i] = strings.TrimSpace(val)
	}
	return r
}

// redirectToApi will redirect to the /api search result page
func redirectToApi(w http.ResponseWriter, r *http.Request) {
	target := "/api/"
	mkpage.ResponseLogger(r, http.StatusTemporaryRedirect, fmt.Errorf("redirected %s to %s", r.URL.Path, target))
	http.Redirect(w, r, target, http.StatusTemporaryRedirect)
}

func init() {
	defaultURL := "http://localhost:8011"

	// Standard Options
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showExamples, "example", false, "display example(s)")

	// App Options
	flag.StringVar(&uri, "u", defaultURL, "The protocol and hostname listen for as a URL")
	flag.StringVar(&uri, "url", defaultURL, "The protocol and hostname listen for as a URL")
	flag.StringVar(&sslKey, "k", "", "Set the path for the SSL Key")
	flag.StringVar(&sslKey, "key", "", "Set the path for the SSL Key")
	flag.StringVar(&sslCert, "c", "", "Set the path for the SSL Cert")
	flag.StringVar(&sslCert, "cert", "", "Set the path for the SSL Cert")
	flag.StringVar(&searchTName, "template", "", "the path to the search result template(s) (colon delimited)")
	flag.StringVar(&searchTName, "t", "", "the path to the search result template(s) (colon delimited)")
	flag.BoolVar(&showTemplates, "show-templates", false, "display the source code of the template(s)")
	flag.BoolVar(&devMode, "dev-mode", false, "reload templates on each page request")
	flag.StringVar(&indexList, "indexes", "", "comma or colon delimited list of index names")
	flag.BoolVar(&letsEncrypt, "acme", false, "Enable Let's Encypt ACME TLS support")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()
	args := flag.Args()

	// Configuration and command line interation
	cfg := cli.New(appName, "DATASET", dataset.Version)
	cfg.LicenseText = fmt.Sprintf(dataset.License, appName, dataset.Version)
	cfg.UsageText = fmt.Sprintf("%s", Help["usage"])
	cfg.DescriptionText = fmt.Sprintf("%s", Help["description"])
	cfg.OptionText = "## OPTIONS\n\n"
	cfg.ExampleText = fmt.Sprintf("%s", Examples["index"])

	// Add help and examples
	for k, v := range Help {
		if k != "nav" {
			cfg.AddHelp(k, fmt.Sprintf("%s", v))
		}
	}
	for k, v := range Examples {
		if k != "nav" {
			cfg.AddExample(k, fmt.Sprintf("%s", v))
		}
	}

	// Process flags and update the environment as needed.
	if showHelp == true {
		if len(args) > 0 {
			fmt.Println(cfg.Help(args...))
		} else {
			fmt.Println(cfg.Usage())
		}
		os.Exit(0)
	}

	if showExamples == true {
		/*
			if len(args) > 0 {
				fmt.Println(cfg.Example(args...))
			} else {
				fmt.Printf("\n%s", cfg.Example())
			}
		*/
		fmt.Println(cfg.ExampleText)
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

	// Load and validate the templates for using in the searchHandler
	searchTName = cfg.CheckOption("template", cfg.MergeEnv("template", searchTName), false)
	templateNames := []string{}
	if searchTName != "" {
		templateNames = strings.Split(searchTName, ":")
	}
	tmpl := tmplfn.New(tmplfn.AllFuncs())

	// Setup templates
	if len(templateNames) > 0 {
		// Load any user supplied templates
		log.Printf("Search templates %q", strings.Join(templateNames, ", "))
		if err := tmpl.ReadFiles(templateNames...); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	} else {
		log.Printf("Using default search templates")
		// Load our default templates from Defaults
		if err := tmpl.ReadMap(Defaults); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}

	if showTemplates == true {
		for name, src := range tmpl.Code {
			if len(tmpl.Code) == 0 {
				fmt.Fprintf(os.Stdout, "%s\n", src)
			} else {
				fmt.Fprintf(os.Stdout,
					"---------------- START: %s -----------------\n",
					name)
				fmt.Fprintf(os.Stdout, "%s\n", src)
				fmt.Fprintf(os.Stdout,
					"---------------- END:   %s -----------------\n",
					name)
			}
		}
		os.Exit(0)
	}

	// Assemble the templates
	searchTmpl, err := tmpl.Assemble()
	if err != nil {
		fmt.Fprintf(os.Stderr, "default search template error, %s\n", err)
		os.Exit(1)
	}

	// Handle the case where indexes were listed with the -indexes option like dsfind
	if indexList != "" {
		var delimiter = ","
		if strings.Contains(indexList, ":") {
			delimiter = ":"
		}
		indexNames = strings.Split(indexList, delimiter)
	}

	// Setup from command line
	pageData := map[string]string{}

	for _, arg := range args {
		ext := path.Ext(arg)
		if ext == ".bleve" {
			indexNames = append(indexNames, arg)
		} else if strings.Contains(arg, "=") {
			kv := strings.SplitN(arg, "=", 2)
			pageData[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		} else {
			docRoot = arg
		}
	}
	if len(indexNames) < 1 {
		fmt.Fprintln(os.Stderr, cfg.UsageText)
		fmt.Fprintf(os.Stderr, "error: one or more Bleve index is required\n")
		os.Exit(1)
	}

	//
	// Final set and start the webservice
	//
	if len(indexNames) == 1 {
		log.Printf("Index %s", strings.Join(indexNames, ", "))
	} else {
		log.Printf("Indexes %s", strings.Join(indexNames, ", "))
	}

	uri = cfg.CheckOption("url", cfg.MergeEnv("url", uri), true)
	u, err := url.Parse(uri)
	if err != nil {
		log.Fatalf("Can't parse %q, %s", uri, err)
	}

	// if the URL starts with http:// then turn off letsEncrypt...
	if u.Scheme == "http" {
		letsEncrypt = false
	}

	if u.Scheme == "https" && letsEncrypt == false {
		sslKey = cfg.CheckOption("ssl_key", cfg.MergeEnv("ssl_key", sslKey), true)
		sslCert = cfg.CheckOption("ssl_cert", cfg.MergeEnv("ssl_cert", sslCert), true)
		log.Printf("SSL Key %s", sslKey)
		log.Printf("SSL Cert %s", sslCert)
	}

	// Open the indexes for reading
	idxAlias, idxFields, err := dataset.OpenIndexes(indexNames)
	if err != nil {
		log.Fatalf("Can't open indexes, %s", err)
	}
	defer idxAlias.Close()

	// Construct our handler
	searchHandler := func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		qformat := values.Get("fmt")
		qString := values.Get("q")
		if qString == "%2A" || qString == "*" {
			http.Error(w, "Missing search terms", 400)
			return
		}
		// Get the options understood by dataset.Find()
		opts := map[string]string{}
		for _, ky := range []string{"size", "from", "ids", "sort", "explain", "fields", "highlight"} {
			if v := values.Get(ky); v != "" {
				// NOTE: we use idxFields for fields' value if no fields or star are passed in
				if ky == "fields" && v == "*" {
					opts[ky] = strings.Join(idxFields, ",")
				} else {
					opts[ky] = v
				}

			}
		}

		//NOTE: If highlight is passed then set the highliter to HTML for web view
		if sVal, ok := opts["highlight"]; ok == true {
			if bVal, err := strconv.ParseBool(sVal); bVal == true && err == nil {
				opts["highlighter"] = "html"
			}
		}
		buf := bytes.NewBufferString("")
		results, err := dataset.Find(buf, idxAlias, []string{qString}, opts)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), 500)
		}
		//FIXME: This is an ugly abuse of a closure to get a developer mode...
		if devMode == true {
			if err := tmpl.ReadFiles(templateNames...); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			} else {
				if t, err := tmpl.Assemble(); err == nil {
					searchTmpl = t
					log.Printf("dev mode: template %s assembled", searchTName)
				} else {
					log.Printf("\n\ndev mode: template %s failed, %s\n\n", strings.Join(templateNames, ", "), err)
				}
			}
		}
		// Based on the request info, format the results appropriately
		var tName string
		switch strings.ToLower(qformat) {
		case "csv":
			fields := trimmedSplit(values.Get("fields"), ",")
			if len(fields) == 0 || (len(fields) == 1 && fields[0] == "*") {
				fields = idxFields
			}
			w.Header().Set("Content-Type", "text/csv")
			if err := dataset.CSVFormatter(w, results, fields, false); err != nil {
				http.Error(w, fmt.Sprintf("%s", err), 500)
			}
			return
		case "json":
			w.Header().Set("Content-Type", "application/json")
			if err := dataset.JSONFormatter(w, results); err != nil {
				http.Error(w, fmt.Sprintf("%s", err), 500)
			}
			return
		case "include":
			w.Header().Set("Content-Type", "text/plain")
			tName = "include.tmpl"
		case "html":
			tName = "page.tmpl"
			w.Header().Set("Content-Type", "text/html")
		default:
			if qformat == "" {
				tName = "page.tmpl"
				w.Header().Set("Content-Type", "text/html")
			} else {
				tName = qformat + ".tmpl"
				//FIXME: Need to pick an appropriate mime type based on format
				//(e.g. BibTeX mime type...)
				w.Header().Set("Content-Type", "text/plain")
			}
		}

		pg := new(bytes.Buffer)
		if err := dataset.Formatter(pg, results, searchTmpl, tName, pageData); err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("Oops, %s formatting error", tName), 500)
		} else {
			pg.WriteTo(w)
		}
	}

	// Define our search API prefix path
	mux := http.NewServeMux()
	//mux.HandleFunc("/api", searchHandler)
	mux.HandleFunc("/api/", searchHandler)
	// FIXME: For each Linux add a /api/INDEXNAME handler

	// Note: If DocRoot is NOT provided we need to redirect to /api
	// instead of using a docRoot with htt.FileServer(http.Dir(docRoot)
	if docRoot == "" {
		log.Printf("Using /api as langing page")
		mux.HandleFunc("/", redirectToApi)
	} else {
		log.Printf("Document root %s", docRoot)
		mux.Handle("/", http.FileServer(http.Dir(docRoot)))
	}

	if letsEncrypt == true {
		// Note: need use a sensible value for data directory
		// this is where cached certificates are stored
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal("Can't determine current working directory where I need to create etc/acme")
		}
		if docRoot == "." || docRoot == cwd {
			log.Fatal("Can't create etc/acme in your shared document root")
		}
		cacheDir := "etc/acme"
		os.MkdirAll(cacheDir, 0700)

		// Setup AMCE TLS Web Server
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(u.Host),
			Cache:      autocert.DirCache(cacheDir),
		}
		sSvr := &http.Server{
			Addr:      ":https",
			TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
			Handler:   mkpage.RequestLogger(mkpage.StaticRouter(mux)),
		}
		go func() {
			log.Printf("Listening for %s (ACME)", u.String())
			log.Fatal(sSvr.ListenAndServeTLS("", ""))
		}()

		// Setup Redirect Web Server
		rmux := http.NewServeMux()
		rmux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			var target string
			if strings.HasPrefix(r.URL.Path, "/") == false {
				target = u.String() + "/" + r.URL.Path
			} else {
				target = u.String() + r.URL.Path
			}
			if len(r.URL.RawQuery) > 0 {
				target += "?" + r.URL.RawQuery
			}
			mkpage.ResponseLogger(r, http.StatusTemporaryRedirect, fmt.Errorf("redirected %s to %s", r.URL.String(), target))
			http.Redirect(w, r, target, http.StatusTemporaryRedirect)
		})
		pSvr := &http.Server{
			Addr:    ":http",
			Handler: mkpage.RequestLogger(rmux),
		}
		log.Printf("Redirecting http://%s to to %s", u.Host, u.String())
		log.Fatal(pSvr.ListenAndServe())
	} else if u.Scheme == "https" {
		log.Printf("Listening for %s", u.String())
		err := http.ListenAndServeTLS(u.Host, sslCert, sslKey, mkpage.RequestLogger(mkpage.StaticRouter(mux)))
		if err != nil {
			log.Fatalf("%s", err)
		}
	} else {
		log.Printf("Listening for %s", u.String())
		err := http.ListenAndServe(u.Host, mkpage.RequestLogger(mkpage.StaticRouter(mux)))
		if err != nil {
			log.Fatalf("%s", err)
		}
	}
}
