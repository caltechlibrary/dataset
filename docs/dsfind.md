
# dsfind

## USAGE

    dsfind [OPTIONS] SEARCH_STRINGS

## SYNOPSIS

_dsfind_ is a command line tool for querying a Bleve indexes based on the records in a 
dataset collection. By default _dsfind_ is assumed there is an index named after the 
collection. An option lets you choose different indexes to query. Results are 
written to standard out and are paged. The query syntax supported is described
at http://www.blevesearch.com/docs/Query-String-Query/.

Options can be used to modify the type of indexes queried as well as how results
are output.

## OPTIONS

```
	-c	sets the collection to be used
	-collection	sets the collection to be used
	-fields	colon delimited list of fields to display in the results, defaults to *
	-h	display help
	-help	display help
	-highlight	display highlight in search results
	-indexes	a colon delimited list of index names
	-l	display license
	-license	display license
	-sort	a colon delimited list of field names to sort by
	-v	display version
	-version	display version
```

## EXAMPLES

In the example the index will be created for a collection called "characters".

```shell
    dsfind -c characters "Jack Flanders"
```

This would search the Bleve index named characters.bleve for the string "Jack Flanders" 
returning records that matched based on how the index was defined.


dsfind v0.0.1-beta11
