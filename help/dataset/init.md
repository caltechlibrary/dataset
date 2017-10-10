+ init - initialize a new collection if none exists, requires a path to collection
  + once collection is created, set the environment variable DATASET
    to collection name
  + if you're using S3 for storing your dataset prefix your path with 's3://'
    'dataset init s3://mybucket/mydataset-collections'
  + if you're using GS (Google Cloud Storage) prefix your path with 'gs://'
    'dataset init gs://mybucket/mydataset-collections'
