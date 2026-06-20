package render

import (
	"strings"
	"testing"

	"github.com/WPitombeira/mini-link/internal/config"
	"github.com/WPitombeira/mini-link/internal/icons"
)

func TestPageRendersInlineSVGAndEscapesText(t *testing.T) {
	cfg := config.Default()
	cfg.Name = `<Mini>`
	cfg.BaseURL = "https://example.com"
	cfg.Links = []config.Link{{Title: "GitHub", URL: "https://github.com/WPitombeira/mini-link", Icon: "github"}}
	page, err := Page(cfg)
	if err != nil {
		t.Fatal(err)
	}
	html := string(page)
	if !strings.Contains(html, "&lt;Mini&gt;") {
		t.Fatal("profile name was not escaped")
	}
	if !strings.Contains(html, `<svg viewBox="0 0 24 24"`) {
		t.Fatal("expected inline svg")
	}
	if !strings.Contains(html, `<link rel="canonical" href="https://example.com/">`) {
		t.Fatal("missing canonical")
	}
	if !strings.Contains(html, `<meta property="og:title"`) || !strings.Contains(html, `<meta name="twitter:card"`) {
		t.Fatal("missing social metadata")
	}
	if !strings.Contains(html, `<script type="application/ld+json">`) || !strings.Contains(html, `"@type":"Person"`) {
		t.Fatal("missing person json-ld")
	}
}

func TestIconCatalogIncludesKeys(t *testing.T) {
	html, err := IconCatalog(icons.All())
	if err != nil {
		t.Fatal(err)
	}
	page := string(html)
	if !strings.Contains(page, "Mini-Link SVG Icon Catalog") {
		t.Fatal("missing catalog title")
	}
	if !strings.Contains(page, `<svg viewBox="0 0 24 24"`) {
		t.Fatal("expected rendered svg markup")
	}
	if !strings.Contains(page, "<code>github</code>") {
		t.Fatal("missing github icon key")
	}
}
