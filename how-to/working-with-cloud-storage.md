
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

# Use _dataset_ with S3

_dataset_ now support integration with S3 storage.  Store _dataset_ content on AWS S3 you should 
download and install the [aws cli sdk](https://aws.amazon.com/cli/), setup your buckets and 
configure permissions, access keys, etc.  _dataset_ will use your local SDK's configuration
(e.g. $HOME/.aws) to configure the connection. You need only set one environment variable, 
run the _dataset_ init option and add the resulting suggested environment variable for working
with your dataset stored at S3.

## Basic steps

1. Set AWS_SDK__LOAD_CONFIG environment variable
2. Envoke the dataaset init command with your "s3://" URL appended with your collectio name
3. Set DATASET environment variable


In the following shell example our bucket is called "dataset.library.exampl.edu" and our dataset
collection is called "mycollection".

```shell
    export AWS_SDK_LOAD_CONFIG=1
    dataset init s3://dataset.library.example.edu/mycollection
    export DATASET=s3://dataset.library.example.edu/mycollection
```

We can now create a JSON record to add called "waldo" and add it to our 
collection. 

```shell
    cat<<EOT>waldo-reading.json
    {
        "reader":"Waldo",
        "author":"Robert Louis Stevenson",
        "title":"The Black Arrow",
        "url":"https://www.gutenberg.org/ebooks/848"
    }
    EOT
    cat waldo-reading.json | dataset create waldo-reading
```

List the keys in our dataset

```shell
    dataset list keys
```

Now let's download a copy of what Waldo is reading and attach it to our "waldo-reading" record.

```shell
    curl -O https://www.gutenberg.org/ebooks/848.txt.utf-8
    dataset attach waldo-reading 848.txt.utf-8
```

To check out attachments

```shell
    dataset attachments waldo-reading
```



