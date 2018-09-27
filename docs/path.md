
# path

## Syntax

```
    dataset path COLLETION_NAME KEY
```

## Description

_path_ will return the full path to a JSON Document with the 
provided KEY.  This is particularly useful when you have your 
_dataset_ collection on local disc. This allows you to process the 
JSON document directory with whatever tools you have at hand.
Use with caution.

## Usage

In this example we are trying to find the full path to a JSON 
document with an KEY of "r1". Our collection name is "data.ds".

```shell
    dataset path data.ds r1
```

Related topics: [keys](keys.html), [read](read.html)

