---
title: dataset
abstract: "Tools for working with JSON documents as a collection hosted on the file system or SQL JSON store"
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
version: 2.2.0
license_url: https://caltechlibrary.github.io/dataset/LICENSE

programming_language:
  - Go

keywords:
  - GitHub
  - metadata
  - data
  - software
  - json

date_released: 2025-04-09
---

About this software
===================

## dataset 2.2.0

This minor release see the addition of two new dataset verbs and
the introduction of SQLite3 as the default storage type. You can
still create a pairtree store but now you need to include that as
a paramter when invoking the init verb.

The added verbs are dump and load. These offer a different
approach than cloning repositories. The dump verb will write a JSONL
object stream to standard out where the objects have two attributes,
key and object. The key attribute corresponds to the object key in the
dataset collection while the object attribute contains the JSON object
in the collection.  The load command can read this stream of objects
and use them to populate a collection.

### Authors

- R. S. Doiel, <https://orcid.org/0000-0003-0900-6903>
- Thomas E Morrell, <https://orcid.org/0000-0001-9266-5146>




### Maintainers

- R. S. Doiel, <https://orcid.org/0000-0003-0900-6903>
- Thomas E Morrell, <https://orcid.org/0000-0001-9266-5146>


Tools for working with JSON documents as a collection hosted on the file system or SQL JSON store

- License: <https://caltechlibrary.github.io/dataset/LICENSE>
- GitHub: <https://github.com/caltechlibrary/dataset>
- Issues: <https://github.com/caltechlibrary/dataset/issues>

### Programming languages

- Go




### Software Requirements

- Golang &gt;&#x3D; 1.24.2
- CMTOlls &gt;&#x3D; 0.0.20
- Pandoc &gt;&#x3D; 3.1
- GNU Make &gt;&#x3D; 3.8

