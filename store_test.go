package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestStore_FetchError(t *testing.T) {
	s := NewStore()
	if _, err := s.fetch("test"); err == nil {
		t.Error("expected error")
	} else if "unknown key: \"test\"" != err.Error() {
		t.Errorf("unexpected error: %q", err)
	}
}

func TestStore_FetchPut(t *testing.T) {
	s := NewStore()
	s.put("test", "wat")
	if v, err := s.fetch("test"); err != nil {
		t.Errorf("unexpected error: %q", err)
	} else if v != "wat" {
		t.Errorf("unexpected value: %q", v)
	}
}

func TestStore_NoKey(t *testing.T) {
	s := NewStore()
	w := httptest.NewRecorder()
	lnk, _ := url.Parse("http://localhost/")
	s.handle(w, &http.Request{
		URL:    lnk,
		Method: http.MethodGet,
	})
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected Err400 on no key")
	}
}

func TestStore_MissingKey(t *testing.T) {
	s := NewStore()
	w := httptest.NewRecorder()
	lnk, _ := url.Parse("http://localhost/?key=wat")
	s.handle(w, &http.Request{
		URL:    lnk,
		Method: http.MethodGet,
	})
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

func TestStore_HappyPath(t *testing.T) {
	s := NewStore()
	expected := "YAY this is the proper value"
	s.put("wat", expected)
	w := httptest.NewRecorder()
	lnk, _ := url.Parse("http://localhost/?key=wat")
	s.handle(w, &http.Request{
		URL:    lnk,
		Method: http.MethodGet,
	})
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

func TestStore_HappyPathRoundtrip(t *testing.T) {
	s := NewStore()
	expected := "YAY this is the proper value"
	lnk, _ := url.Parse("http://localhost/?key=wat")

	w := httptest.NewRecorder()
	s.handle(w, &http.Request{
		URL:    lnk,
		Method: http.MethodGet,
	})
	if w.Code != http.StatusNotFound {
		t.Errorf("should be not found initially")
	}

	w = httptest.NewRecorder()
	post, _ := http.NewRequest(
		http.MethodPost,
		lnk.String(),
		bytes.NewBuffer([]byte(expected)),
	)
	s.handle(w, post)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200OK saving the value")
	}

	w = httptest.NewRecorder()
	s.handle(w, &http.Request{
		URL:    lnk,
		Method: http.MethodGet,
	})
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
