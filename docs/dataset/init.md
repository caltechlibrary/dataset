
# init

## Syntax

```
    dataset init COLLECT_NAME
```

## Description

init create a collection. Collections can be create on local disc, on Amazon's S3 and
in Google Cloud Storage. If you are initializing a collection on S3 or in Google
Cloud Storage then the bucket (where the collection will reside) needs to already 
exist and you need to have been authenticated.

To store your collection in S3 prefix the path with s3://, likewise for Google
Cloud Storage the prefix is gs://.

## Usage

The following three example commands create a dataset collection named "data".
First one local disc in the current currectory, then in S3 and again in Google
Cloud Storage. In the case of S3 and Google Cloud Storage the buckets exist and
are named "stuff.example.org". Also for both remote storage options it is assumed
you've authenticated and have your environment setup correctly.

```
    dataset init data
    dataset init s3://stuff.example.org/data
    dataset init gs://stuff.example.org/data
```

NOTE: After each envocation of `dataset init` if all went well you will be
shown an `export DATASET=...` line. This cut and pasted into your shell
to save using the `-c` option explicitly indicating which dataset to operate on.

### S3 environment example

You can refernce loading the environment for AWS S3 access previous setup with
the AWS SDK tool with by exporting the "AWS_SDK_LOAD_CONFIG" environment variable
with the a value of "1".

```shell
    export AWS_SDK_LOAD_CONFIG=1
```

### Google Cloud Platform

Google Cloud Platform authentication can be done via the _gsutil_ command available
with Google Cloud SDK setup.

