package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	// CaltechLibrary Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/dataset"
)

var (
	usage = `USAGE: %s [OPTIONS] SEARCH_STRINGS`

	description = `
SYNOPSIS

%s is a command line tool for querying a Bleve indexes based on the records in a 
dataset collection. By default %s is assumed there is an index named after the 
collection. An option lets you choose different indexes to query. Results are 
written to standard out and are paged. The query syntax supported is described
at http://www.blevesearch.com/docs/Query-String-Query/.

Options can be used to modify the type of indexes queried as well as how results
are output.`

	examples = `
EXAMPLES

In the example the index will be created for a collection called "characters".

    %s -c characters "Jack Flanders"

This would search the Bleve index named characters.bleve for the string "Jack Flanders" 
returning records that matched based on how the index was defined.`

	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool

	// App Specific Options
	collectionName string
	indexNames     string
	showHighlight  bool
	resultFields   string
	sortBy         string
	jsonFormat     bool
	csvFormat      bool
	idsOnly        bool
	size           string
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
	flag.StringVar(&collectionName, "c", "", "sets the collection to be used")
	flag.StringVar(&collectionName, "collection", "", "sets the collection to be used")
	flag.StringVar(&indexNames, "indexes", "", "a colon delimited list of index names")
	flag.StringVar(&sortBy, "sort", "", "a colon delimited list of field names to sort by")
	flag.BoolVar(&showHighlight, "highlight", false, "display highlight in search results")
	flag.StringVar(&resultFields, "fields", "", "colon delimited list of fields to display in the results")
	flag.BoolVar(&jsonFormat, "json", false, "format results as a JSON document")
	flag.BoolVar(&csvFormat, "csv", false, "format results as a CSV document, used with fields option")
	flag.BoolVar(&idsOnly, "ids", false, "output only a list of ids from results")
	flag.StringVar(&size, "size", "", "number of results returned per request or the word \"all\"")
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

	// Merge environment
	datasetEnv := os.Getenv("DATASET")
	if datasetEnv != "" && collectionName == "" {
		collectionName = datasetEnv
	}
	if len(indexNames) == 0 {
		indexNames = fmt.Sprintf("%s.bleve", collectionName)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println(cfg.Usage())
		os.Exit(1)
	}
	options := map[string]string{}
	if explain != "" {
		options["explain"] = "true"
		jsonFormat = true
	}

	if sortBy != "" {
		options["sort_by"] = sortBy
	}
	if showHighlight == true {
		options["highlight"] = "true"
		options["highlighter"] = "ansi"
	}

	if resultFields != "" {
		options["result_fields"] = strings.TrimSpace(resultFields)
	} else {
		//options["result_fields"] = "*"
	}
	if size == "all" {
		options["size"] = "1000000000000"
	} else if size != "" {
		options["size"] = size
	}
	if from != 0 {
		options["from"] = fmt.Sprintf("%d", from)
	}

	results, err := dataset.Find(os.Stdout, strings.Split(indexNames, ":"), args, options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't search index %s, %s\n", indexNames, err)
		os.Exit(1)
	}

	//
	// Handle results formatting choices
	//

	if jsonFormat == true {
		src, err := json.Marshal(results)
		if err != nil {
			fmt.Fprintf(os.Stderr, "JSON conversion error, %s\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%s\n", src)
		os.Exit(0)
	}

	if csvFormat == true {
		if resultFields == "" {
			fmt.Fprintf(os.Stderr, "-csv must be used with -fields option\n")
			os.Exit(1)
		}
		// Note: we need to provide the fieldnames that will be come columns
		w := csv.NewWriter(os.Stdout)

		// write a header row
		colNames := strings.Split(resultFields, ":")
		if err := w.Write(colNames); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		for _, hit := range results.Hits {
			row := []string{}
			for _, col := range colNames {
				if val, ok := hit.Fields[col]; ok == true {
					switch val := val.(type) {
					case int:
						row = append(row, strconv.FormatInt(int64(val), 10))
					case uint:
						row = append(row, strconv.FormatUint(uint64(val), 10))
					case int64:
						row = append(row, strconv.FormatInt(val, 10))
					case uint64:
						row = append(row, strconv.FormatUint(val, 10))
					case float64:
						row = append(row, strconv.FormatFloat(val, 'G', -1, 64))
					case string:
						row = append(row, strings.TrimSpace(val))
					default:
						row = append(row, strings.TrimSpace(fmt.Sprintf("%s", val)))
					}
				} else {
					row = append(row, "")
				}
			}
			if err := w.Write(row); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
		}
		w.Flush()
		if err := w.Error(); err != nil {
			fmt.Fprint(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if idsOnly == true {
		for _, hit := range results.Hits {
			fmt.Fprintf(os.Stdout, "%s\n", hit.ID)
		}
		os.Exit(0)
	}

	fmt.Fprintf(os.Stdout, "%s\n", results)
}
