
# USAGE

```
    dsws [OPTIONS] [KEY_VALUE_PAIRS] [DOC_ROOT] BLEVE_INDEXES
```

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
	-acme	Enable Let's Encypt ACME TLS support
	-c	Set the path for the SSL Cert
	-cert	Set the path for the SSL Cert
	-dev-mode	reload templates on each page request
	-example	display example(s)
	-h	display help
	-help	display help
	-indexes	comma or colon delimited list of index names
	-k	Set the path for the SSL Key
	-key	Set the path for the SSL Key
	-l	display license
	-license	display license
	-show-templates	display the source code of the template(s)
	-t	the path to the search result template(s) (colon delimited)
	-template	the path to the search result template(s) (colon delimited)
	-u	The protocal and hostname listen for as a URL
	-url	The protocal and hostname listen for as a URL
	-v	display version
	-version	display version
```


## EXAMPLES

Run web server using the content in the current directory
(assumes the environment variables DATASET_DOCROOT are not defined).

```
   dsws
```

Run web service using "index.bleve" index, results templates in 
"templates/search.tmpl" and a "htdocs" directory for static files.

```
   dsws -template=templates/search.tmpl htdocs index.bleve
```

Run a web service with custom navigation taken from a Markdown file

```
   dsws -template=templates/search.tmpl "Nav=nav.md" index.bleve
```

Running above web service using ACME TLS support (i.e. Let's Encrypt).
Note will only include the hostname as the ACME setup is for
listenning on port 443. This may require privilaged account
and will require that the hostname listed matches the public
DNS for the machine (this is need by the ACME protocol to
issue the cert, see https://letsencrypt.org for details)

```
   dsws -acme -template=templates/search.tmpl "Nav=nav.md" index.bleve
```

