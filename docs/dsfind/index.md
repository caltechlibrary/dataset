
# USAGE

	dsfind [OPTIONS]

## SYNOPSIS


## Description

dsfind is a command line tool for querying a Bleve indexes based on the records in a 
dataset collection. By default dsfind is assumed there is an index named after the 
collection. An option lets you choose different indexes to query. Results are 
written to standard out and are paged. The query syntax supported is described
at http://www.blevesearch.com/docs/Query-String-Query/.

Options can be used to modify the type of indexes queried as well as how results
are output.



## OPTIONS

```
    -csv                      format results as a CSV document, used with fields option
    -csv-skip-header          don't output a header row, only values for csv output
    -e, -examples             display examples
    -explain                  explain results in a verbose JSON document
    -fields                   comma delimited list of fields to display in the results
    -from                     return the result starting with this result number
    -generate-markdown-docs   output documentation in Markdown
    -h, -help                 display help
    -highlight                display highlight in search results
    -highlighter              set the highlighter (ansi,html) for search results
    -i, -input                input file name
    -ids                      output only a list of ids from results
    -indexes                  colon or comma delimited list of index names
    -json                     format results as a JSON document
    -l, -license              display license
    -nl, -newline             if set to false suppress the trailing newline
    -o, -output               output file name
    -p, -pretty               pretty print output
    -quiet                    suppress error messages
    -sample                   return a sample of size N of results
    -size                     number of results returned for request
    -sort                     a comma delimited list of field names to sort by
    -v, -version              display version
```


dsfind v0.0.21-dev
