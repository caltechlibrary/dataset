/**
 * libdataset.ts — Deno wrapper for the dataset WASM module.
 *
 * Requires Deno 1.40+ with --allow-read and WASI support.
 *
 * NOTE: The WASM module prints one informational line to stdout on startup:
 *   "If you're reading this, you're unnecessarily importing
 *    github.com/ncruces/go-sqlite3/embed."
 * This is a known cosmetic message from the ncruces SQLite driver and can
 * be safely ignored.
 *
 * Usage:
 *   import { LibDataset, DatasetError } from "./libdataset.ts";
 *
 *   const ds = await LibDataset.load("libdataset.wasm", {
 *     preopens: { "/data": "/data" },
 *   });
 *   await ds.setupFile("/data/libdataset.yaml");
 *   await ds.create("people.ds", "doiel-r-s", { family_name: "Doiel" });
 *   const obj = await ds.read("people.ds", "doiel-r-s");
 *   ds.close();
 */

export class DatasetError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "DatasetError";
  }
}

interface Response {
  ok: boolean;
  result?: unknown;
  error?: string;
}

export interface LoadOptions {
  /** Host-to-guest path mappings granted to the WASM module via WASI. */
  preopens?: Record<string, string>;
}

export class LibDataset {
  private instance: WebAssembly.Instance;
  private memory: WebAssembly.Memory;
  private encoder = new TextEncoder();
  private decoder = new TextDecoder();

  private constructor(instance: WebAssembly.Instance) {
    this.instance = instance;
    this.memory = instance.exports.memory as WebAssembly.Memory;
  }

  /**
   * Instantiate the WASM module.
   *
   * @param wasmPath  Path to libdataset.wasm (must be readable).
   * @param options   WASI options, including preopens.
   */
  static async load(
    wasmPath: string,
    options: LoadOptions = {},
  ): Promise<LibDataset> {
    const wasmBytes = await Deno.readFile(wasmPath);
    const module = await WebAssembly.compile(wasmBytes);

    // Build WASI import object.
    // @ts-ignore — Deno.build.target exposes WASI via the wasi_snapshot_preview1 namespace.
    const wasi = new Deno.Wasi({
      preopens: options.preopens ?? {},
      env: Deno.env.toObject(),
    });

    const instance = await WebAssembly.instantiate(module, {
      // deno-lint-ignore no-explicit-any
      wasi_snapshot_preview1: (wasi as any).exports,
    });

    // WASI requires _start to be called to initialise the Go runtime.
    // deno-lint-ignore no-explicit-any
    (wasi as any).initialize(instance);

    return new LibDataset(instance);
  }

  // ------------------------------------------------------------------
  // Internal helpers
  // ------------------------------------------------------------------

  private memView(): Uint8Array {
    return new Uint8Array(this.memory.buffer);
  }

  private malloc(size: number): number {
    return (this.instance.exports.ds_malloc as CallableFunction)(size) as number;
  }

  private free(ptr: number, size: number): void {
    (this.instance.exports.ds_free as CallableFunction)(ptr, size);
  }

  private exec(cmdPtr: number, cmdLen: number): number {
    return (this.instance.exports.ds_exec as CallableFunction)(
      cmdPtr,
      cmdLen,
    ) as number;
  }

  private result(outPtr: number): void {
    (this.instance.exports.ds_result as CallableFunction)(outPtr);
  }

  private call(cmd: Record<string, unknown>): unknown {
    const mem = this.memView();
    const cmdBytes = this.encoder.encode(JSON.stringify(cmd));

    const cmdPtr = this.malloc(cmdBytes.length);
    mem.set(cmdBytes, cmdPtr);

    const resultLen = this.exec(cmdPtr, cmdBytes.length);
    this.free(cmdPtr, cmdBytes.length);

    const outPtr = this.malloc(resultLen);
    this.result(outPtr);
    const raw = this.decoder.decode(
      // slice to get a copy before freeing
      mem.slice(outPtr, outPtr + resultLen),
    );
    this.free(outPtr, resultLen);

    const resp = JSON.parse(raw) as Response;
    if (!resp.ok) {
      throw new DatasetError(resp.error ?? "unknown error");
    }
    return resp.result;
  }

  /** Release all resources. Call when finished with the instance. */
  close(): void {
    for (const key of Object.keys(this.instance.exports)) {
      if (key.startsWith("_")) continue;
    }
    // Collections are closed by the WASM runtime on exit; this is a no-op
    // in terms of the JS side but signals intent in calling code.
  }

