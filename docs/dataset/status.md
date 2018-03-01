
# status

## Syntax

```
    dataset status COLLECTION_NAME [COLLECTION_NAME ...]
```

## Description

Checks to see if a `collection.json` file is associated with the COLLECTION_NAME. Can work on multiple
collection names. Returns "OK" if it is.

## Usage

```
    dataset status MyRecordCollection.ds
    dataset status MyRecordCollection.ds MyBookCollection.ds
```

