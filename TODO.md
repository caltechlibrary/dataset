
Action Items
============

- [X] Done for v2.5.2. Drop the 32-bit Raspberry Pi build (`Linux-armv7l`, `GOARM=7`) from the Makefile and any release check the next time this module is worked on (workspace DR-0018, proposed 2026-10-26). 64-bit Raspberry Pi OS stays covered by `Linux-aarch64`. Already-published releases are left as they are.

X
: completed

D
: declined, decided not to implement

Bugs
----

- [ ] A failed `query` is indistinguishable from an empty result, end to end.
      **Server half done in v2.5.2** (reasons in the body, `[]` documented in `datasetd_api.5.md`, and
      `strict_status`). Note the "fourth argument" below is really `errorRedirect`, the fifth parameter; the
      fix added `statusIsErrorReason`/`apiError`. **Still open: the client half**, `ts_dataset`'s `Dataset`
      class, planned for ts_dataset 2.5.2 behind a `strictStatus` option (plan in
      `agents/projects/dataset/plans/datasetd_http_status_v2_5_2_plan.md`, Phase 4).
      Found 2026-09-18 while building a COLD report against `rdm_requests.ds`.

      **Server half.** `api_routes.go` answers four distinct query failures with
      the identical opaque body — `statusIsError(..., http.StatusBadRequest, "")`
      at lines 215, 224, 230 and 236: missing query name, collection not found,
      no queries defined for the collection, and undefined query name. The real
      cause goes to the server log only, so a client cannot tell an operator
      *which* of the four happened, and a caller with no access to the log gets
      nothing.

      The mechanism to fix this already exists in the same function: lines 266,
      292, 439 and 477 pass a real reason as `statusIsError`'s fourth argument.
      The four query-routing cases just pass `""`. It is also inconsistent with
      datasetd's object-write path, which supplies a reason in the body or an
      `x-validation-errors` header — a third-party consumer already depends on
      that difference (`content-dashboard/ds-client.ts`, and see its comment
      above `writeObject`).

      **Client half, and the reason this bites.** `ts_dataset`'s
      `Dataset.query()` (`dataset.ts` ~line 413) returns `undefined` for any
      non-ok response and drops the cause. A genuinely empty result returns
      `[]` with HTTP 200. Both verified against a running `datasetd`: a bogus
      query name gives 400, and a parameterised query matching nothing gives
      200 with `[]`. So `undefined` means *failed* and `[]` means *empty*, but
      nothing in the API's shape or documentation says so, and the natural
      `results ?? []` idiom silently converts a failure into a successful empty
      answer. `DatasetApiClient.query()` and the five other `resp.ok` sites in
      that file (~263, 286, 314, 346, 374) have the same shape.

      **`ts_dataset` is the outlier, not the standard.** Two of the three
      independent clients in this ecosystem already fail loudly:
      `dataset/wrappers/typescript/libdataset.ts` throws
      `new DatasetError(resp.error ?? "unknown error")`, and
      `content-dashboard/ds-client.ts` throws on `!res.ok` and reads the body
      for a cause. Whatever is decided, prefer aligning `ts_dataset` with those
      two rather than the reverse.

      **Recommended changes.**

      1. *Server, low risk and additive.* Give the query route the same
         treatment as the write path: pass a real reason as the fourth argument
         to `statusIsError` instead of `""`, distinguishing the four cases.
         Nothing that reads only the status code changes behaviour, so no
         existing client breaks.
      2. *Client, the actual fix.* Make `ts_dataset` throw a `DatasetError`
         carrying the status and body, matching `libdataset.ts`. Returning a
         discriminated result (`{ok, value} | {ok, error}`) is the alternative
         and is more honest for a library, but it changes every call site,
         whereas throwing only changes the sites that were silently wrong.
      3. *Either way, document the contract.* State in `datasetd_api.5.md`
         that an empty result set is `200` with `[]`, so "no rows" and "request
         failed" are documented as different answers rather than inferred.
      4. *Do not "fix" this by having the server return `200` with `[]` on a
         bad query.* That would make the failure permanently invisible and is
         the opposite of what is wanted.

      **Projects to check before changing the client.** Each assumes, or may
      assume, that a falsy/empty return means "no data" rather than "the
      request failed":

      - `cold` — the largest consumer, 20+ `ds.query()`/`ds.read()` call sites,
        most casting the result straight to a type. The worst class is the RDM
        vocabulary generators: `group_vocabulary.ts` checks for `undefined`,
        then prints an **empty vocabulary** to stdout and exits 0, so a failed
        query yields a valid-looking file that would wipe the vocabulary it
        replaces. `people_vocabulary.ts`, `thesis_option_vocabulary.ts` and
        `journal_vocabulary.ts` share the shape. Also
        `generate_country_collaboration_rpt.ts:147` (`return results ?? []`),
        `cold_reports.ts:243` and `:693`, `journals.ts:227`, `subjects.ts:159`,
        `groups.ts:296`, `funders.ts:193`, `utils.ts:162`.
        `generate_technical_reports_rpt.ts` is the one place that already
        distinguishes them, deliberately.
      - `content-dashboard` — **not affected by the client change** (it has its
        own `ds-client.ts` and does not use `ts_dataset`), but it is a direct
        beneficiary of the server change: it already reads the response body to
        report a cause and currently receives only `Bad Request`. It is also
        not our project, so the server change should stay backward compatible
        at the status-code level. Note its `getAllObjects` deliberately skips
        keys whose fetch fails, which is a separate silent-failure choice.
      - Any other `datasetd` consumer outside this workspace. The server change
        is safe for them; the client change is not a concern unless they import
        `ts_dataset`.

      Filed from the COLD side as cold DR-0022's sibling finding; the COLD
      workaround is recorded in
      `agents/projects/cold/plans/technical_reports_report_plan.md`, Phase D.

