libdataset Python wrapper
=========================

A [wasmtime-py](https://github.com/bytecodealliance/wasmtime-py) wrapper
around `libdataset.wasm` that provides a Pythonic interface to
[dataset](https://github.com/caltechlibrary/dataset) collections.

Requirements
------------

- Python 3.9+
- wasmtime-py: `pip install wasmtime`
- `libdataset.wasm` from the dataset release

Installation
------------

Copy the `libdataset/` directory alongside your project, or install it
as a package:

```shell
pip install .   # from this directory
```

Quick start
-----------

```python
from libdataset import LibDataset, DatasetError

# Instantiate — preopens grant WASM access to host directories
ds = LibDataset(
    "path/to/libdataset.wasm",
    preopens={"/data": "/data"},
)

# (Optional) load named query config
ds.setup_file("/data/libdataset.yaml")

# Basic CRUD
ds.create("people.ds", "doe-j", {"family_name": "Doe", "given_name": "Jane"})
obj  = ds.read("people.ds", "doe-j")
keys = ds.keys("people.ds")
ds.update("people.ds", "doe-j", {**obj, "email": "jdoe@example.org"})
ds.delete("people.ds", "doe-j")

# Named query (must be configured in libdataset.yaml)
results = ds.query("people.ds", "by_family", ["Doe"])

# Attachments
ds.attach("people.ds", "doe-j", "/data/headshot.jpg")
ds.retrieve("people.ds", "doe-j", "headshot.jpg", "/tmp/headshot.jpg")
```

All methods raise `DatasetError` on failure.

See [libdataset.md](../../libdataset.md) for the full command reference
and configuration format.
