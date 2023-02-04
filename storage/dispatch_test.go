package storage

import (
	"disco/network"
	"disco/store"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestDispatch_ErrorsOnNilPeers(t *testing.T) {
	d := NewDispatch(nil, NewStorageConfig("storage", ":6660"))
	w := httptest.NewRecorder()
	d.Handle(w, nil)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected Err500 on empty peers dispatch")
	}
}

func TestDispatch_ErrorsOnPeersInit(t *testing.T) {
	d := NewDispatch(network.NewPeers(), NewStorageConfig("storage", ":6660"))
	w := httptest.NewRecorder()
	d.Handle(w, nil)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected Err500 on empty peers dispatch")
	}

	r := w.Result()
	defer r.Body.Close()
	resp, _ := ioutil.ReadAll(r.Body)
	if string(resp) != network.Init.String() {
		t.Errorf("unexpected response: %q", resp)
	}
}

func TestDispatch_InstanceGetting(t *testing.T) {
	p := network.NewPeers()
	p.Confirm("test1", "test2")

	d := NewDispatch(p, NewStorageConfig("storage", ":6660"))
	suite := map[string]string{
		"AAA": "test1",
		"aaa": "test1",
		"ZZZ": "test2",
		"zzz": "test2",
	}
	for test, want := range suite {
		t.Run(test, func(t *testing.T) {
			key, err := store.NewKey(test)
			if err != nil {
				t.Error(err)
				return
			}
			got := d.getInstance(key)
			if got != want {
				t.Errorf("%q: want %q, got %q",
					test, want, got)
			}
		})
	}
}

func TestDispatch_InstanceUrlGetting(t *testing.T) {
	p := network.NewPeers()
	p.Confirm("test1", "test2")

	d := NewDispatch(p, NewStorageConfig("storage", ":6660"))
	suite := map[string]struct {
		host  string
		query string
	}{
		"AAA": {host: "test1:6660", query: "AAA"},
		"aaa": {host: "test1:6660", query: "aaa"},
		"ZZZ": {host: "test2:6660", query: "ZZZ"},
		"zzz": {host: "test2:6660", query: "zzz"},
	}
	for test, want := range suite {
		t.Run(test, func(t *testing.T) {
			key, err := store.NewKey(test)
			if err != nil {
				t.Error(err)
				return
			}
			got := d.getInstanceURL(key)
			if got.Host != want.host {
				t.Errorf("%q host: want %q, got %q",
					test, want.host, got.Host)
			}
			if got.Query().Get("key") != want.query {
				t.Errorf("%q key: want %q, got %q",
					test, want.query, got.Query().Get("key"))
			}
		})
	}
}

func TestDispatch_NoKey(t *testing.T) {
	p := network.NewPeers()
	p.Confirm("test1", "test2")
	p.SetReady(true)

	d := NewDispatch(p, NewStorageConfig("storage", ":6660"))
	w := httptest.NewRecorder()
	lnk, _ := url.Parse("http://localhost/")
	d.Handle(w, &http.Request{
		URL: lnk,
	})
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected Err400 on empty no key")
	}
}

func TestDispatch_ErrorsWithNoStorageServer(t *testing.T) {
	p := network.NewPeers()
	p.Confirm("test1", "test2")
	p.SetReady(true)

	d := NewDispatch(p, NewStorageConfig("storage", ":6660"))
	d.client = &http.Client{
		Timeout: time.Millisecond,
	}
	w := httptest.NewRecorder()
	lnk, _ := url.Parse("http://whatever-fake-host/?key=AAA")
	d.Handle(w, &http.Request{
		URL:    lnk,
		Method: http.MethodGet,
	})
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected Err500 on transport error")
	}

	r := w.Result()
	defer r.Body.Close()
	resp, _ := ioutil.ReadAll(r.Body)
	if !strings.Contains(string(resp), "deadline exceeded") {
		t.Errorf("unexpected response: %q", resp)
	}
}

func TestDispatch_HappyPath(t *testing.T) {
	handle := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, r.URL.Query().Get("key"))
	}
	srv := httptest.NewServer(http.HandlerFunc(handle))
	defer srv.Close()

	raw, _ := url.Parse(srv.URL)
	split := strings.SplitN(raw.Host, ":", 2)
	host := split[0]
	port := split[1]
	p := network.NewPeers()
	p.Confirm(host, "test2")
	p.SetReady(true)

	d := NewDispatch(p, NewStorageConfig("storage", ":6660"))
	d.storagePort = port
	w := httptest.NewRecorder()
	r := httptest.NewRequest(
		http.MethodGet,
		"/?key=AAA",
		nil)

	d.Handle(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200OK on transport error")
	}

	resp := w.Result()
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if strings.TrimSpace(string(body)) != "AAA" {
		t.Errorf("unexpected response: %q", string(body))
	}
}
