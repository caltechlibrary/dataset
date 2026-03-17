---
title: dataset
abstract: "The Dataset Project provides tools for working with collections of JSON documents. It uses a simple key and object pair to organize JSON documents into a collection. It supports SQL querying of the objects stored in a collection.

It was designed for temporary storage of JSON objects in data processing pipelines. It can be used as persistent storage mechanism for collections of JSON objects you wish to distribute when used in conjuction with pairtree or SQLite3 storage.

The Dataset Project provides command line programs and a web service for working with JSON objects as a collection or individual objects. As such it is well suited for data science projects as well as building web applications that work with metadata."
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
version: 2.3.4
license_url: https://caltechlibrary.github.io/dataset/LICENSE

programming_language:
  - Go
  - SQL

keywords:
  - metadata
  - data storage
  - json

date_released: 2026-03-17
---

About this software
===================

## dataset 2.3.4

- Improved loading large JSON objects from jsonl files
- Fixed issue #164 where the queries in COLD would work in v2.2.0 but fail in v2.3.x.
- Removed support for SQL parameters in dsquery due to encoding issues and lack of practical use cases
- Removed duplicated code from dsquery.go and api_routes.go in favor of collection.go's implementation of query functionality.
- Added tailing semi-colon removal for SQL queries due to changes in behavior of SQLite3 driver

### Authors

- R. S. Doiel, <https://orcid.org/0000-0003-0900-6903>
- Thomas E Morrell, <https://orcid.org/0000-0001-9266-5146>




### Maintainers

- R. S. Doiel, <https://orcid.org/0000-0003-0900-6903>
- Thomas E Morrell, <https://orcid.org/0000-0001-9266-5146>


The Dataset Project provides tools for working with collections of JSON documents. It uses a simple key and object pair to organize JSON documents into a collection. It supports SQL querying of the objects stored in a collection.

It was designed for temporary storage of JSON objects in data processing pipelines. It can be used as persistent storage mechanism for collections of JSON objects you wish to distribute when used in conjuction with pairtree or SQLite3 storage.

The Dataset Project provides command line programs and a web service for working with JSON objects as a collection or individual objects. As such it is well suited for data science projects as well as building web applications that work with metadata.

- License: <https://caltechlibrary.github.io/dataset/LICENSE>
- GitHub: <https://github.com/caltechlibrary/dataset>
- Issues: <https://github.com/caltechlibrary/dataset/issues>

### Programming languages

- Go
- SQL




### Software Requirements

- Golang >= 1.26.1
- CMTools >= 0.0.40


### Software Suggestions

- Pandoc &gt;&#x3D; 3.1
- GNU Make &gt;&#x3D; 3.8


