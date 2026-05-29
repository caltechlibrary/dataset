---
title: dataset
abstract: |-
  The Dataset Project provides tools for working with collections of JSON documents. It uses a simple key and object pair to organize JSON documents into a collection. It supports SQL querying of the objects stored in a collection.

  It was designed for temporary storage of JSON objects in data processing pipelines. It can be used as persistent storage mechanism for collections of JSON objects you wish to distribute when used in conjunction with pairtree or SQLite3 storage.

  The Dataset Project provides command line programs and a web service for working with JSON objects as a collection or individual objects. As such it is well suited for data science and web base applications.
authors:
  - family_name: Doiel
    given_name: R. S.
    id: https://orcid.org/0000-0003-0900-6903
  - family_name: Morrell
    given_name: Thomas E
    id: https://orcid.org/0000-0001-9266-5146


maintainer:
  - family_name: Doiel
    given_name: R. S.
    id: https://orcid.org/0000-0003-0900-6903
  - family_name: Morrell
    given_name: Thomas E
    id: https://orcid.org/0000-0001-9266-5146

repository_code: https://github.com/caltechlibrary/dataset
version: 2.5.0
license_url: https://caltechlibrary.github.io/dataset/LICENSE

programming_language:
  - Go
  - SQL

keywords:
  - metadata
  - data storage
  - JSON

date_released: 2026-05-29
---

About this software
===================

## dataset 2.5.0

- Added libdataset WASM module (GOOS=wasip1) replacing the old cgo-based C shared library, with Python (wasmtime-py) and TypeScript/Deno wrappers
- Named SQL queries pre-configured via settings.yaml-style config; new release artifact dataset-vVERSION-libdataset.zip
- Switched SQLite driver to ncruces/go-sqlite3 (SQLite 3.53.1, up from 3.41.2), enabling WASM builds while retaining FTS5 and full feature set
- Fixed datasetd request body size limits not being enforced (http.MaxBytesReader now applied)
- Fixed datasetd error redirect using HTTP 304 instead of HTTP 303
- Replaced deprecated io/ioutil throughout

## Authors

- [R. S. Doiel](https://orcid.org/0000-0003-0900-6903)
- [Thomas E Morrell](https://orcid.org/0000-0001-9266-5146)




## Maintainers

- [R. S. Doiel](https://orcid.org/0000-0003-0900-6903)
- [Thomas E Morrell](https://orcid.org/0000-0001-9266-5146)


The Dataset Project provides tools for working with collections of JSON documents. It uses a simple key and object pair to organize JSON documents into a collection. It supports SQL querying of the objects stored in a collection.

It was designed for temporary storage of JSON objects in data processing pipelines. It can be used as persistent storage mechanism for collections of JSON objects you wish to distribute when used in conjunction with pairtree or SQLite3 storage.

The Dataset Project provides command line programs and a web service for working with JSON objects as a collection or individual objects. As such it is well suited for data science and web base applications.

- [License](https://caltechlibrary.github.io/dataset/LICENSE)
- [Code Repository](https://github.com/caltechlibrary/dataset)
  - [Issue Tracker](https://github.com/caltechlibrary/dataset/issues)

## Programming languages

- Go
- SQL




## Software Requirements

- Golang >= 1.26
- CMTools >= 0.0.43


## Software Suggestions

- Pandoc >= 3.9
- GNU Make >= 3.8


