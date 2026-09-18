
Action Items
============

X
: completed

D
: declined, decided not to implement

Bugs
----

- [ ] A failed `query` is indistinguishable from an empty result, end to end.
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

- [ ] Add a YAML attribute for application config. This would allow me to include non-datasetd configuration for use with other middleware.

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