Next (prep for v2.6)
--------------------

- [X] Done for v2.5.2 as the opt-in `strict_status` YAML setting; default-on stays for the next major version. Review the HTTP status code returned by datasetd to be more effective, see https://mglaman.dev/blog/status-codes-your-api-getting-wrong for example of the topic

      **Measured 2026-11-13** against datasetd 2.5.1 (scratch SQLite collection, one writable and one
      read-only configuration) while building TMS, which is written to the right codes and carries a
      shim for the current behaviour (workspace `tms` DR-0006). What datasetd answers today, and what it
      should answer:

      | Situation | datasetd 2.5.1 | Should be |
      |---|---|---|
      | `GET /object/<missing key>` | `200` with `{}` and no Content-Type | `404` |
      | `DELETE /object/<missing key>` | `200`, `"action":"delete","status":"OK"` | `404` |
      | `POST /object/<existing key>` | `200`, `"action":"updated"`; **silently overwrites** | `409 Conflict`, nothing written |
      | `PUT /object/<missing key>` | `200`, `"action":"updated"`; **creates nothing** | `404` |
      | `POST` invalid JSON | `406` | `400` |
      | `POST` empty body | `406` | `400` |
      | `POST` with `Content-Type: text/plain` | `201`, accepted | `415` (the API docs say JSON is required) |
      | `POST` body over the 1 MiB limit | `400` | `413` |
      | query with a missing parameter | `406` | `400` |
      | query that fails on the user's input (e.g. bad FTS5 syntax) | `406` | `400` (or `422`) |
      | query name not defined | `400`, opaque body (see the Bugs item above) | `404`, with a reason in the body |
      | unknown collection, unknown route | `403` | `404` |
      | method disabled for the collection | `405` | `405`, adding an `Allow` header |
      | create | `201` | `201` (already right) |
      | update, delete of an existing key | `200` | `200` or `204` |

      `406 Not Acceptable` is the code for failed content negotiation (`Accept`), not for malformed input;
      three different failures use it today. The first four rows are the ones that matter most: with them
      a client cannot tell "done" from "nothing happened", and a repeated create destroys the record that
      was there.

      **This requires a change in datasetd.** It changes answers existing clients may depend on (`cold`
      reads `{}` as "no such object"; `ts_dataset` and `content-dashboard` read status codes), so do not
      just flip it. Suggested: an opt-in setting in the YAML (`strict_status: true`, globally or per
      collection) that selects the codes above, default off through v2.x so nothing breaks, default on in
      the next major version; document both in `datasetd_yaml.5.md` and `datasetd_api.5.md` (an empty
      result set stays `200` with `[]`, as the Bugs item above also asks). TMS needs nothing from it to
      work: its client speaks in outcomes (created, conflict, not found, invalid, too large) and treats
      the current behaviour as a removable compatibility layer pinned by contract tests, so the day this
      lands TMS flips one setting and deletes the shims.

