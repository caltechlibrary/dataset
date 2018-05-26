
# Cloud Storage

These days it is often useful to leverage cloud storage. _dataset_ currently supports Amazon's 
cloud storage and Google's cloud storage sollutions. These can be configured through the
setting the DATASET environment variable or via the command line. The buckets for either Amazon
or Google need to have already been setup including authenticated configs (e.g. with the aws or gsutil
tools respectively).  The aws tool is available
[here](https://aws.amazon.com/cli), and can be set up using `aws configure` and
entering the Access key information from your AWS user accounts page (under
"Security credentials".  If the prefix for the path to the collection is prefixed with s3:// then the
collection is stored at AWS S3, if the prefix is gs:// then it is stored on Google Cloud Storage and
if there is now prefix it is stored on local disc.

## Local Storage setup

```shell
    #!/bin/bash
    
    #
    # Local Disc setup
    #
    export DATASET="my-test-bucket"
```

## S3 Storage setup

```shell
    #!/bin/bash
    
    #
    # S3 test setup example
    #
    
    # Load the config and credentials in ~/.aws if found
    export AWS_SDK_LOAD_CONFIG=1
    # You will need to define your bucket name
    export DATASET="s3://my-test-bucket"
```

## GS Storage setup

```shell
    #!/bin/bash
    
    #
    # Google Cloud Storage test setup example
    #
    export DATASET="gs://my-test-bucket"
```
