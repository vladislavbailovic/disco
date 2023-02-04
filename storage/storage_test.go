package storage

import (
	"bytes"
	"disco/store"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestStorage_FetchError(t *testing.T) {
	s := NewStorage(nil)
	key, _ := store.NewKey("test")
	if _, err := s.Fetch(key); err == nil {
		t.Error("expected error")
	} else if "unknown key: \"test\"" != err.Error() {
		t.Errorf("unexpected error: %q", err)
	}
}

func TestStorage_FetchPut(t *testing.T) {
	s := NewStorage(nil)
	key, _ := store.NewKey("test")
	s.Put(key, "wat")
	if v, err := s.Fetch(key); err != nil {
		t.Errorf("unexpected error: %q", err)
	} else if v.Value() != "wat" {
		t.Errorf("unexpected value: %q", v)
	}
}

func TestStorage_NoKey(t *testing.T) {
	s := NewStorage(nil)
	w := httptest.NewRecorder()
	lnk, _ := url.Parse("http://localhost/")
	s.Handle(w, &http.Request{
		URL:    lnk,
		Method: http.MethodGet,
	})
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected Err400 on no key")
	}
}

func TestStorage_MissingKey(t *testing.T) {
	s := NewStorage(nil)
	w := httptest.NewRecorder()
	lnk, _ := url.Parse("http://localhost/?key=wat")
	s.Handle(w, &http.Request{
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

func TestStorage_HappyPath(t *testing.T) {
	s := NewStorage(nil)
	expected := "YAY this is the proper value"
	key, _ := store.NewKey("wat")
	s.Put(key, expected)
	w := httptest.NewRecorder()
	lnk, _ := url.Parse("http://localhost/?key=" + key.String())
	s.Handle(w, &http.Request{
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

func TestStorage_HappyPathRoundtrip(t *testing.T) {
	s := NewStorage(nil)
	expected := "YAY this is the proper value"
	lnk, _ := url.Parse("http://localhost/?key=wat")

	w := httptest.NewRecorder()
	s.Handle(w, &http.Request{
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
	s.Handle(w, post)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200OK saving the value")
	}

	w = httptest.NewRecorder()
	s.Handle(w, &http.Request{
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
