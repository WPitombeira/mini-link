package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestWriteHTMLCacheValidation(t *testing.T) {
	body := []byte("<!doctype html><title>x</title>")
	etag := strongETag(body)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	writeHTML(rec, req, body, etag, time.Unix(10, 0).UTC(), 300)
	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d", rec.Code)
	}
	if got := rec.Header().Get("ETag"); got != etag {
		t.Fatalf("etag = %q", got)
	}
	if !strings.Contains(rec.Header().Get("Cache-Control"), "stale-while-revalidate") {
		t.Fatal("missing stale-while-revalidate")
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("If-None-Match", etag)
	rec = httptest.NewRecorder()
	writeHTML(rec, req, body, etag, time.Unix(10, 0).UTC(), 300)
	if rec.Code != http.StatusNotModified {
		t.Fatalf("code = %d", rec.Code)
	}
}
