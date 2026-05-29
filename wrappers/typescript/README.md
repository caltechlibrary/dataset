libdataset TypeScript wrapper (Deno)
=====================================

A Deno wrapper around `libdataset.wasm` that provides a TypeScript
interface to [dataset](https://github.com/caltechlibrary/dataset)
collections.

Requirements
------------

- [Deno](https://deno.land/) 1.40+
- `libdataset.wasm` from the dataset release

Quick start
-----------

```typescript
import { LibDataset, DatasetError } from "./libdataset.ts";

// Instantiate — preopens grant WASM access to host directories
const ds = await LibDataset.load("path/to/libdataset.wasm", {
  preopens: { "/data": "/data" },
});

// (Optional) load named query config
ds.setupFile("/data/libdataset.yaml");

// Basic CRUD
ds.create("people.ds", "doe-j", { family_name: "Doe", given_name: "Jane" });
const obj  = ds.read("people.ds", "doe-j");
const keys = ds.keys("people.ds");
ds.update("people.ds", "doe-j", { ...obj, email: "jdoe@example.org" });
ds.delete("people.ds", "doe-j");

// Named query (must be configured in libdataset.yaml)
const results = ds.query("people.ds", "by_family", ["Doe"]);

// Attachments
ds.attach("people.ds", "doe-j", "/data/headshot.jpg");
ds.retrieve("people.ds", "doe-j", "headshot.jpg", "/tmp/headshot.jpg");
```

Run with the permissions the WASM module needs:

```shell
deno run --allow-read --allow-write your_script.ts
```

All methods throw `DatasetError` on failure.

See [libdataset.md](../../libdataset.md) for the full command reference
and configuration format.
