# Dataset v3

The goal of Dataset v3 is to provide a simple core level of functionality making it trivial to implement light weight metadata repositories. A simple core may be easier to port to other languages or provide as a WASM library.

Version 3's focus on reducing the functionality of Dataset v2 and simplifying the codebase in the process. It is a distillation of the ideas and concepts that have guided the Dataset Project over since 2016. V3 is intended to be the penultimate implementation of dataset and datasetd. Subsequent version should be stable and feature complete.

## On the chopping block

### Dropped features from v2

- [ ] drop or rethink pairtree storage for JSON documents
- [X] frames related verbs (has been superseded by query, dsquery)
- [X] clone, clone sample (superseded by dump/load of json lines support)
- [X] join (should be handled via external tooling or via SQLite3 query support)
- [X] libdataset is being abandoned, too hard to maintain Windows build
- [ ] dsquery (merged into dataset command, already supported in datasetd)
- [X] dsimport (replaced with dump/load of json line documents)

## Revisions

- [ ] attachments and related verbs should store versioned objects in a common, bag friendly layout (OCFL v1.1 or RO-create v1.1)
- [X] default storage of metadata is in an SQLite3 database with support for PostgreSQL and MySQL maintained from v2
- [ ] a simplified model for versioned metadata.
  - [ ] one table is "current" metadata
  - [ ] second is a "history" table of versioned metadata
  - [ ] same database schema except the history uses a composite of key version number for primary key
  - [ ] history is always enabled
- [ ] dsquery merged into dataset cli
- datasetd
  - [ ] Support a read/write (GET/POST) model in addition to dogmatically following REST

## Simplified documentation

- documentations has grown organically and is difficulty to keep acurate
- manual pages should be generated from the command help
- better tutorials are needed
- examples use cases are needed
- One or more workshop presentations could be included to easy adoption of dataset

## Under consideration

- model support
  - Review model approach, integrate this into the dataset and datasetd
    - does it make sense to continue with GitHub Issue Template or is JSON
      schema stable enough to easily implement
- consider integrating a data model support
- evaluate what features might allow a dataset collection to become a turn-key collection driven web app 
- evaluate YAML data integration
  - while dataset will remain a JSON document management system, YAML support
    should be included for input and outputs

## Implementation language(s)

One of the goals of v3 will be to reduce the feature set such that implementation in Go, Rust, TypeScript and Python with feature parity of the reference Go implementation is easy.

## Targeted platforms

- aarch64
  - Linux
  - macOS
  - Windows 11
- x86_64
  - Linux
  - macOS
  - Windows 11
- Explore WASM support as an option for replacing shared library
