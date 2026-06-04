%dataset(5) user manual | version 2.5.X
% R. S. Doiel and Tom Morrell
% 2026-06-03

# Compatibility in v2.5

## v2.5.1 — SQLite driver split (2026-06-03)

### Problem

In v2.5.0 `ncruces/go-sqlite3` was adopted as the sole SQLite driver to
enable the `GOOS=wasip1 GOARCH=wasm` build of `libdataset`. The ncruces
driver is WASM-based by design: even for ordinary Linux/macOS/Windows
builds it runs SQLite inside an embedded WASM runtime, carrying runtime
overhead on every invocation. Additionally, when the driver is imported
without its companion `embed` sub-package the WASM binary is not bundled
and the driver prints a spurious warning to stderr on every program
startup.

### Decision

Split the SQLite driver selection across two build-tagged files so each
build target gets the best-fit driver:

| Build target | File | Driver | Registered name |
|---|---|---|---|
| All non-WASM targets | `sqlstore_sqlite.go` (`!wasip1`) | `glebarez/go-sqlite` (pure-Go via `modernc.org/sqlite`) | `"sqlite"` |
| `GOOS=wasip1 GOARCH=wasm` | `sqlstore_sqlite_wasm.go` (`wasip1`) | `ncruces/go-sqlite3/driver` + `/embed` | `"sqlite3"` |

`glebarez/go-sqlite` is a thin shim over `modernc.org/sqlite`, which is
SQLite compiled to pure Go via ccgo. It requires no CGO, no external
binary, and no WASM runtime for native builds. The `ncruces` driver with
the `embed` sub-package bundles the SQLite WASM binary directly into the
compiled WASM module, which is required for the Python (wasmtime-py) and
TypeScript/Deno `libdataset` wrappers.

`modernc.org/libc` (pulled in by glebarez) does not support
`GOOS=wasip1`, so it cannot be used for WASM builds. Keeping ncruces for
that target is therefore mandatory for the foreseeable future.

### Impact on callers

- **DSN URIs** (`sqlite://collection.db`) are unchanged. The `"sqlite"`
  scheme in the URI is always mapped to whichever registered driver name
  is in use via `driverNameFixUp`.
- **Non-WASM builds**: the internal constant `Sqlite3DriverName` changes
  from `"sqlite3"` to `"sqlite"`. This constant is not part of the public
  API; callers that open collections through `Init` / `SQLStoreOpen` are
  unaffected.
- **WASM builds**: `Sqlite3DriverName` remains `"sqlite3"` and behaviour
  is identical to v2.5.0.
- Existing `.db` files are fully compatible; no migration is required.

## v2.5.0 — libdataset WASM module and ncruces SQLite driver (2025-07-XX)

Replaced the cgo-based C shared library with a pure WASM module built
with `GOOS=wasip1 GOARCH=wasm`. Switched SQLite driver from
`glebarez/go-sqlite` to `ncruces/go-sqlite3` to enable the WASM build.
Existing collection files and `sqlite://...` DSN URIs are unchanged.
