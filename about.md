---
title: dataset
abstract: "The Dataset Project provides tools for working with collections of JSON documents easily. It uses a simple key and object pair to organize JSON documents into a collection. It supports SQL querying of the objects stored in a collection.

It is suitable for temporary storage of JSON objects in data processing pipelines as well as a persistent storage mechanism for collections of JSON objects.

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
version: 2.3.2
license_url: https://caltechlibrary.github.io/dataset/LICENSE

programming_language:
  - Go
  - SQL

keywords:
  - metadata
  - data
  - software
  - json

date_released: 2025-07-11
---

About this software
===================

## dataset 2.3.2

Issue #161 fix for handling GET with query were data is passed via URL parameters.

Removed support for frame, clone, sample, sync and join support removed. The dsimporter cli removed (use jsonl dump and load instead).

### Authors

- R. S. Doiel, <https://orcid.org/0000-0003-0900-6903>
- Thomas E Morrell, <https://orcid.org/0000-0001-9266-5146>




### Maintainers

- R. S. Doiel, <https://orcid.org/0000-0003-0900-6903>
- Thomas E Morrell, <https://orcid.org/0000-0001-9266-5146>


The Dataset Project provides tools for working with collections of JSON documents easily. It uses a simple key and object pair to organize JSON documents into a collection. It supports SQL querying of the objects stored in a collection.

It is suitable for temporary storage of JSON objects in data processing pipelines as well as a persistent storage mechanism for collections of JSON objects.

The Dataset Project provides command line programs and a web service for working with JSON objects as a collection or individual objects. As such it is well suited for data science projects as well as building web applications that work with metadata.

- License: <https://caltechlibrary.github.io/dataset/LICENSE>
- GitHub: <https://github.com/caltechlibrary/dataset>
- Issues: <https://github.com/caltechlibrary/dataset/issues>

### Programming languages

- Go
- SQL




### Software Requirements

- Golang &gt;&#x3D; 1.24.5
- CMTools &gt;&#x3D; 0.0.35


### Software Suggestions

- Pandoc &gt;&#x3D; 3.1
- GNU Make &gt;&#x3D; 3.8


