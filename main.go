package main

/// diskey: Distributed in-memory key-value storage

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	DefaultPort  string = "6660"
	DispatchMode string = "dispatch"
	StorageMode  string = "storage"
)

var _port string = DefaultPort

func main() {
	mode := DispatchMode
	if len(os.Args) > 1 && (os.Args[1] == DispatchMode || os.Args[1] == StorageMode) {
		mode = os.Args[1]
	}
	if len(os.Args) > 2 {
		_port = os.Args[2]
	}
	fmt.Println("yo", _port, mode)

	if mode == DispatchMode {
		http.HandleFunc("/", handleDispatch)
	} else {
		http.HandleFunc("/", handleStore)
	}
	http.ListenAndServe(":"+_port, nil)
}

func handleDispatch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		key := strings.Trim(r.URL.String(), "/")

		instance := getInstance(key)
		instanceUrl, err := url.Parse("http://localhost" + instance)
		if err != nil {
			serverError(w, http.StatusInternalServerError, err.Error())
			return
		}
		url, err := url.JoinPath(instanceUrl.String(), key)
		if err != nil {
			serverError(w, http.StatusInternalServerError, err.Error())
			return
		}

		resp, err := http.Get(url)
		if err != nil {
			serverError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if resp.StatusCode != http.StatusOK {
			serverError(w, http.StatusInternalServerError, resp.Status)
			return
		}
		defer resp.Body.Close()
		io.Copy(w, resp.Body)
		return
	case http.MethodPost:
		key := strings.Trim(r.URL.String(), "/")

		instance := getInstance(key)
		instanceUrl, err := url.Parse("http://localhost" + instance)
		if err != nil {
			serverError(w, http.StatusInternalServerError, err.Error())
			return
		}
		url, err := url.JoinPath(instanceUrl.String(), key)
		if err != nil {
			serverError(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer r.Body.Close()
		req, err := http.NewRequest(http.MethodPost, url, r.Body)
		if err != nil {
			serverError(w, http.StatusInternalServerError, err.Error())
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			serverError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if resp.StatusCode != http.StatusOK {
			serverError(w, http.StatusInternalServerError, "NOT 200OK: "+resp.Status)
			return
		}
		defer resp.Body.Close()
		io.Copy(w, resp.Body)
		return
	}

	serverError(w, http.StatusBadRequest, "Bad Request")
}

func handleStore(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		key := strings.Trim(r.URL.String(), "/")
		v, err := getStored(key)
		if err != nil {
			serverError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if v == "" {
			serverError(w, http.StatusNotFound, "Not Found")
			return
		}
		respond(w, http.StatusOK, v)
		return
	case http.MethodPost:
		key := strings.Trim(r.URL.String(), "/")
		defer r.Body.Close()
		value, err := ioutil.ReadAll(r.Body)
		if err != nil {
			serverError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if err := setStored(key, string(value)); err != nil {
			serverError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respond(w, http.StatusOK, "stored")
		return
	}

	serverError(w, http.StatusBadRequest, "Bad Request")
}

func serverError(w http.ResponseWriter, code int, msg string) {
	respond(w, code, ResponseError{Msg: msg})
}

func respond(w http.ResponseWriter, code int, msg any) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(msg)
}

type ResponseError struct{ Msg string }

var _instances = []string{
	":6661",
	":6662",
	":6663",
	":6664",
}

func getInstance(key string) string {
	keyspaces := []struct {
		min, max int
	}{
		{min: int('0'), max: int('9')},
		{min: int('A'), max: int('Z')},
		{min: int('a'), max: int('z')},
	}
	for _, keyspace := range keyspaces {
		test := int(key[0])
		if test >= keyspace.min && test <= keyspace.max {
			total := keyspace.max - keyspace.min + 1
			test -= keyspace.min + 1
			stride := total / len(_instances)
			idx := (test - 1) / stride
			return _instances[idx]
		}
	}
	panic("GTFO")
}

var _storage map[string]string = make(map[string]string)

func getStored(key string) (string, error) {
	return _storage[key], nil
}

func setStored(key, value string) error {
	_storage[key] = _port + value
	return nil
}
