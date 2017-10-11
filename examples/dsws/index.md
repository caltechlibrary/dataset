
## EXAMPLES

Run web server using the content in the current directory
(assumes the environment variables DATASET_DOCROOT are not defined).

   dsws

Run web service using "index.bleve" index, results templates in 
"templates/search.tmpl" and a "htdocs" directory for static files.

   dsws -template=templates/search.tmpl htdocs index.bleve

Run a web service with custom navigation taken from a Markdown file

   dsws -template=templates/search.tmpl "Nav=nav.md" index.bleve

Running above web service using ACME TLS support (i.e. Let's Encrypt).
Note will only include the hostname as the ACME setup is for
listenning on port 443. This may require privilaged account
and will require that the hostname listed matches the public
DNS for the machine (this is need by the ACME protocol to
issue the cert, see https://letsencrypt.org for details)

   dsws -acme -template=templates/search.tmpl "Nav=nav.md" index.bleve

