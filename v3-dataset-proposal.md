
# Dataset v3

The goal of Dataset v3 is to provide a simple core level of functionality making it trivial to implement light weight metadata repositories.

Version 3 will focus on reducing the functionality of Dataset (i.e. dataset and datasetd) to it to a core that is readily implemented in other programming languages with basic JSON support. It is a distillation of the ideas and concepts that have guided the Dataset Project over since 2016. V3 is intended to be the penultimate implementation of dataset and datasetd. Subsequent version should be stable and feature complete.


## On the chopping block

### Dropped features from v2

- pairtree storage
- frames related verbs (has been superseded by query, dsquery)
- clone, clone sample (superseded by dump/load of json lines support)
- join (should be handled via external tooling or via SQLite3 query support)
- libdataset is being abandoned, too hard to maintain Windows build
- dsquery (merged into dataset command, already supported in datasetd)
- dsimport (replaced with dump/load of json line documents)

## Revisions

- default storage of metadata is in an SQLite3 database
- a simplified model for versioned metadata.
  - one table is "current" metadata
  - second is a "history" table of versioned metadata
  - same database schema except the history uses a composite of key version number for primary key
- attachments support moved out of the dataset directory structure
 - attachment support targeting OCFL via local file system or S3 protocol
- dsquery merged into dataset cli
- datasetd
  - Look at improving and documenting the JSON API such that it is clearer how
    SQLite3 queries are mapped with named parameters

## Simplified documentation

- documentations has grown organically and is difficulty to keep acurate
- manual pages should be generated from the command help
- better tutorials are needed
- examples use cases are needed

## Under consideration

- model support 
  - Review model approach, integrate this into the dataset and datasetd
    - does it make sense to continue with GitHub Issue Template or is JSON schema stable enough to easily implement
- YAML integration
  - while dataset will remain a JSON document management system, YAML support should be included for input and outputs

## Implementation language(s)

One of the goals of v3 will be to reduce the feature set such that implementation in Go, Rust, TypeScript and Python with feature parity of the reference implementation is easy.

## Targeted platforms

- aarch64
  - Linux
  - macOS
  - Windows 11
- x86_64
  - Linux
  - macOS
  - Windows 11
- Explore WASM support
