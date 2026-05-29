libdataset
==========

libdataset is a [WebAssembly](https://webassembly.org/) (WASM) module that
exposes the full [dataset](README.md) collection API to Python, TypeScript,
JavaScript, and any other language with a WASI-compatible WASM runtime.

It replaces the old cgo-based `libdataset` C shared library with a single
portable binary (`libdataset.wasm`) that requires no compilation or
platform-specific packaging.

Overview
--------

`libdataset.wasm` is built for the [WASI Preview 1](https://wasi.dev/)
target (`GOOS=wasip1 GOARCH=wasm`). It exposes four low-level exports
(`ds_malloc`, `ds_free`, `ds_exec`, `ds_result`) through which the host
exchanges JSON-encoded commands and responses via shared linear memory.
The provided Python and TypeScript wrappers handle all memory management
so application code never touches the ABI directly.

The module supports the same operations as the `dataset` CLI:

- Collection lifecycle: init, open, close, codemeta
- CRUD: create, read, update, delete
- Query: keys, has\_key, count, updated\_keys, named SQL queries
- Versioning: get/set versioning, list versions, read a specific version
- Attachments: list, attach, retrieve, prune
- Bulk: dump (export JSONL), load (import JSONL)
- Maintenance: check, repair

Prerequisites
-------------

### Python

- Python 3.9+
- [wasmtime-py](https://github.com/bytecodealliance/wasmtime-py)

```shell
pip install wasmtime
```

### TypeScript / Deno

- [Deno](https://deno.land/) 1.40+

No additional packages are required; Deno includes native WASI support.

Installation
------------

Download the `dataset-vVERSION-libdataset.zip` file from the
[releases page](https://github.com/caltechlibrary/dataset/releases) and
unzip it. The archive contains:

```
libdataset.wasm          the WASM module (platform-independent)
python/libdataset/       Python package
typescript/libdataset.ts TypeScript module for Deno
LICENSE
libdataset.md            this file
```

Place `libdataset.wasm` somewhere accessible to your project, then copy
or install the wrapper that matches your language.

Building from source
--------------------

```shell
git clone https://github.com/caltechlibrary/dataset
cd dataset
make libdataset.wasm        # builds dist/libdataset.wasm
```

Requires Go 1.24+ with the standard WASM/WASI toolchain (no extra tools).

WASI filesystem access
----------------------

WASM modules cannot read or write the host filesystem unless the host
explicitly grants access via **preopens**. You must preopen every
directory that contains collections, attachment files, or JSONL export
paths.

```python
# Python — grant /data read/write access
ds = LibDataset("libdataset.wasm", preopens={"/data": "/data"})
```

```typescript
// Deno — grant /data read/write access
const ds = await LibDataset.load("libdataset.wasm", {
  preopens: { "/data": "/data" },
});
```

A common pattern during development is to preopen `"/"` to grant full
filesystem access, then restrict to specific paths in production.

Configuration
-------------

Named SQL queries must be pre-configured in a YAML file before use.
The format is identical to `datasetd`'s `settings.yaml`, minus the
HTTP-specific fields:

```yaml
# libdataset.yaml
collections:
  - dataset: /data/people.ds
    dsn_uri: sqlite://collection.db
    query:
      by_family: >-
        SELECT src FROM people
        WHERE json_extract(src, '$.family_name') = ?
      recent: >-
        SELECT src FROM people
        WHERE updated > ?
```

Load the configuration with `setup_file()` before calling `query()`.
All collections listed in the config are opened automatically.

### Python usage

```python
from libdataset import LibDataset, DatasetError

ds = LibDataset("libdataset.wasm", preopens={"/data": "/data"})
ds.setup_file("/data/libdataset.yaml")

# Create a record
ds.create("people.ds", "doiel-r-s", {
    "family_name": "Doiel",
    "given_name":  "Robert",
    "orcid":       "0000-0003-0900-6903",
})

# Read it back
obj = ds.read("people.ds", "doiel-r-s")
print(obj["family_name"])   # Doiel

# List all keys
keys = ds.keys("people.ds")

# Named query (requires libdataset.yaml)
results = ds.query("people.ds", "by_family", ["Doiel"])

# Attach a file
ds.attach("people.ds", "doiel-r-s", "/data/photo.jpg")

# Retrieve it
ds.retrieve("people.ds", "doiel-r-s", "photo.jpg", "/tmp/photo.jpg")

# Update, delete
obj["email"] = "rsdoiel@example.org"
ds.update("people.ds", "doiel-r-s", obj)
ds.delete("people.ds", "doiel-r-s")
```

Errors raise `DatasetError`.

### TypeScript / Deno usage

```typescript
import { LibDataset, DatasetError } from "./libdataset.ts";

const ds = await LibDataset.load("libdataset.wasm", {
  preopens: { "/data": "/data" },
});
ds.setupFile("/data/libdataset.yaml");

ds.create("people.ds", "doiel-r-s", {
  family_name: "Doiel",
  given_name:  "Robert",
  orcid:       "0000-0003-0900-6903",
});

const obj   = ds.read("people.ds", "doiel-r-s");
const keys  = ds.keys("people.ds");
const rows  = ds.query("people.ds", "by_family", ["Doiel"]);

ds.update("people.ds", "doiel-r-s", { ...obj, email: "rsdoiel@example.org" });
ds.delete("people.ds", "doiel-r-s");
```

Run with:

```shell
deno run --allow-read --allow-write your_script.ts
```

Errors throw `DatasetError`.

Command reference
-----------------

All commands are issued through the JSON ABI. The wrapper methods map
directly to the operations listed here.

### Lifecycle

| Method | Op | Parameters |
|--------|-----|------------|
| `version()` | `version` | — |
| `setup_file(config)` | `setup` | `config`: path to YAML |
| `collection_init(collection, dsn_uri)` | `collection_init` | `collection`, `dsn_uri` |
| `collection_open(collection)` | `collection_open` | `collection` |
| `collection_close(collection)` | `collection_close` | `collection` |
| `codemeta(collection)` | `codemeta` | `collection` |

### CRUD

| Method | Op | Parameters |
|--------|-----|------------|
| `create(collection, key, object[, overwrite])` | `create` | |
| `read(collection, key)` | `read` | |
| `update(collection, key, object)` | `update` | |
| `delete(collection, key)` | `delete` | |
| `keys(collection)` | `keys` | |
| `has_key(collection, key)` | `has_key` | |
| `count(collection)` | `count` | |
| `updated_keys(collection, start, end)` | `updated_keys` | ISO datetime strings |
| `query(collection, query_name, params)` | `query` | `params`: ordered list |

### Versioning

| Method | Op |
|--------|-----|
| `get_versioning(collection)` | `get_versioning` |
| `set_versioning(collection, versioning)` | `set_versioning` — `""`, `"patch"`, `"minor"`, `"major"` |
| `versions(collection, key)` | `versions` |
| `read_version(collection, key, version)` | `read_version` |

### Attachments

| Method | Op |
|--------|-----|
| `attachments(collection, key)` | `attachments` |
| `attach(collection, key, filename)` | `attach` — `filename` must be WASI-accessible |
| `retrieve(collection, key, filename, output)` | `retrieve` — `output` must be WASI-accessible |
| `prune(collection, key, filename)` | `prune` |

### Bulk / maintenance

| Method | Op |
|--------|-----|
| `dump(collection, output)` | `dump` — writes JSONL to `output` |
| `load(collection, input[, overwrite])` | `load` — reads JSONL from `input` |
| `check(collection)` | `check` |
| `repair(collection)` | `repair` |

Known issues
------------

When the WASM module starts it prints one informational line:

```
If you're reading this, you're unnecessarily importing github.com/ncruces/go-sqlite3/embed.
```

This is a cosmetic message from the SQLite driver embedded in the module
and does not indicate an error. It can be safely ignored. It will be
suppressed in a future release of the ncruces/go-sqlite3 driver.
