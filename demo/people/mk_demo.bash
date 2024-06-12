#!/bin/bash

#
# Create a people demo using an SQLite3 database
# and a simple collection with two people records
#

if [ -d "people.ds" ]; then
	rm -fR people.ds
fi
dataset init people.ds "sqlite://collection.db"
dataset create people.ds doe-jane <<EOT
{
	"family": "Doe",
	"lived": "Jane",
	"orcid": "0000-0000-0000-0000"
}
EOT

dataset create people.ds doe-jim <<EOT
{
	"family": "Doe",
	"lived": "Jimmy",
	"orcid": "9999-9999-9999-9999"
}
EOT

cat <<EOT >people.yaml
#
# people.yaml descrives a datasetd web service demonstrating
# how the query functionality works.
#
host: localhost:8484
collections: 
  - dataset: people.ds
    query:
      list: |
         select src
         from people
      full_name: |
        select src
        from people
        where src->>'family' like ?
           and src->>'lived' like ? 
      servertime: |
        select datetime()
    read: true
    keys: true
EOT

cat <<EOT >run_demo.bash
echo
echo "The following uses curl to demonstrate the query API provided by datasetd"
echo "Content type is set to application/json and the http method is POST"
echo
echo "List all the object using the query path"
echo "     url: http://localhost:8484/api/people.ds/query/list"
echo
curl -X POST \
     -H 'Content-type: application/json' \
	 -H 'Accepted: application/json' \
	 http://localhost:8484/api/people.ds/query/list
echo
echo
echo "List people named Jane Doe"
echo "     url: http://localhost:8484/api/people.ds/query/full_name/family/lived"
echo
curl -X POST \
     -d '{"family": "Doe", "lived": "Jane"}' \
     -H 'Content-type: application/json' \
	 -H 'Accepted: application/json' \
	 http://localhost:8484/api/people.ds/query/full_name/family/lived
echo
echo "List the server datetime"
echo "     url: http://localhost:8484/api/people.ds/query/servertime"
echo
curl -X POST \
     -H 'Content-type: application/json' \
	 -H 'Accepted: application/json' \
	 http://localhost:8484/api/people.ds/query/servertime

echo
echo
EOT

echo
chmod 775 run_demo.bash
cat <<EOT

  Start the web API with 

     datasetd people.yaml

  Then run the demo of the API with

     ./run_demo.bash

EOT
