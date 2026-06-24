package main

import (
	"context"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WPitombeira/mini-link/internal/config"
)

func TestBuildFaviconDefaultsToInitialsSVG(t *testing.T) {
	cfg := config.Default()
	result, err := buildFavicon(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Assets) != 1 || result.Assets[0].Name != "favicon.svg" {
		t.Fatalf("assets = %#v", result.Assets)
	}
	if !strings.Contains(string(result.Assets[0].Body), cfg.Avatar) {
		t.Fatal("default favicon should include avatar initials")
	}
}

func TestBuildFaviconProcessesLocalImage(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source.png")
	img := image.NewRGBA(image.Rect(0, 0, 64, 48))
	for y := 0; y < 48; y++ {
		for x := 0; x < 64; x++ {
			img.Set(x, y, color.RGBA{R: 15, G: 118, B: 110, A: 255})
		}
	}
	file, err := os.Create(source)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(file, img); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	cfg := config.Default()
	cfg.Favicon.SourcePath = source
	result, err := buildFavicon(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Assets) != 2 {
		t.Fatalf("assets = %d", len(result.Assets))
	}
	names := map[string]bool{}
	for _, asset := range result.Assets {
		names[asset.Name] = true
		if asset.ContentType != "image/png" {
			t.Fatalf("content type = %q", asset.ContentType)
		}
	}
	if !names["favicon.png"] || !names["apple-touch-icon.png"] {
		t.Fatalf("missing png assets: %#v", names)
	}
}
