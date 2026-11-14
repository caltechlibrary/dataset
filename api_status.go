package dataset

// api_status.go holds the one place the legacy and strict HTTP status codes
// are chosen between (workspace DR-0020; design brief
// agents/projects/dataset/design/datasetd_http_status_v2_5_2.md).
//
// With `strict_status: true` in the datasetd YAML the object, query and keys
// routes answer with the precise status and a JSON error body. Without it they
// answer as datasetd 2.5.1 did. Every handler passes both codes to apiError, so
// removing the legacy half when strict becomes the default is a mechanical edit.

import (
	"encoding/json"
	"errors"
	"mime"
	"net/http"
	"strings"
)

const (
	mediaJSON = "application/json"
	mediaForm = "application/x-www-form-urlencoded"
)

// strict says whether this service answers with the DR-0020 status codes.
//
// @returns {bool} true when strict_status is set in the settings
// @example
//
//	if api.strict() { ... }
func (api *API) strict() bool {
	return api != nil && api.Settings != nil && api.Settings.StrictStatus
}

// outcomeName gives the DR-0020 outcome for an HTTP status.
//
// @param status {int} the HTTP status
// @returns {string} the outcome name, e.g. "not_found"
// @example
//
//	outcomeName(409) // "conflict"
func outcomeName(status int) string {
	switch {
	case status == http.StatusUnauthorized:
		return "unauthenticated"
	case status == http.StatusForbidden:
		return "forbidden"
	case status == http.StatusNotFound:
		return "not_found"
	case status == http.StatusConflict:
		return "conflict"
	case status == http.StatusRequestEntityTooLarge:
		return "too_large"
	case status == http.StatusTooManyRequests, status == http.StatusBadGateway,
		status == http.StatusServiceUnavailable, status == http.StatusGatewayTimeout:
		return "unavailable"
	case status >= 500:
		return "internal"
	case status >= 400:
		return "invalid"
	}
	return "ok"
}

// writeStrictError writes the DR-0020 JSON error body, never cached.
//
// @param w {http.ResponseWriter} the response
// @param status {int} the HTTP status
// @param message {string} what is wrong, for a person; never a stack trace or a secret
// @example
//
//	writeStrictError(w, 404, "no such record: k1")
func writeStrictError(w http.ResponseWriter, status int, message string) {
	if message == "" {
		message = http.StatusText(status)
	}
	body := map[string]map[string]interface{}{"error": {
		"status":  status,
		"outcome": outcomeName(status),
		"message": message,
	}}
	src, _ := json.Marshal(body)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	w.Write(src)
	w.Write([]byte("\n"))
}

// apiError answers an error with the legacy status, or the strict status when
// strict_status is on. A configured error redirect is honoured in both modes.
//
// @param w {http.ResponseWriter} the response
// @param r {*http.Request} the request
// @param api {*API} the service, which knows the mode
// @param legacy {int} the status datasetd 2.5.1 sent
// @param strict {int} the DR-0020 status
// @param reason {string} why, for the client; the strict body always carries it
// @param redirect {string} the error redirect from the YAML, or ""
// @param legacyReason {bool} true to put the reason in the legacy body too
// @example
//
//	apiError(w, r, api, http.StatusNotAcceptable, http.StatusBadRequest, "request body is not valid JSON", "", false)
func apiError(w http.ResponseWriter, r *http.Request, api *API, legacy, strict int, reason, redirect string, legacyReason bool) {
	if redirect != "" {
		statusIsError(w, r, http.StatusText(legacy), legacy, redirect)
		return
	}
	if api.strict() {
		writeStrictError(w, strict, reason)
		return
	}
	if legacyReason && reason != "" {
		statusIsErrorReason(w, legacy, reason)
		return
	}
	http.Error(w, http.StatusText(legacy), legacy)
}

// mediaType returns the media type of a Content-Type header without its
// parameters, lower case ("application/json; charset=utf-8" is
// "application/json").
//
// @param contentType {string} a Content-Type header value
// @returns {string} the media type, or "" when there is none
// @example
//
//	mediaType("Application/JSON; charset=utf-8") // "application/json"
func mediaType(contentType string) string {
	if strings.TrimSpace(contentType) == "" {
		return ""
	}
	mt, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return strings.ToLower(strings.TrimSpace(contentType))
	}
	return mt
}

// bodyReadStatus is the strict status for a request body that could not be
// read: 413 when it was over the size limit, otherwise 400.
//
// @param err {error} what reading the body returned
// @returns {int} 413 or 400
// @example
//
//	apiError(w, r, api, 400, bodyReadStatus(err), err.Error(), "", false)
func bodyReadStatus(err error) int {
	var tooBig *http.MaxBytesError
	if errors.As(err, &tooBig) {
		return http.StatusRequestEntityTooLarge
	}
	return http.StatusBadRequest
}

// redirectStatus is the status a successful form post redirects with: 303 See
// Other when strict, so a reload does not post again, 302 as datasetd 2.5.1 did.
//
// @param api {*API} the service
// @returns {int} 303 or 302
// @example
//
//	http.Redirect(w, r, to, redirectStatus(api))
func redirectStatus(api *API) int {
	if api.strict() {
		return http.StatusSeeOther
	}
	return http.StatusFound
}
