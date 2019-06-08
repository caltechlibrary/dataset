
# USAGE

	dataset [OPTIONS] VERB COLLECTION_NAME [ACTION PARAMETERS...]

## SYNOPSIS


dataset is a command line tool demonstrating dataset package for 
managing JSON documents stored on disc. A dataset is organized 
around collections, collections contain a pairtree holding specific 
JSON documents and related content.  In addition to the JSON 
documents dataset maintains metadata for management of the 
documents, their attachments as well as a ability to generate 
select lists based JSON document keys (aka JSON document names).



## OPTIONS

Options are shared between all actions and must precede the action 
on the command line.

```
    -c, -collection           sets the collection to be used
    -client-secret            (import-gsheet, export-gsheet) set the client secret path and filename for GSheet access
    -csv                      (find) format results as a CSV document, used with fields option
    -csv-skip-header          (find) don't output a header row, only values for csv output
    -e, -examples             display examples
    -explain                  (find) explain results in a verbose JSON document
    -fields                   (find) comma delimited list of fields to display in the results
    -from                     (find) return the result starting with this result number
    -generate-markdown-docs   output documentation in Markdown
    -h, -help                 display help
    -highlight                (find) display highlight in search results
    -highlighter              (find) set the highlighter (ansi,html) for search results
    -i, -input                input file name
    -ids, -ids-only           (find) output only a list of ids from results
    -json                     (find) format results as a JSON document
    -key-file                 operate on the record keys contained in file, one key per line
    -l, -license              display license
    -layout                   set file layout for a new collection (i.e. "buckets" or "pairtree")
    -nl, -newline             if set to false suppress the trailing newline
    -o, -output               output file name
    -overwrite                overwrite will treat a create as update if the record exists
    -p, -pretty               pretty print output
    -quiet                    suppress error messages
    -sample                   set the sample size when listing keys
    -sort                     (find) a comma delimited list of field names to sort by
    -use-header-row           (import) use the header row as attribute names in the JSON document
    -v, -version              display version
    -verbose                  output rows processed on importing from CSV
```


## VERBS

```
    attach         Attach a document (file) to a JSON record in a collection
    attachments    List of attachments associated with a JSON record in a collection
    check          Check the health of a dataset collection
    clone          Clone a collection from a list of keys into a new collection
    clone-sample   Clone a collection into a sample size based training collection and test collection
    count          Counts the number of records in a collection, accepts a filter for sub-counts
    create         Create a JSON record in a collection
    delete         Delete a JSON record (and attachments) from a collection
    delete-frame   remove a frame from a collection
    detach         Copy an attach out of an associated JSON record in a collection
    export-csv     Export a JSON records from a collection to a CSV file
    export-gsheet  Export a collection's JSON records to a GSheet
    frame          define or retrieve a frame from a collection
    frame-labels   define explicitly the labels associated with a frame
    frame-types    define explicitly the column type names associated with a frame
    frames         list the available frames in a collection
    grid           Creates a data grid from a list keys of dot paths
    haskey         Returns true if key is in collection, false otherwise
    import-csv     Import a CSV file's rows as JSON records into a collection
    import-gsheet  Import a GSheet rows as JSON records into a collection
    init           Initialize a dataset collection
    join           Join a JSON record with a new JSON object in a collection
    keys           List the keys in a collection, support filtering and sorting
    list           List the JSON records as an array for provided record ids
    path           Show the file system path to a JSON record in a collection
    prune          Remove attachments from a JSON record in a collection
    read           Read back a JSON record from a collection
    reframe        re-generate a frame with existing or provided key list
    repair         Try to repair a damaged dataset collection
    status         Checks for a collection's layout and version
    update         Update a JSON record in a collection
```


Related: [attach](attach.html), [attachments](attachments.html), [check](check.html), [clone](clone.html), [clone-sample](clone-sample.html), [count](count.html), [create](create.html), [delete](delete.html), [delete-frame](delete-frame.html), [detach](detach.html), [export-csv](export-csv.html), [export-gsheet](export-gsheet.html), [frame](frame.html), [frame-labels](frame-labels.html), [frame-types](frame-types.html), [frames](frames.html), [grid](grid.html), [haskey](haskey.html), [import-csv](import-csv.html), [import-gsheet](import-gsheet.html),  [init](init.html), [join](join.html), [keys](keys.html), [list](list.html), [path](path.html), [prune](prune.html), [read](read.html), [reframe](reframe.html), [repair](repair.html), [status](status.html), [update](update.html)

dataset v0.0.45
