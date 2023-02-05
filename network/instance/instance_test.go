package instance

import (
	"bytes"
	"disco/network"
	"disco/storage"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestInstance_FetchError(t *testing.T) {
	s := NewInstance(nil, network.NewConfig("", ""))
	key, _ := storage.NewKey("test")
	if _, err := s.Fetch(key); err == nil {
		t.Error("expected error")
	} else if "unknown key: \"test\"" != err.Error() {
		t.Errorf("unexpected error: %q", err)
	}
}

func TestInstance_FetchPut(t *testing.T) {
	s := NewInstance(nil, network.NewConfig("", ""))
	key, _ := storage.NewKey("test")
	s.Put(key, "wat")
	if v, err := s.Fetch(key); err != nil {
		t.Errorf("unexpected error: %q", err)
	} else if v.Value() != "wat" {
		t.Errorf("unexpected value: %q", v)
	}
}

func TestInstance_NoKey(t *testing.T) {
	s := NewInstance(nil, network.NewConfig("", ""))
	w := httptest.NewRecorder()
	lnk, _ := url.Parse("http://localhost/")
	req, _ := http.NewRequest(
		http.MethodGet, lnk.String(), nil)
	req.Header.Add("x-relay-key", s.apiKey.String())
	s.Handle(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected Err400 on no key")
	}
}

func TestInstance_MissingKey(t *testing.T) {
	s := NewInstance(nil, network.NewConfig("", ""))
	w := httptest.NewRecorder()
	lnk, _ := url.Parse("http://localhost/?key=wat")
	req, _ := http.NewRequest(
		http.MethodGet, lnk.String(), nil)
	req.Header.Add("x-relay-key", s.apiKey.String())
	s.Handle(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected Err404 on missing key")
	}

	r := w.Result()
	defer r.Body.Close()
	resp, _ := ioutil.ReadAll(r.Body)
	if string(resp) != `unknown key: "wat"` {
		t.Errorf("unexpected response: %q", resp)
	}
}

func TestInstance_HappyPath(t *testing.T) {
	s := NewInstance(nil, network.NewConfig("", ""))
	expected := "YAY this is the proper value"
	key, _ := storage.NewKey("wat")
	s.Put(key, expected)
	w := httptest.NewRecorder()
	lnk, _ := url.Parse("http://localhost/?key=" + key.String())
	req, _ := http.NewRequest(
		http.MethodGet, lnk.String(), nil)
	req.Header.Add("x-relay-key", s.apiKey.String())
	s.Handle(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200OK with proper key")
	}

	r := w.Result()
	defer r.Body.Close()
	resp, _ := ioutil.ReadAll(r.Body)
	if string(resp) != expected {
		t.Errorf("unexpected response: %q", resp)
	}
}

func TestInstance_HappyPathRoundtrip(t *testing.T) {
	s := NewInstance(nil, network.NewConfig("", ""))
	expected := "YAY this is the proper value"
	lnk, _ := url.Parse("http://localhost/?key=wat")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		http.MethodGet, lnk.String(), nil)
	req.Header.Add("x-relay-key", s.apiKey.String())
	s.Handle(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("should be not found initially: %d", w.Code)
	}

	w = httptest.NewRecorder()
	post, _ := http.NewRequest(
		http.MethodPost,
		lnk.String(),
		bytes.NewBuffer([]byte(expected)),
	)
	post.Header.Add("x-relay-key", s.apiKey.String())
	s.Handle(w, post)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200OK saving the value: %d", w.Code)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(
		http.MethodGet, lnk.String(), nil)
	req.Header.Add("x-relay-key", s.apiKey.String())
	s.Handle(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("should be found now")
	}
	r := w.Result()
	defer r.Body.Close()
	resp, _ := ioutil.ReadAll(r.Body)
	if string(resp) != expected {
		t.Errorf("unexpected response: %q", resp)
	}
}
