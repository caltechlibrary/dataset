package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func ApiVersion(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
	fmt.Fprintf(w, "%s %s", api.AppName, api.Version)
}

func ApiCollections(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
	collections := []string{}
	w.Header().Add("Content-Type", "application/json")
	if len(api.CMap) > 0 {
		for k := range api.CMap {
			collections = append(collections, k)
		}
		src, err := json.MarshalIndent(collections, "", "     ")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "%s", src)
		return
	}
	fmt.Fprintf(w, "[]")
}

func ApiKeys(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
	if c, ok := api.CMap[cName]; ok {
		keys, err := c.Keys()
		if err != nil {
			log.Printf("c.Keys() returned error %s", err)
			http.NotFound(w, r)
			return
		}
		src, err := json.MarshalIndent(keys, "", "    ")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// Set header to application/json
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", src)
		return
	}
	fmt.Fprintf(w, "ApiKeys(w, r, api, %q, %q, %s) not implemented", cName, verb, strings.Join(options, " "))
	return
}

func ApiCreate(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
	defer r.Body.Close()
	if len(options) != 1 {
		log.Printf("DEBUG request missing key value")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	key := options[0]

	if c, ok := api.CMap[cName]; ok {
		src, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		o := map[string]interface{}{}
		err = json.Unmarshal(src, &o)
		if err != nil {
			log.Printf("unmarshal error %+v, %s", o, err)
			http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
			return
		}
		if err := c.Create(key, o); err != nil {
			log.Printf("DEBUG c.Create(%q, %s), %s", key, src, err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// Set header to application/json
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok", "key": %q, "action": "created"}`, key)
		return
	}
	http.NotFound(w, r)
	return
}

func ApiRead(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
	if len(options) != 1 {
		log.Printf("DEBUG request missing key value")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	key := options[0]

	if c, ok := api.CMap[cName]; ok {
		o := map[string]interface{}{}
		err := c.Read(key, o)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		src, err := json.MarshalIndent(o, "", "    ")
		if err != nil {
			log.Printf("marshal error %+v, %s", o, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// Set header to application/json
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `%s`, src)
		return
	}
	http.NotFound(w, r)
	return
}

func ApiUpdate(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
	defer r.Body.Close()
	if len(options) != 1 {
		log.Printf("DEBUG request missing key value")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	key := options[0]

	if c, ok := api.CMap[cName]; ok {
		src, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		o := map[string]interface{}{}
		err = json.Unmarshal(src, &o)
		if err != nil {
			log.Printf("unmarshal error %+v, %s", o, err)
			http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
			return
		}
		if err := c.Update(key, o); err != nil {
			log.Printf("DEBUG c.Update(%q, %s), %s", key, src, err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// Set header to application/json
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok", "key": %q, "action": "updated"}`, key)
		return
	}
	http.NotFound(w, r)
	return
}

func ApiDelete(w http.ResponseWriter, r *http.Request, api *API, cName string, verb string, options []string) {
	defer r.Body.Close()
	if len(options) != 1 {
		log.Printf("DEBUG request missing key value")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	key := options[0]

	if c, ok := api.CMap[cName]; ok {
		if err := c.Delete(key); err != nil {
			log.Printf("DEBUG c.Delete(%q), %s", key, err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// Set header to application/json
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok", "key": %q, "action": "deleted"}`, key)
		return
	}
	http.NotFound(w, r)
	return
}
