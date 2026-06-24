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

func TestBuildStaticAssets(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "avatar.png")
	if err := os.WriteFile(source, []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}, 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := config.Default()
	cfg.StaticAssets = []config.StaticAsset{
		{SourcePath: source, OutputPath: "assets/avatar.png"},
	}

	assets, err := buildStaticAssets(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(assets) != 1 {
		t.Fatalf("assets = %d", len(assets))
	}
	if assets[0].Name != "assets/avatar.png" {
		t.Fatalf("asset name = %q", assets[0].Name)
	}
	if assets[0].ContentType != "image/png" {
		t.Fatalf("content type = %q", assets[0].ContentType)
	}
}
