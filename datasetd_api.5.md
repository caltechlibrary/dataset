%datasetd(1) user manual | version 2.3.1 656569c
% R. S. Doiel
% 2025-07-10


# datasetd REST API

datasetd provides a RESTful JSON API for working with a dataset collection. This document describes the path expressions and to interact with the API.  Note some of the methods and paths require permissions as set in the datasetd YAML or JSON [settings file](datasetd_yaml.5.md).

## basic path expressions

There are three basic forms of the URL paths supported by the API.

- `/api/<COLLECTION_NAME>/keys`, get a list of all keys in the the collection
- `/api/<COLLECTION_NAME>/object/<OPTIONS>`, interact with an object in the collection (e.g. create, read, update, delete)
- `/api/<COLLECTION_NAME>/query/<QUERY_NAME>/<FIELDS>`, query the collection and receive a list of objects in response

The "`<COLLECTION_NAME>`" would be the name of the dataset collection, e.g. "mydata.ds".

The "`<OPTIONS>`" holds any additional parameters related to the verb. Options are separated by the path delimiter (i.e. "/"). The options are optional. They do not require a trailing slash.

The "`<QUERY_NAME>`" is the query name defined in the YAML configuration for the specific collection.

The "`<FIELDS>`" holds the set of fields being passed into the query. These are delimited with the path separator like with options (i.e. "/"). Fields are optional and they do not require a trailing slash.

## HTTP Methods

The datasetd REST API follows the rest practices. Good examples are POST creates, GET reads, PUT updates, and DELETE removes. It is important to remember that the HTTP method and path expression need to match form the actions you'd take using the command line version of dataset. For example to create a new object you'd use the object path without any options and a POST expression. You can do a read of an object using the GET method along withe object style path.

## Content Type and the API

The REST API works with JSON data. The service does not support multipart urlencoded content. You MUST use the content type of `application/json` when performing a POST, or PUT. This means if you are building a user interface for a collections datasetd service you need to appropriately use JavaScript to send content into the API and set the content type to `application/json`.

## Examples

Here's an example of a list, in YAML, of people in a collection called "people.ds". There are some fields for the name, sorted name, display name and orcid. The pid is the "key" used to store the objects in our collection.

~~~yaml
people:
  - pid: doe-jane
    family: Doe
    lived: Jane
    orcid: 9999-9999-9999-9999
~~~

In JSON this would look like

~~~json
{
  "people": [
    {
      "pid": "doe-jane",
      "family": "Doe",
      "lived": "Jane",
      "orcid": "9999-9999-9999-9999"
    }
  ]
}
~~~

### create

The create action is formed with the object URL path, the POST http method and the content type of "application/json". It POST data is expressed as a JSON object.

The object path includes the dataset key you'll assign in the collection. The key must be unique and not currently exist in the collection.

If we're adding an object with the key of "doe-jane" to our collection called "people.ds" then the object URL path would be  `/api/people.ds/object/doe-jane`. NOTE: the object key is included as a single parameter after "object" path element.

Adding an object to our collection using curl looks like the following.

~~~shell
curl -X POST \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -d '{"pid": "doe-jane", "family": "Doe", "lived": "Jane", "orcid": "9999-9999-9999-9999" }' \
  http://localhost:8485/api/people.ds/object/doe-jane  
~~~

### read

The read action is formed with the object URL path, the GET http method and the content type of "application/json".  There is no data
aside from the URL to request the object. Here's what it would look like using curl to access the API.

~~~shell
curl http://localhost:8485/api/people.ds/object/doe-jane  
~~~

### update

Like create update is formed from the object URL path, content type of "application/json" the data is expressed as a JSON object.
Onlike create we use the PUT http method.

Here's how you would use curl to get the JSON expression of the object called "doe-jane" in your collection.

~~~shell
curl -X PUT \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -d '{"pid": "doe-jane", "family": "Doe", "lived": "Jane", "orcid": "9999-9999-9999-9999" }' \
  http://localhost:8485/api/people.ds/object/doe-jane  
~~~

This will overwrite the existing "doe-jane". NOTE the record must exist or you will get an error.

### delete

If you want to delete the "doe-jane" record in "people.ds" you perform an http DELETE method and form the url like a read.

~~~shell
curl -X DELETE http://localhost:8485/api/people.ds/object/doe-jane  
~~~

## query

The query path lets you run a predefined query from your settings YAML file. The http method used is a POST. This is becaue we need to send data inorder to receive a response. The resulting data is expressed as a JSON array of object. Like with create, read, update and delete you use the content type of "application/json".

In the settings file the queries are named. The query names are unique. One or many queries may be defined. The SQL expression associated with the name run as a prepared statement and parameters are mapped into based on the URL path provided. This allows you use many fields in forming your query.

Let's say we have a query called "full_name". It is defined to run the following SQL.

~~~sql
select src
from people
where src->>'family' like ?
  and src->>'lived' like ?
order by family, lived
~~~

NOTE: The SQL is has to retain the constraint of a single object per row, normally this will be "src" for dataset collections.

When you form a query path we need to indicate that the parameter for family and lived names need to get mapped to their respect positional references in the SQL. This is done as following url path. In this example "full_name" is the name of the query while "family" and "lived" are the values mapped into the parameters.

~~~
/api/people.ds/query/full_name/family/lived
~~~

The web form could look like this.  

~~~
<form id="query_name">
   <label for="family">Family</label> <input id="family" name="family" ><br/>
   <label for="lived">Lived</label> <input id="lived" name="lived" ><br/>
   <button type="submit">Search</button>
</form>
~~~

REMEMBER: the JSON API only supports the content type of "application/json" so you can use the browser's action and method in the form.

You would include JavaScript in the your HTML to pull the values out of the form and create a JSON object. If I searched
for someone who had the family name "Doe" and he lived name of "Jane" the object submitted to query might look like the following. 

~~~json
{
    "family": "Doe"
    "lived": "Jane"
}
~~~

The curl expression would look like the following simulating the form submission would look like the following.


~~~shell
curl -X POST \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -d '{"family": "Doe", "lived": "Jane" }' \
  http://localhost:8485/api/people.ds/query/full_name/family/lived
~~~



