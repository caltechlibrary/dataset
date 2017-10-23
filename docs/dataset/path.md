
# path

## Syntax

```
    dataset path JSON_RECORD_ID
```

## Description

path will return the full path to a JSON Document with the provided JSON_RECORD_ID.
This is particularly useful when you have your _dataset_ collection on local disc. This
allows you to process the JSON document directory with whatever tools you have at hand.
Use with caution.

## Usage

In this example we are trying to find the full path to a JSON document with an JSON_RECORD_ID
of "r1".

```shell
    dataset path r1
```

