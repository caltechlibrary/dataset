---
title: dataset
abstract: "Dataset Project provides tools for working with JSON documents as a collection hosted in a SQL JSON store (e.g. SQLite3). It is an exploration of the use and usefulness of JSON documents in the setting of libraries, archives and museums. It may have application further afield as metadata management benefits a wide range of application domains.

Dataset v3 stores JSON documents using a SQL engine with JSON column support. These days that most popular SQL implementations support JSON columns (e.g. SQlite3, PostgreSQL, MySQL). V3 uses a extremely simple table structure for both the current object state and object history as well as integration with the SQL’s growing support for working with JSON objects generally.

Two tools are provided by the Dataset Project v3

[dataset3](dataset3.1.md)
: is a command line interface for managing JSON documents. 

[dataset3d](dataset3d.1.md)
: JSON Web Service. This service supports file attachments."
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
version: 3.0.0-alpha
license_url: https://caltechlibrary.github.io/dataset/LICENSE

programming_language:
  - Go
  - SQL

keywords:
  - metadata
  - collections
  - json
  - web service

date_released: 2025-05-28
---

About this software
===================

## dataset 3.0.0-alpha

Dataset v3 is a breaking change from v2. It focuses on feature reduction and code simplification. Many features have been removed. The v3 collections are not compatible with v2 and earlier collections.

JSON lines support in query results and read operations has been improved. The provided programs have been reduced to __dataset3__ and __dataset3d__. As such they can be installed along side Dataset v2 and earlier.

### Authors

- R. S. Doiel, <https://orcid.org/0000-0003-0900-6903>
- Thomas E Morrell, <https://orcid.org/0000-0001-9266-5146>




### Maintainers

- R. S. Doiel, <https://orcid.org/0000-0003-0900-6903>
- Thomas E Morrell, <https://orcid.org/0000-0001-9266-5146>


Dataset Project provides tools for working with JSON documents as a collection hosted in a SQL JSON store (e.g. SQLite3). It is an exploration of the use and usefulness of JSON documents in the setting of libraries, archives and museums. It may have application further afield as metadata management benefits a wide range of application domains.

Dataset v3 stores JSON documents using a SQL engine with JSON column support. These days that most popular SQL implementations support JSON columns (e.g. SQlite3, PostgreSQL, MySQL). V3 uses a extremely simple table structure for both the current object state and object history as well as integration with the SQL’s growing support for working with JSON objects generally.

Two tools are provided by the Dataset Project v3

[dataset3](dataset3.1.md)
: is a command line interface for managing JSON documents. 

[dataset3d](dataset3d.1.md)
: JSON Web Service. This service supports file attachments.

- License: <https://caltechlibrary.github.io/dataset/LICENSE>
- GitHub: <https://github.com/caltechlibrary/dataset>
- Issues: <https://github.com/caltechlibrary/dataset/issues>

### Programming languages

- Go
- SQL




### Software Requirements

- Golang &gt;&#x3D; 1.24.2
- CMTools &gt;&#x3D; 0.0.29
- Pandoc &gt;&#x3D; 3.1
- GNU Make &gt;&#x3D; 3.8


### Software Suggestions

- jq &gt;&#x3D; 1.7


