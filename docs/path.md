
# path

## Syntax

```
    dataset COLLETION_NAME path KEY
```

## Description

path will return the full path to a JSON Document with the provided KEY.
This is particularly useful when you have your _dataset_ collection on local disc. This
allows you to process the JSON document directory with whatever tools you have at hand.
Use with caution.

## Usage

In this example we are trying to find the full path to a JSON document with an KEY
of "r1". Our collection name is "data.ds".

```shell
    dataset data.ds path r1
```

