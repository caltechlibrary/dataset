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
version: 2.3.0
license_url: https://caltechlibrary.github.io/dataset/LICENSE

programming_language:
  - Go
  - SQL

keywords:
  - metadata
  - data
  - software
  - json

date_released: 2025-06-26
---

About this software
===================

## dataset 2.3.0

The object versioning problem identified in issue #149 persisted after the release of v2.2.8.  The root problem is the versioning.json document is not reliable given the sequences of opartions that cause it to come into existence or be updated. This lead to instability issues with v2 versioning.  In v2.3 the versioning.json file is being ignored and nolonger is saved. Instead the collection.go was updated to explicitly set the versioning state when a collection is open. This means the collections.json reflects the true state of versioning not the inferred state.

The result of changes simplifies the codebase some. It appears to be backward compatible while maintaining the flawed semver approach to versioning JSON documents. In v3 versioning will be automatic and use a simple incrementing integer value.

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

- Golang &gt;&#x3D; 1.24.4
- CMTools &gt;&#x3D; 0.0.33


### Software Suggestions

- Pandoc &gt;&#x3D; 3.1
- GNU Make &gt;&#x3D; 3.8


