//go:build wasip1

// libdataset is a WASM module (GOOS=wasip1 GOARCH=wasm) exposing the dataset
// package through a four-function ABI. Callers exchange JSON-encoded commands
// and responses through shared linear memory.
//
// # ABI
//
//	ptr := ds_malloc(len)          allocate len bytes, return pointer
//	ds_free(ptr, len)              release a previous allocation
//	rlen := ds_exec(cmdPtr, cmdLen) dispatch a command, return result length
//	ds_result(outPtr)              copy last result into caller-provided buffer
//
// # Command format
//
//	{"op": "<name>", ... fields ...}
//
// # Response format
//
//	{"ok": true,  "result": <value>}
//	{"ok": false, "error": "<message>"}
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"unsafe"

	dataset "github.com/caltechlibrary/dataset/v2"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// pins prevents the GC from collecting buffers we hand to the host.
var pins = map[int32][]byte{}

// lastResult holds the JSON-encoded result of the most recent ds_exec call.
var lastResult []byte

//go:wasmexport ds_malloc
func ds_malloc(size int32) int32 {
	b := make([]byte, size)
	p := int32(uintptr(unsafe.Pointer(&b[0])))
	pins[p] = b
	return p
}

//go:wasmexport ds_free
func ds_free(ptr, size int32) {
	delete(pins, ptr)
}

//go:wasmexport ds_exec
func ds_exec(cmdPtr, cmdLen int32) int32 {
	src := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(cmdPtr))), cmdLen)
	var c command
	if err := json.Unmarshal(src, &c); err != nil {
		lastResult = errResult(fmt.Sprintf("invalid command JSON: %s", err))
		return int32(len(lastResult))
	}
	fn, ok := dispatch[c.Op]
	if !ok {
		lastResult = errResult(fmt.Sprintf("unknown op %q", c.Op))
		return int32(len(lastResult))
	}
	lastResult = fn(&c)
	return int32(len(lastResult))
}

//go:wasmexport ds_result
func ds_result(outPtr int32) {
	dest := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(outPtr))), len(lastResult))
	copy(dest, lastResult)
}

// command is the unified request structure for all ops.
type command struct {
	Op         string                 `json:"op"`
	Config     string                 `json:"config,omitempty"`
	Collection string                 `json:"collection,omitempty"`
	DsnURI     string                 `json:"dsn_uri,omitempty"`
	Key        string                 `json:"key,omitempty"`
	Object     map[string]interface{} `json:"object,omitempty"`
	QueryName  string                 `json:"query_name,omitempty"`
	Params     []interface{}          `json:"params,omitempty"`
	Versioning string                 `json:"versioning,omitempty"`
	Version    string                 `json:"version,omitempty"`
	Filename   string                 `json:"filename,omitempty"`
	Output     string                 `json:"output,omitempty"`
	Input      string                 `json:"input,omitempty"`
	Start      string                 `json:"start,omitempty"`
	End        string                 `json:"end,omitempty"`
	Overwrite  bool                   `json:"overwrite,omitempty"`
}

