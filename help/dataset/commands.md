
# COMMANDS

Documentation on individual commands can be see with `dataset -help COMMAND_NAME` where
"COMMAND_NAME" is replaced with one of the commands below --

+ init - initialize a new collection if none exists, requires a path to collection
  + once collection is created, set the environment variable DATASET
    to collection name
  + if you're using S3 for storing your dataset prefix your path with 's3://'
    'dataset init s3://mybucket/mydataset-collections'
  + if you're using GS (Google Cloud Storage) prefix your path with 'gs://'
    'dataset init gs://mybucket/mydataset-collections'
+ create - creates a new JSON document or replace an existing one in collection
  + requires JSON document name followed by JSON blob or JSON blob read from stdin
+ read - displays a JSON document to stdout
  + requires JSON document name
+ update - updates a JSON document in collection
  + requires JSON document name, followed by replacement JSON document name or 
    JSON document read from stdin
  + JSON document must already exist
+ delete - removes a JSON document from collection
  + requires JSON document name
+ join - brings the functionality of jsonjoin to the dataset command.
  + option update will only add unique key/values not in the existing stored document
  + option overwrite will overwrite all key/values in the existing document
+ filter - takes a filter and returns an unordered list of keys that match filter expression
  + if filter expression not provided as a command line parameter then it is read from stdin
+ keys - returns the keys to stdout, one key per line
+ haskey - returns true is key is in collection, false otherwise
+ path - given a document name return the full path to document
+ attach - attaches a non-JSON content to a JSON record 
  + "dataset attach k1 stats.xlsx" would attach the stats.xlsx file to JSON document named k1
  + (stores content in a related tar file)
+ attachments - lists any attached content for JSON document
  + "dataset attachments k1" would list all the attachments for k1
+ attached - returns attachments for a JSON document 
  + "dataset attached k1" would write out all the attached files for k1
  + "dataset attached k1 stats.xlsx" would write out only the stats.xlsx file attached to k1
+ detach - remove attachments to a JSON document
  + "dataset detach k1 stats.xlsx" would rewrite the attachments tar file without including stats.xlsx
  + "dataset detach k1" would remove ALL attachments to k1
+ import - import a CSV file's rows as JSON documents
  + "dataset import mydata.csv 1" would import the CSV file mydata.csv using column one's value as key
+ export - export a CSV file based on filtered results of collection records rendering dotpaths associated with column names
  + "dataset export titles.csv 'true' '._id,.title,.pubDate' 'id,title,publication date'" 
    this would export all the ids, titles and publication dates as a CSV fiile named titles.csv
+ extract - will return a unique list of unique values based on the associated dot path described in the JSON docs
  + "dataset extract true .authors[:].orcid" would extract a list of authors' orcid ids in collection

