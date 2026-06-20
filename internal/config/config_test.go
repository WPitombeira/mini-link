package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadYAML(t *testing.T) {
	cfg, err := Load(filepath.Join("..", "..", "examples", "mini-link.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Name != "WPitombeira" {
		t.Fatalf("name = %q", cfg.Name)
	}
	if len(cfg.Links) != 6 {
		t.Fatalf("links = %d", len(cfg.Links))
	}
	if !cfg.Links[0].Featured {
		t.Fatal("first link should be featured")
	}
	if cfg.Template != "glass" {
		t.Fatalf("template = %q", cfg.Template)
	}
}

func TestLoadJSON(t *testing.T) {
	cfg, err := Load(filepath.Join("..", "..", "examples", "mini-link.json"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Accent != "#0f766e" {
		t.Fatalf("accent = %q", cfg.Accent)
	}
}

func TestLoadDotEnv(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(`MINI_LINK_NAME=Env Person
MINI_LINK_BIO=Env Bio
MINI_LINK_LINK_1_TITLE=Website
MINI_LINK_LINK_1_URL=https://example.com
MINI_LINK_LINK_1_ICON=globe
`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Name != "Env Person" {
		t.Fatalf("name = %q", cfg.Name)
	}
	if cfg.Links[0].Icon != "globe" {
		t.Fatalf("icon = %q", cfg.Links[0].Icon)
	}
}

func TestLoadRejectsUnknownIcon(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte(`{
  "name": "Bad Icon",
  "links": [
    {
      "title": "Website",
      "url": "https://example.com",
      "icon": "missing-brand"
    }
  ]
}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("expected unknown icon error")
	}
}

func TestLoadRejectsLowContrastAccent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte(`{
  "name": "Low Contrast",
  "accent": "#ffffff",
  "links": [
    {
      "title": "Website",
      "url": "https://example.com",
      "icon": "globe"
    }
  ]
}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("expected low contrast accent error")
	}
}