  // ------------------------------------------------------------------
  // Lifecycle
  // ------------------------------------------------------------------

  version(): string {
    return this.call({ op: "version" }) as string;
  }

  setupFile(configPath: string): void {
    this.call({ op: "setup", config: configPath });
  }

  collectionInit(collection: string, dsnUri = ""): void {
    this.call({ op: "collection_init", collection, dsn_uri: dsnUri });
  }

  collectionOpen(collection: string): void {
    this.call({ op: "collection_open", collection });
  }

  collectionClose(collection: string): void {
    this.call({ op: "collection_close", collection });
  }

  codemeta(collection: string): Record<string, unknown> {
    return this.call({ op: "codemeta", collection }) as Record<string, unknown>;
  }

  // ------------------------------------------------------------------
  // CRUD
  // ------------------------------------------------------------------

  create(
    collection: string,
    key: string,
    obj: Record<string, unknown>,
    overwrite = false,
  ): void {
    this.call({ op: "create", collection, key, object: obj, overwrite });
  }

  read(collection: string, key: string): Record<string, unknown> {
    return this.call({ op: "read", collection, key }) as Record<string, unknown>;
  }

  update(
    collection: string,
    key: string,
    obj: Record<string, unknown>,
  ): void {
    this.call({ op: "update", collection, key, object: obj });
  }

  delete(collection: string, key: string): void {
    this.call({ op: "delete", collection, key });
  }

  keys(collection: string): string[] {
    return this.call({ op: "keys", collection }) as string[];
  }

  hasKey(collection: string, key: string): boolean {
    return this.call({ op: "has_key", collection, key }) as boolean;
  }

  count(collection: string): number {
    return this.call({ op: "count", collection }) as number;
  }

  updatedKeys(collection: string, start: string, end: string): string[] {
    return this.call({
      op: "updated_keys",
      collection,
      start,
      end,
    }) as string[];
  }

  query(
    collection: string,
    queryName: string,
    params: unknown[] = [],
  ): Record<string, unknown>[] {
    return this.call({
      op: "query",
      collection,
      query_name: queryName,
      params,
    }) as Record<string, unknown>[];
  }

  // ------------------------------------------------------------------
  // Versioning
  // ------------------------------------------------------------------

  getVersioning(collection: string): string {
    return this.call({ op: "get_versioning", collection }) as string;
  }

  /** Set versioning to "", "patch", "minor", or "major". */
  setVersioning(collection: string, versioning: string): void {
    this.call({ op: "set_versioning", collection, versioning });
  }

  versions(collection: string, key: string): string[] {
    return this.call({ op: "versions", collection, key }) as string[];
  }

  readVersion(
    collection: string,
    key: string,
    version: string,
  ): Record<string, unknown> {
    return this.call({
      op: "read_version",
      collection,
      key,
      version,
    }) as Record<string, unknown>;
  }

  // ------------------------------------------------------------------
  // Attachments
  // ------------------------------------------------------------------

  attachments(collection: string, key: string): string[] {
    return this.call({ op: "attachments", collection, key }) as string[];
  }

  /** Attach the file at *filename* (WASI-accessible path) to *key*. */
  attach(collection: string, key: string, filename: string): void {
    this.call({ op: "attach", collection, key, filename });
  }

  /** Write attachment *filename* to *output* (WASI-accessible path). */
  retrieve(
    collection: string,
    key: string,
    filename: string,
    output: string,
  ): void {
    this.call({ op: "retrieve", collection, key, filename, output });
  }

  /** Remove all versions of attachment *filename* from *key*. */
  prune(collection: string, key: string, filename: string): void {
    this.call({ op: "prune", collection, key, filename });
  }

  // ------------------------------------------------------------------
  // Bulk / maintenance
  // ------------------------------------------------------------------

  /** Write the collection as JSONL to *output* (WASI-accessible path). */
  dump(collection: string, output: string): void {
    this.call({ op: "dump", collection, output });
  }

  /** Load JSONL from *inputPath* (WASI-accessible path) into collection. */
  load(collection: string, inputPath: string, overwrite = false): void {
    this.call({ op: "load", collection, input: inputPath, overwrite });
  }

  /** Verify collection integrity; throws DatasetError on problems. */
  check(collection: string): void {
    this.call({ op: "check", collection });
  }

  /** Attempt to repair collection. */
  repair(collection: string): void {
    this.call({ op: "repair", collection });
  }
}
