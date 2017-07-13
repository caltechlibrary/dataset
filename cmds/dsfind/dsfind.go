package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	// CaltechLibrary Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/dataset"
)

var (
	usage = `USAGE: %s [OPTIONS] [INDEX_LIST] SEARCH_STRINGS`

	description = `
SYNOPSIS

%s is a command line tool for querying a Bleve indexes based on the records in a 
dataset collection. By default %s is assumed there is an index named after the 
collection. An option lets you choose different indexes to query. Results are 
written to standard out and are paged. The query syntax supported is described
at http://www.blevesearch.com/docs/Query-String-Query/.

Options can be used to modify the type of indexes queried as well as how results
are output.
`

	examples = `
EXAMPLES

In the example the index will be created for a collection called "characters".

    %s characters.bleve "Jack Flanders"

This would search the Bleve index named characters.bleve for the string "Jack Flanders" 
returning records that matched based on how the index was defined.
`

	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool

	// App Specific Options
	indexList      string
	showHighlight  bool
	setHighlighter string
	resultFields   string
	sortBy         string
	jsonFormat     bool
	csvFormat      bool
	idsOnly        bool
	size           int
	from           int
	explain        string // Note: will be converted to boolean so expecting 1,0,T,F,true,false, etc.
)

func init() {
	// Standard Options
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")

	// Application Options
	flag.StringVar(&indexList, "indexes", "", "colon or comma delimited list of index names")
	flag.StringVar(&sortBy, "sort", "", "a comma delimited list of field names to sort by")
	flag.BoolVar(&showHighlight, "highlight", false, "display highlight in search results")
	flag.StringVar(&setHighlighter, "highlighter", "", "set the highlighter (ansi,html) for search results")
	flag.StringVar(&resultFields, "fields", "", "comma delimited list of fields to display in the results")
	flag.BoolVar(&jsonFormat, "json", false, "format results as a JSON document")
	flag.BoolVar(&csvFormat, "csv", false, "format results as a CSV document, used with fields option")
	flag.BoolVar(&idsOnly, "ids", false, "output only a list of ids from results")
	flag.IntVar(&size, "size", 0, "number of results returned for request")
	flag.IntVar(&from, "from", 0, "return the result starting with this result number")
	flag.StringVar(&explain, "explain", "", "explain results in a verbose JSON document")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()

	cfg := cli.New(appName, appName, fmt.Sprintf(dataset.License, appName, dataset.Version), dataset.Version)
	cfg.UsageText = fmt.Sprintf(usage, appName)
	cfg.DescriptionText = fmt.Sprintf(description, appName, appName)
	cfg.ExampleText = fmt.Sprintf(examples, appName)

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

	// Handle the case where indexes were listed with the -indexes option like dsfind
	var indexNames []string
	if indexList != "" {
		var delimiter = ","
		if strings.Contains(indexList, ":") {
			delimiter = ":"
		}
		indexNames = strings.Split(indexList, delimiter)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println(cfg.Usage())
		os.Exit(1)
	}

	// Collect any additional index names from the remaining args
	for _, arg := range args {
		if path.Ext(arg) == ".bleve" {
			indexNames = append(indexNames, arg)
		}
	}
	if len(indexNames) == 0 {
		fmt.Printf("Do not know what index to use")
		os.Exit(1)
	}

	options := map[string]string{}
	if explain != "" {
		options["explain"] = "true"
		jsonFormat = true
	}

	if from != 0 {
		options["from"] = fmt.Sprintf("%d", from)
	}
	if size > 0 {
		options["size"] = fmt.Sprintf("%d", size)
	}
	if sortBy != "" {
		options["sort"] = sortBy
	}
	if showHighlight == true {
		options["highlight"] = "true"
		options["highlighter"] = "ansi"
	}
	if setHighlighter != "" {
		options["highlight"] = "true"
		options["highlighter"] = setHighlighter
	}

	if resultFields != "" {
		options["fields"] = strings.TrimSpace(resultFields)
	} else {
		options["fields"] = "*"
	}

	idxAlias, idxFields, err := dataset.OpenIndexes(indexNames)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't open index %s, %s\n", strings.Join(indexNames, ", "), err)
		os.Exit(1)
	}
	defer idxAlias.Close()

	results, err := dataset.Find(os.Stdout, idxAlias, args, options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't search index %s, %s\n", strings.Join(indexNames, ", "), err)
		os.Exit(1)
	}

	//
	// Handle results formatting choices
	//
	switch {
	case jsonFormat == true:
		if err := dataset.JSONFormatter(os.Stdout, results); err != nil {
			fmt.Fprintf(os.Stderr, "JSON formatting error, %s\n", err)
			os.Exit(1)
		}
	case csvFormat == true:
		var fields []string
		if resultFields == "" {
			fields = idxFields
		} else {
			fields = strings.Split(resultFields, ",")
		}
		if err := dataset.CSVFormatter(os.Stdout, results, fields); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	case idsOnly == true:
		for _, hit := range results.Hits {
			fmt.Fprintf(os.Stdout, "%s\n", hit.ID)
		}
	default:
		fmt.Fprintf(os.Stdout, "%s\n", results)
	}
}
