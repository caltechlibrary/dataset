
Action Items
============

X
: completed

D
: declined, decided not to implement

Bugs
----

Next (prep for v2.6)
--------------------

- [X] Review SQLite3 driver: replace `github.com/ncruces/go-sqlite3` (WASM-based,
      prints spurious stderr warning on every invocation) with `github.com/glebarez/go-sqlite`
      (pure-Go, no embed overhead). Changed `Sqlite3DriverName` from `"sqlite3"` to `"sqlite"`
      in `sqlstore.go`; removed ncruces blank imports from `sqlstore.go` and dropped the sqlite
      import from `libdataset/libdataset.go` entirely (WASM/libdataset deferred).

Someday, Maybe
--------------

- [ ] dsbagit would generate a "BagIt" bag for preservation of collection
      objects
- [ ] OAI-PMH importer to prototype iiif service based on Islandora
      content driven by a dataset collection
- [ ] Implement an integrated a web UI for managing dataset collections and their data structures
  - [ ] Form pages could be expressed in Markdown+YAML for forms and embedded in the datasetd settings YAML file
    - See my notes on my text oreinted web experiment, yaml2webform.go
    - Forms could be render into the htdocs auto-magically saving development effort
    - The same forms could then be used server side for validation based on descriptors and JavaScript converted to WASM code
