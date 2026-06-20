package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WPitombeira/mini-link/internal/config"
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
	if rec.Header().Get("Content-Security-Policy") == "" {
		t.Fatal("missing content security policy")
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("If-None-Match", etag)
	rec = httptest.NewRecorder()
	writeHTML(rec, req, body, etag, time.Unix(10, 0).UTC(), 300)
	if rec.Code != http.StatusNotModified {
		t.Fatalf("code = %d", rec.Code)
	}
	if rec.Header().Get("ETag") != etag {
		t.Fatal("304 should include etag")
	}
}

func TestRobotsAndSitemap(t *testing.T) {
	cfg := config.Default()
	cfg.BaseURL = "https://example.com"
	robots := robotsTXT(cfg)
	if !strings.Contains(robots, "Sitemap: https://example.com/sitemap.xml") {
		t.Fatalf("robots = %q", robots)
	}
	sitemap := sitemapXML(cfg, time.Date(2026, 6, 20, 0, 0, 0, 0, time.UTC))
	if !strings.Contains(sitemap, "<loc>https://example.com</loc>") {
		t.Fatalf("sitemap = %q", sitemap)
	}
	if !strings.Contains(sitemap, "<lastmod>2026-06-20</lastmod>") {
		t.Fatalf("sitemap missing lastmod: %q", sitemap)
	}
}
