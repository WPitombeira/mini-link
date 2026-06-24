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
	if CountLinks(cfg.Links) != 7 {
		t.Fatalf("links = %d", CountLinks(cfg.Links))
	}
	if !cfg.Links[0].Featured {
		t.Fatal("first link should be featured")
	}
	if len(cfg.Links[4].Links) != 2 {
		t.Fatalf("dropdown links = %d", len(cfg.Links[4].Links))
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
	if CountLinks(cfg.Links) != 7 {
		t.Fatalf("links = %d", CountLinks(cfg.Links))
	}
}

func TestLoadDotEnv(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(`MINI_LINK_NAME=Env Person
MINI_LINK_BIO=Env Bio
MINI_LINK_LINK_1_TITLE=Website
MINI_LINK_LINK_1_URL=https://example.com
MINI_LINK_LINK_1_ICON_URL=https://cdn.example.com/icon.svg
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
	if cfg.Links[0].IconURL != "https://cdn.example.com/icon.svg" {
		t.Fatalf("icon_url = %q", cfg.Links[0].IconURL)
	}
}

func TestLoadCustomIcons(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "custom.json")
	if err := os.WriteFile(path, []byte(`{
  "name": "Custom Icons",
  "custom_icons": [
    {
      "name": "spark",
      "label": "Spark",
      "view_box": "0 0 24 24",
      "path": "M12 2 15 9 22 12 15 15 12 22 9 15 2 12 9 9Z"
    },
    {
      "name": "remote_logo",
      "label": "Remote Logo",
      "url": "https://cdn.example.com/logo.svg"
    }
  ],
  "links": [
    {
      "title": "Inline",
      "url": "https://example.com",
      "icon": "spark"
    },
    {
      "title": "Remote",
      "url": "https://example.com/remote",
      "icon": "remote_logo"
    },
    {
      "title": "Direct Remote",
      "url": "https://example.com/direct",
      "icon_url": "https://cdn.example.com/direct.svg"
    }
  ]
}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.CustomIcons) != 2 {
		t.Fatalf("custom icons = %d", len(cfg.CustomIcons))
	}
	if !HasExternalIcons(cfg) {
		t.Fatal("expected external icons")
	}
}

func TestLoadYAMLCustomIcons(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "custom.yaml")
	if err := os.WriteFile(path, []byte(`name: YAML Custom Icons
custom_icons:
  - name: spark
    label: Spark
    view_box: "0 0 24 24"
    path: "M12 2 15 9 22 12 15 15 12 22 9 15 2 12 9 9Z"
links:
  - title: Inline
    url: https://example.com
    icon: spark
`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Links[0].Icon != "spark" {
		t.Fatalf("icon = %q", cfg.Links[0].Icon)
	}
}

func TestLoadDropdownExamples(t *testing.T) {
	for _, name := range []string{"dropdowns.yaml", "dropdowns.json"} {
		cfg, err := Load(filepath.Join("..", "..", "examples", name))
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if CountLinks(cfg.Links) != 8 {
			t.Fatalf("%s links = %d", name, CountLinks(cfg.Links))
		}
		if !cfg.Links[1].Open {
			t.Fatalf("%s work dropdown should be open", name)
		}
	}
}

func TestLoadRejectsURLAndNestedLinks(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte(`{
  "name": "Bad Dropdown",
  "links": [
    {
      "title": "Bad",
      "url": "https://example.com",
      "links": [
        {
          "title": "Website",
          "url": "https://example.com",
          "icon": "globe"
        }
      ]
    }
  ]
}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("expected url plus nested links error")
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

func TestLoadRejectsUnsafeCustomSVG(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte(`{
  "name": "Bad SVG",
  "custom_icons": [
    {
      "name": "bad",
      "svg": "<svg viewBox=\"0 0 24 24\"><script>alert(1)</script></svg>"
    }
  ],
  "links": [
    {
      "title": "Website",
      "url": "https://example.com",
      "icon": "bad"
    }
  ]
}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("expected unsafe svg error")
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