- [X] Done for v2.5.2 (`GetCfg` matches the base name; duplicate base names refused). A `dataset:` path in the settings that is absolute makes the collection routes work (keys, objects) but every named query answer `400 Bad Request` with "not found" in the log. Found 2026-11-13 building TMS: the same YAML with `dataset: etd.ds` (relative to datasetd's working directory) works. Either resolve the query map by the collection as configured and as routed, or refuse an absolute path at start with a clear message. Pinned by `tms/datasetd_contract_test.ts`.

- [X] Done for v2.5.2 (the body says which, and `strict_status` gives 404). A named query that finds the collection but not the query name, an absolute-path collection, and an unknown collection all look alike to a client (`400`, `400`, `403`): see the status-code item above.

- [X] Done for v2.5.2 (both spellings accepted, help fixed). `dataset help load` says each JSONL line needs `__key__` and `__object__`, but the installed binary (2.5.1) reads `key` and `object` and rejects the documented names ("missing key or object"). Found 2026-11-13 building the TMS migration, which uses `key`/`object` (as `cold`'s harvest SQL does). Fix the help, or accept both.

- [ ] Decide the future of libdataset. It is deprecated (dataset DR-0004) and no binaries are published after
      v2.5.2; the source, `libdataset.md` and the wrappers in `wrappers/` are still in the tree. Either remove them
      in a later release or replace them with a purpose-built WASM module. py_dataset is archived; if it is revived
      it will be a client of datasetd like ts_dataset, so `wrappers/python` would go.

- [ ] Review the exit codes for dataset and datasetd to provide more expressive results

- [ ] `dataset` and `datasetd` have no `help topics` listing (`dataset help` with no topic answers "Unable to find
      help"). The workspace practice (irdmtools' Makefile) writes one manual page per help topic from
      `NAME -help topics` and `NAME -help TOPIC`; `release-review` also asks for a help index. Add the listing,
      then the per-topic pages in the Makefile's program rule (see the NOTE there). Raised 2026-11-13 during the
      v2.5.2 release prep; the author has not yet ruled on releasing without it.

- [ ] 66 generated `.html` files are committed in this repository (the project pages, `docs/` and `how-to/`).
      CLAUDE.md says not to commit generated content and that websites are built in CI. Untrack them in a commit
      of their own, and check what the site build needs first. The six manual page `.html` files were removed
      with the move of their sources to `docs/` in v2.5.2.

- [ ] `codemeta.json` `softwareRequirements` still lists "CMTools >= 0.0.43". `make generate` and `make release`
      still run `cmt`; drop the entry when codemetatools replaces it (CLAUDE.md: CMTools is no longer a dependency).


- [ ] Add a YAML attribute for application config. This would allow me to include non-datasetd configuration for use with other middleware.

- [X] Review SQLite3 driver: replace `github.com/ncruces/go-sqlite3` (WASM-based,
      prints spurious stderr warning on every invocation) with `github.com/glebarez/go-sqlite`
      (pure-Go, no embed overhead). Changed `Sqlite3DriverName` from `"sqlite3"` to `"sqlite"`
      in `sqlstore.go`; removed ncruces blank imports from `sqlstore.go` and dropped the sqlite
      import from `libdataset/libdataset.go` entirely (WASM/libdataset deferred).

- [ ] datasetd's own request log must carry the fields logagent reads
      (added 2026-10-07 from the logagent work). Today `requestLogger` and
      `responseLogger` in `api.go` (lines ~91-111) write two separate
      `log.Printf` lines per request, "Request: ..." and "Response: ...",
      with method, path, `r.RemoteAddr`, user agent and the whole query map,
      but no time-to-serve, no response size, no status on the request line,
      and no way to tell the real client from a proxy. logagent's tiers need,
      per request on one line: client address (honouring a trusted
      `X-Forwarded-For` or `CF-Connecting-IP` when datasetd sits behind a
      proxy, and the peer address as well), time, request line, status, body
      bytes, user agent and request time. The authoritative list, with each
      field's purpose, is the field table in the logagent repository,
      `internal/fields/fields.json`, documented in `logagent-fields(5)`. The
      upstream response time field does not apply to datasetd. Also decide
      what to do with the query string: logagent never keeps one (logagent
      DR-0002), and datasetd currently logs it in full, which can include
      search terms a reader typed. Needs a design note first (log line shape,
      whether it is on by default or a `datasetd.yaml` setting, how the
      trusted-proxy list is configured) and tests before code.
- [ ] Explore making datasetd's log output compatible with what NginX writes,
      not only Apache 2 (added 2026-10-07). "Apache compatible" usually means
      the combined log format; NginX's `combined` is the same shape, but a
      real NginX deployment (for example CaltechAUTHORS) appends labelled
      extended fields (`rt=`, `urt=`, `cf_ray=`, `peer=` and so on), so a
      reader written for one server's format may not read the other's. Settle
      one of: emit exactly the stock `combined` line; emit `combined` plus the
      same `label=value` extension logagent generates for nginx
      (`logagent-fields(5)`, "THE GENERATED FORMAT"); or let `datasetd.yaml`
      name a format. Check the result against logagent's reader once that
      exists, and against the existing Apache and NginX documentation in
      `docs/datasetd.md`, which today says only that datasetd can sit behind
      either server.

Someday, Maybe
--------------

- [ ] Provide a TUI based on experience with clasm, each verb becomes a top level menu item, subsequent args are either sub menus or leaf points for actions. We can use Charm and Huh to implement much of the TUI as we're using in the clasm project.
- [ ] dsbagit would generate a "BagIt" bag for preservation of collection
      objects
- [ ] OAI-PMH importer to prototype iiif service based on Islandora
      content driven by a dataset collection
- [ ] Implement an integrated a web UI for managing dataset collections and their data structures
  - [ ] Form pages could be expressed in Markdown+YAML for forms and embedded in the datasetd settings YAML file
    - See my notes on my text oreinted web experiment, yaml2webform.go
    - Forms could be render into the htdocs auto-magically saving development effort
    - The same forms could then be used server side for validation based on descriptors and JavaScript converted to WASM code
