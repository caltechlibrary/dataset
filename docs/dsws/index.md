
# USAGE

	dsws [OPTIONS]

## SYNOPSIS


## Description

dsws is a web server and provides a web search service for indexes 
built from a dataset collection.

### CONFIGURATION

dsws can be configurated through environment settings. The following are
supported.

+ DATASET_URL  - (optional) sets the URL to listen on (e.g. http://localhost:8011)
+ DATASET_SSL_KEY - (optional) the path to the SSL key if using https
+ DATASET_SSL_CERT - (optional) the path to the SSL cert if using https
+ DATASET_TEMPLATE - (optional) path to search results template(s)



## OPTIONS

```
    -acme                     Enable Let's Encypt ACME TLS support
    -c, -cert                 Set the path for the SSL Cert
    -cors-origin              Set the restriction for CORS origin headers
    -dev-mode                 reload templates on each page request
    -e, -examples             display examples
    -generate-markdown-docs   output documentation in Markdown
    -h, -help                 display help
    -i, -input                input file name
    -indexes                  comma or colon delimited list of index names
    -k, -key                  Set the path for the SSL Key
    -l, -license              display license
    -nl, -newline             if set to false suppress the trailing newline
    -o, -output               output file name
    -p, -pretty               pretty print output
    -quiet                    suppress error messages
    -show-templates           display the source code of the template(s)
    -t, -template             the path to the search result template(s) (colon delimited)
    -u, -url                  The protocol and hostname listen for as a URL
    -v, -version              display version
```


dsws v0.0.27-dev
