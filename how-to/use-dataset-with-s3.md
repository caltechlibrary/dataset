
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



