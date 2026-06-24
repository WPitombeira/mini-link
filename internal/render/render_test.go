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

func TestPageRendersDropdowns(t *testing.T) {
	cfg := config.Default()
	cfg.BaseURL = "https://example.com"
	cfg.Links = []config.Link{
		{
			Title: "Projects",
			Icon:  "briefcase",
			Open:  true,
			Links: []config.Link{
				{Title: "Mini-Link", URL: "https://github.com/WPitombeira/mini-link", Icon: "github"},
				{Title: "Contact", URL: "mailto:hello@example.com", Icon: "mail"},
			},
		},
	}
	page, err := Page(cfg)
	if err != nil {
		t.Fatal(err)
	}
	html := string(page)
	if !strings.Contains(html, `<details class="dropdown" open>`) {
		t.Fatal("missing open dropdown")
	}
	if !strings.Contains(html, `<summary>`) {
		t.Fatal("missing summary")
	}
	if !strings.Contains(html, `href="https://github.com/WPitombeira/mini-link"`) {
		t.Fatal("missing child link")
	}
	if !strings.Contains(html, `"sameAs":["https://github.com/WPitombeira/mini-link"]`) {
		t.Fatal("external child link should be included in schema")
	}
}

func TestPageRendersGoogleAdsTrackingWhenConfigured(t *testing.T) {
	cfg := config.Default()
	cfg.Tracking.GoogleAdsID = "AW-123456789"
	cfg.Links = []config.Link{{Title: "Website", URL: "https://example.com", Icon: "globe"}}
	page, err := Page(cfg)
	if err != nil {
		t.Fatal(err)
	}
	html := string(page)
	if !strings.Contains(html, "https://www.googletagmanager.com/gtag/js?id=AW-123456789") {
		t.Fatal("missing google tag script")
	}
	if !strings.Contains(html, `gtag("config","AW-123456789")`) {
		t.Fatal("missing google ads config")
	}
}

func TestPageRendersCustomAndExternalIcons(t *testing.T) {
	cfg := config.Default()
	cfg.CustomIcons = []config.CustomIcon{
		{
			Name:    "spark",
			Label:   "Spark",
			ViewBox: "0 0 24 24",
			Path:    "M12 2 15 9 22 12 15 15 12 22 9 15 2 12 9 9Z",
		},
		{
			Name: "remote",
			URL:  "https://cdn.example.com/icon.svg",
		},
	}
	cfg.Links = []config.Link{
		{Title: "Inline", URL: "https://example.com", Icon: "spark"},
		{Title: "Remote", URL: "https://example.com/remote", Icon: "remote"},
		{Title: "Direct", URL: "https://example.com/direct", IconURL: "https://cdn.example.com/direct.svg"},
	}
	page, err := Page(cfg)
	if err != nil {
		t.Fatal(err)
	}
	html := string(page)
	if !strings.Contains(html, `<path d="M12 2 15 9 22 12 15 15 12 22 9 15 2 12 9 9Z"`) {
		t.Fatal("missing path-based custom icon")
	}
	if !strings.Contains(html, `<img src="https://cdn.example.com/icon.svg" alt="" loading="eager" decoding="async">`) {
		t.Fatal("missing custom external icon")
	}
	if !strings.Contains(html, `<img src="https://cdn.example.com/direct.svg" alt="" loading="eager" decoding="async">`) {
		t.Fatal("missing direct external icon")
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
