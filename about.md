---
title: dataset
abstract: "The Dataset Project provides tools for working with collections of JSON documents. It uses a simple key and object pair to organize JSON documents into a collection. It supports SQL querying of the objects stored in a collection.

It was designed for temporary storage of JSON objects in data processing pipelines. It can be used as persistent storage mechanism for collections of JSON objects you wish to distribute when used in conjunction with pairtree or SQLite3 storage.

The Dataset Project provides command line programs and a web service for working with JSON objects as a collection or individual objects. As such it is well suited for data science and web base applications."
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
version: 2.4.0
license_url: https://caltechlibrary.github.io/dataset/LICENSE

programming_language:
  - Go
  - SQL

keywords:
  - metadata
  - data storage
  - JSON

date_released: 2026-03-20
---

About this software
===================

## dataset 2.4.0

- Removed MySQL support
- Added schema-based validation for dataset collections using the models package
- Added support for nested object and list structures in schema definitions
- Added identifier types (ISBN, ISSN, DOI, ORCID, ROR, ISNI, PMID, PMCID, FundRef, LCNAF, VIAF, SNAC, ArXiv, EAN) for CrossRef/DataCite record validation
- Added schemas configuration to settings.yaml for defining reusable validation schemas
- Added schema_name and validate fields to collection configuration for enabling per-collection validation
- API now returns X-Validation-Errors header with detailed validation error information on create/update failures

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


