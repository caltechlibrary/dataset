# libdataset.py — wasmtime-py wrapper for the dataset WASM module.
#
# Requires: wasmtime-py  (pip install wasmtime)
#
# Usage:
#   from libdataset import LibDataset, DatasetError
#
#   ds = LibDataset("path/to/libdataset.wasm", preopens={"/data": "/data"})
#   ds.setup_file("/data/libdataset.yaml")
#   ds.create("people.ds", "doiel-r-s", {"family_name": "Doiel"})
#   obj = ds.read("people.ds", "doiel-r-s")
#   ds.close()
#
# NOTE: The WASM module prints one informational line to stdout on startup:
#   "If you're reading this, you're unnecessarily importing
#    github.com/ncruces/go-sqlite3/embed."
# This is a known cosmetic message from the ncruces SQLite driver and can
# be safely ignored. It does not affect correctness.

import json
from ctypes import c_uint8
from typing import Any, Optional

import wasmtime
from wasmtime import Store, Module, Linker, Engine, WasiConfig


class DatasetError(Exception):
    """Raised when a dataset operation returns ok=false."""
    pass


class LibDataset:
    """
    Wrapper around libdataset.wasm providing a Pythonic interface to the
    dataset collection operations.

    Parameters
    ----------
    wasm_path:
        Path to libdataset.wasm.
    preopens:
        Mapping of host paths to guest paths granted to the WASM module.
        Collections and attachment files must live under a preopened path.
        Example: {"/data": "/data"}
    """

    def __init__(self, wasm_path: str, preopens: Optional[dict] = None):
        engine = Engine()
        self._store = Store(engine)

        wasi = WasiConfig()
        wasi.inherit_env()
        if preopens:
            for host_path, guest_path in preopens.items():
                wasi.preopen_dir(host_path, guest_path)
        self._store.set_wasi(wasi)

        module = Module.from_file(engine, wasm_path)
        linker = Linker(engine)
        linker.define_wasi()
        instance = linker.instantiate(self._store, module)

        exports = instance.exports(self._store)
        self._memory  = exports["memory"]
        self._malloc  = exports["ds_malloc"]
        self._free    = exports["ds_free"]
        self._exec    = exports["ds_exec"]
        self._result  = exports["ds_result"]

    # ------------------------------------------------------------------
    # Internal helpers
    # ------------------------------------------------------------------

    def _mem_view(self):
        """Return a mutable uint8 ctypes array over WASM linear memory."""
        ptr = self._memory.data_ptr(self._store)
        length = self._memory.data_len(self._store)
        return (c_uint8 * length).from_address(ptr)  # type: ignore[return-value]

    def _call(self, cmd: dict) -> Any:
        """Serialise cmd to JSON, dispatch through ds_exec, return result value."""
        cmd_bytes = json.dumps(cmd).encode()
        mem = self._mem_view()

        # Write command into WASM memory.
        cmd_ptr = self._malloc(self._store, len(cmd_bytes))
        for i, b in enumerate(cmd_bytes):
            mem[cmd_ptr + i] = b

        result_len = self._exec(self._store, cmd_ptr, len(cmd_bytes))
        self._free(self._store, cmd_ptr, len(cmd_bytes))

        # Read result from WASM memory.
        out_ptr = self._malloc(self._store, result_len)
        self._result(self._store, out_ptr)
        raw = bytes(mem[out_ptr:out_ptr + result_len])
        self._free(self._store, out_ptr, result_len)

        resp = json.loads(raw)
        if not resp.get("ok"):
            raise DatasetError(resp.get("error", "unknown error"))
        return resp.get("result")

    # ------------------------------------------------------------------
    # Lifecycle
    # ------------------------------------------------------------------

    def version(self) -> str:
        return self._call({"op": "version"})

    def setup_file(self, config_path: str) -> None:
        """Load a libdataset YAML config and open all listed collections."""
        self._call({"op": "setup", "config": config_path})

    def collection_init(self, collection: str, dsn_uri: str = "") -> None:
        """Create a new collection at *collection* path."""
        self._call({"op": "collection_init", "collection": collection,
                    "dsn_uri": dsn_uri})

    def collection_open(self, collection: str) -> None:
        """Open an existing collection."""
        self._call({"op": "collection_open", "collection": collection})

    def collection_close(self, collection: str) -> None:
        """Close an open collection."""
        self._call({"op": "collection_close", "collection": collection})

    def codemeta(self, collection: str) -> dict:
        """Return the collection's codemeta.json as a dict."""
        return self._call({"op": "codemeta", "collection": collection})

    # ------------------------------------------------------------------
    # CRUD
    # ------------------------------------------------------------------

    def create(self, collection: str, key: str, obj: dict,
               overwrite: bool = False) -> None:
        self._call({"op": "create", "collection": collection,
                    "key": key, "object": obj, "overwrite": overwrite})

    def read(self, collection: str, key: str) -> dict:
        return self._call({"op": "read", "collection": collection, "key": key})

    def update(self, collection: str, key: str, obj: dict) -> None:
        self._call({"op": "update", "collection": collection,
                    "key": key, "object": obj})

    def delete(self, collection: str, key: str) -> None:
        self._call({"op": "delete", "collection": collection, "key": key})

    def keys(self, collection: str) -> list:
        return self._call({"op": "keys", "collection": collection})

    def has_key(self, collection: str, key: str) -> bool:
        return self._call({"op": "has_key", "collection": collection,
                           "key": key})

    def count(self, collection: str) -> int:
        return self._call({"op": "count", "collection": collection})

    def updated_keys(self, collection: str, start: str, end: str) -> list:
        return self._call({"op": "updated_keys", "collection": collection,
                           "start": start, "end": end})

    def query(self, collection: str, query_name: str,
              params: Optional[list] = None) -> list:
        return self._call({"op": "query", "collection": collection,
                           "query_name": query_name,
                           "params": params or []})

    # ------------------------------------------------------------------
    # Versioning
    # ------------------------------------------------------------------

    def get_versioning(self, collection: str) -> str:
        return self._call({"op": "get_versioning", "collection": collection})

    def set_versioning(self, collection: str, versioning: str) -> None:
        """Set versioning to '', 'patch', 'minor', or 'major'."""
        self._call({"op": "set_versioning", "collection": collection,
                    "versioning": versioning})

    def versions(self, collection: str, key: str) -> list:
        return self._call({"op": "versions", "collection": collection,
                           "key": key})

    def read_version(self, collection: str, key: str, version: str) -> dict:
        return self._call({"op": "read_version", "collection": collection,
                           "key": key, "version": version})

    # ------------------------------------------------------------------
    # Attachments
    # ------------------------------------------------------------------

    def attachments(self, collection: str, key: str) -> list:
        return self._call({"op": "attachments", "collection": collection,
                           "key": key})

    def attach(self, collection: str, key: str, filename: str) -> None:
        """Attach the file at *filename* (WASI-accessible path) to *key*."""
        self._call({"op": "attach", "collection": collection,
                    "key": key, "filename": filename})

    def retrieve(self, collection: str, key: str, filename: str,
                 output: str) -> None:
        """Write attachment *filename* to *output* (WASI-accessible path)."""
        self._call({"op": "retrieve", "collection": collection,
                    "key": key, "filename": filename, "output": output})

    def prune(self, collection: str, key: str, filename: str) -> None:
        """Remove all versions of attachment *filename* from *key*."""
        self._call({"op": "prune", "collection": collection,
                    "key": key, "filename": filename})

    # ------------------------------------------------------------------
    # Bulk / maintenance
    # ------------------------------------------------------------------

    def dump(self, collection: str, output: str) -> None:
        """Write the collection as JSONL to *output* (WASI-accessible path)."""
        self._call({"op": "dump", "collection": collection, "output": output})

    def load(self, collection: str, input_path: str,
             overwrite: bool = False) -> None:
        """Load JSONL from *input_path* (WASI-accessible path) into collection."""
        self._call({"op": "load", "collection": collection,
                    "input": input_path, "overwrite": overwrite})

    def check(self, collection: str) -> None:
        """Verify collection integrity; raises DatasetError on problems."""
        self._call({"op": "check", "collection": collection})

    def repair(self, collection: str) -> None:
        """Attempt to repair collection."""
        self._call({"op": "repair", "collection": collection})
