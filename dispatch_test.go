package main

import (
	"disco/network"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDispatch_ErrorsOnNilPeers(t *testing.T) {
	d := NewDispatch(nil)
	w := httptest.NewRecorder()
	d.handle(w, nil)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected Err500 on empty peers dispatch")
	}
}

func TestDispatch_ErrorsOnPeersInit(t *testing.T) {
	d := NewDispatch(network.NewPeers())
	w := httptest.NewRecorder()
	d.handle(w, nil)
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