type response struct {
	OK     bool        `json:"ok"`
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func okResult(v interface{}) []byte {
	b, _ := json.Marshal(response{OK: true, Result: v})
	return b
}

func errResult(msg string) []byte {
	b, _ := json.Marshal(response{OK: false, Error: msg})
	return b
}

var dispatch = map[string]func(*command) []byte{
	// Lifecycle
	"version":          opVersion,
	"setup":            opSetup,
	"collection_init":  opCollectionInit,
	"collection_open":  opCollectionOpen,
	"collection_close": opCollectionClose,
	"codemeta":         opCodemeta,
	// CRUD
	"create":       opCreate,
	"read":         opRead,
	"update":       opUpdate,
	"delete":       opDelete,
	"keys":         opKeys,
	"has_key":      opHasKey,
	"count":        opCount,
	"updated_keys": opUpdatedKeys,
	"query":        opQuery,
	// Versioning
	"get_versioning": opGetVersioning,
	"set_versioning": opSetVersioning,
	"versions":       opVersions,
	"read_version":   opReadVersion,
	// Attachments
	"attachments": opAttachments,
	"attach":      opAttach,
	"retrieve":    opRetrieve,
	"prune":       opPrune,
	// Bulk / maintenance
	"dump":   opDump,
	"load":   opLoad,
	"check":  opCheck,
	"repair": opRepair,
}

// ---- Lifecycle ----

func opVersion(_ *command) []byte {
	return okResult(dataset.Version)
}

func opSetup(c *command) []byte {
	if c.Config == "" {
		return errResult("missing config path")
	}
	if err := setupFromConfig(c.Config); err != nil {
		return errResult(err.Error())
	}
	return okResult(true)
}

func opCollectionInit(c *command) []byte {
	if c.Collection == "" {
		return errResult("missing collection")
	}
	if err := initCollection(c.Collection, c.DsnURI); err != nil {
		return errResult(err.Error())
	}
	return okResult(true)
}

func opCollectionOpen(c *command) []byte {
	if c.Collection == "" {
		return errResult("missing collection")
	}
	if err := openCollection(c.Collection); err != nil {
		return errResult(err.Error())
	}
	return okResult(true)
}

func opCollectionClose(c *command) []byte {
	if c.Collection == "" {
		return errResult("missing collection")
	}
	if err := closeCollection(c.Collection); err != nil {
		return errResult(err.Error())
	}
	return okResult(true)
}

func opCodemeta(c *command) []byte {
	e, err := getEntry(c.Collection)
	if err != nil {
		return errResult(err.Error())
	}
	src, err := e.c.Codemeta()
	if err != nil {
		return errResult(err.Error())
	}
	var obj interface{}
	if err := json.Unmarshal(src, &obj); err != nil {
		return errResult(err.Error())
	}
	return okResult(obj)
}

// ---- CRUD ----

func opCreate(c *command) []byte {
	if c.Key == "" {
		return errResult("missing key")
	}
	e, err := getEntry(c.Collection)
	if err != nil {
		return errResult(err.Error())
	}
	if c.Overwrite && e.c.HasKey(c.Key) {
		if err := e.c.Update(c.Key, c.Object); err != nil {
			return errResult(err.Error())
		}
		return okResult(true)
	}
	if err := e.c.Create(c.Key, c.Object); err != nil {
		return errResult(err.Error())
	}
	return okResult(true)
}

func opRead(c *command) []byte {
	if c.Key == "" {
		return errResult("missing key")
	}
	e, err := getEntry(c.Collection)
	if err != nil {
		return errResult(err.Error())
	}
	src, err := e.c.ReadJSON(c.Key)
	if err != nil {
		return errResult(err.Error())
	}
	var obj interface{}
	if err := json.Unmarshal(src, &obj); err != nil {
		return errResult(err.Error())
	}
	return okResult(obj)
}

func opUpdate(c *command) []byte {
	if c.Key == "" {
		return errResult("missing key")
	}
	e, err := getEntry(c.Collection)
	if err != nil {
		return errResult(err.Error())
	}
	if err := e.c.Update(c.Key, c.Object); err != nil {
		return errResult(err.Error())
	}
	return okResult(true)
}

func opDelete(c *command) []byte {
	if c.Key == "" {
		return errResult("missing key")
	}
	e, err := getEntry(c.Collection)
	if err != nil {
		return errResult(err.Error())
	}
	if err := e.c.Delete(c.Key); err != nil {
		return errResult(err.Error())
	}
	return okResult(true)
}

func opKeys(c *command) []byte {
	e, err := getEntry(c.Collection)
	if err != nil {
		return errResult(err.Error())
	}
	keys, err := e.c.Keys()
	if err != nil {
		return errResult(err.Error())
	}
	if keys == nil {
		keys = []string{}
	}
	return okResult(keys)
}

func opHasKey(c *command) []byte {
	e, err := getEntry(c.Collection)
	if err != nil {
		return errResult(err.Error())
	}
	return okResult(e.c.HasKey(c.Key))
}

func opCount(c *command) []byte {
	e, err := getEntry(c.Collection)
	if err != nil {
		return errResult(err.Error())
	}
	return okResult(e.c.Length())
}

func opUpdatedKeys(c *command) []byte {
	e, err := getEntry(c.Collection)
	if err != nil {
		return errResult(err.Error())
	}
	keys, err := e.c.UpdatedKeys(c.Start, c.End)
	if err != nil {
		return errResult(err.Error())
	}
	if keys == nil {
		keys = []string{}
	}
	return okResult(keys)
}

func opQuery(c *command) []byte {
	e, err := getEntry(c.Collection)
	if err != nil {
		return errResult(err.Error())
	}
	if len(e.queryFn) == 0 {
		return errResult("no queries configured for this collection")
	}
	sql, ok := e.queryFn[c.QueryName]
	if !ok {
		return errResult(fmt.Sprintf("query %q not found", c.QueryName))
	}
	results, err := e.c.Query(sql, false, c.Params)
	if err != nil {
		return errResult(err.Error())
	}
	if results == nil {
		results = []interface{}{}
	}
	return okResult(results)
}

// ---- Versioning ----

func opGetVersioning(c *command) []byte {
	e, err := getEntry(c.Collection)
	if err != nil {
		return errResult(err.Error())
	}
	return okResult(e.c.Versioning)
}

func opSetVersioning(c *command) []byte {
	e, err := getEntry(c.Collection)
	if err != nil {
		return errResult(err.Error())
	}
	if err := e.c.SetVersioning(c.Versioning); err != nil {
		return errResult(err.Error())
	}
	return okResult(true)
}

func opVersions(c *command) []byte {
	if c.Key == "" {
		return errResult("missing key")
	}
	e, err := getEntry(c.Collection)
	if err != nil {
		return errResult(err.Error())
	}
	versions, err := e.c.Versions(c.Key)
	if err != nil {
		return errResult(err.Error())
	}
	if versions == nil {
		versions = []string{}
	}
	return okResult(versions)
}

func opReadVersion(c *command) []byte {
	if c.Key == "" {
		return errResult("missing key")
	}
	if c.Version == "" {
		return errResult("missing version")
	}
	e, err := getEntry(c.Collection)
	if err != nil {
		return errResult(err.Error())
	}
	src, err := e.c.ReadJSONVersion(c.Key, c.Version)
	if err != nil {
		return errResult(err.Error())
	}
	var obj interface{}
	if err := json.Unmarshal(src, &obj); err != nil {
		return errResult(err.Error())
	}
	return okResult(obj)
}

// ---- Attachments ----

func opAttachments(c *command) []byte {
	if c.Key == "" {
		return errResult("missing key")
	}
	e, err := getEntry(c.Collection)
	if err != nil {
		return errResult(err.Error())
	}
	fNames, err := e.c.Attachments(c.Key)
	if err != nil {
		return errResult(err.Error())
	}
	if fNames == nil {
		fNames = []string{}
	}
	return okResult(fNames)
}

func opAttach(c *command) []byte {
	if c.Key == "" {
		return errResult("missing key")
	}
	if c.Filename == "" {
		return errResult("missing filename")
	}
	e, err := getEntry(c.Collection)
	if err != nil {
		return errResult(err.Error())
	}
	f, err := os.Open(c.Filename)
	if err != nil {
		return errResult(fmt.Sprintf("opening %q: %s", c.Filename, err))
	}
	defer f.Close()
	if err := e.c.AttachStream(c.Key, c.Filename, f); err != nil {
		return errResult(err.Error())
	}
	return okResult(true)
}

func opRetrieve(c *command) []byte {
	if c.Key == "" {
		return errResult("missing key")
	}
	if c.Filename == "" {
		return errResult("missing filename")
	}
	if c.Output == "" {
		return errResult("missing output path")
	}
	e, err := getEntry(c.Collection)
	if err != nil {
		return errResult(err.Error())
	}
	f, err := os.Create(c.Output)
	if err != nil {
		return errResult(fmt.Sprintf("creating %q: %s", c.Output, err))
	}
	defer f.Close()
	if err := e.c.RetrieveStream(c.Key, c.Filename, f); err != nil {
		return errResult(err.Error())
	}
	return okResult(true)
}

func opPrune(c *command) []byte {
	if c.Key == "" {
		return errResult("missing key")
	}
	if c.Filename == "" {
		return errResult("missing filename")
	}
	e, err := getEntry(c.Collection)
	if err != nil {
		return errResult(err.Error())
	}
	if err := e.c.Prune(c.Key, c.Filename); err != nil {
		return errResult(err.Error())
	}
	return okResult(true)
}

// ---- Bulk / maintenance ----

func opDump(c *command) []byte {
	if c.Output == "" {
		return errResult("missing output path")
	}
	e, err := getEntry(c.Collection)
	if err != nil {
		return errResult(err.Error())
	}
	f, err := os.Create(c.Output)
	if err != nil {
		return errResult(fmt.Sprintf("creating %q: %s", c.Output, err))
	}
	defer f.Close()
	if err := e.c.Dump(f); err != nil {
		return errResult(err.Error())
	}
	return okResult(true)
}

func opLoad(c *command) []byte {
	if c.Input == "" {
		return errResult("missing input path")
	}
	e, err := getEntry(c.Collection)
	if err != nil {
		return errResult(err.Error())
	}
	f, err := os.Open(c.Input)
	if err != nil {
		return errResult(fmt.Sprintf("opening %q: %s", c.Input, err))
	}
	defer f.Close()
	if err := e.c.Load(f, c.Overwrite, 0); err != nil {
		return errResult(err.Error())
	}
	return okResult(true)
}

func opCheck(c *command) []byte {
	if c.Collection == "" {
		return errResult("missing collection")
	}
	if err := dataset.Analyzer(c.Collection, false); err != nil {
		return errResult(err.Error())
	}
	return okResult(true)
}

func opRepair(c *command) []byte {
	if c.Collection == "" {
		return errResult("missing collection")
	}
	if err := dataset.Repair(c.Collection, false); err != nil {
		return errResult(err.Error())
	}
	return okResult(true)
}

func main() {}
