package api

import (
	"fmt"
	"net/http"
)

func ApiVersion(w http.ResponseWriter, req *http.Request, api *API, verb string, options []string) {
	fmt.Fprintf(w, "%s %s", api.AppName, api.Version)
}

func ApiKeys(w http.ResponseWriter, req *http.Request, api *API, verb string, options []string) {
	http.Error(w, "ApiKeys() not implemented", http.StatusNotImplemented)
}

func ApiCreate(w http.ResponseWriter, req *http.Request, api *API, verb string, options []string) {
	http.Error(w, "ApiCreate() not implemented", http.StatusNotImplemented)
}

func ApiRead(w http.ResponseWriter, req *http.Request, api *API, verb string, options []string) {
	http.Error(w, "ApiRead() not implemented", http.StatusNotImplemented)
}

func ApiUpdate(w http.ResponseWriter, req *http.Request, api *API, verb string, options []string) {
	http.Error(w, "ApiUpdate() not implemented", http.StatusNotImplemented)
}

func ApiDelete(w http.ResponseWriter, req *http.Request, api *API, verb string, options []string) {
	http.Error(w, "ApiDelete() not implemented", http.StatusNotImplemented)
}
