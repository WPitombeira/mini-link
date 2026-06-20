package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/WPitombeira/mini-link/internal/config"
	"github.com/WPitombeira/mini-link/internal/icons"
	"github.com/WPitombeira/mini-link/internal/render"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	var err error
	switch os.Args[1] {
	case "serve":
		err = serve(os.Args[2:])
	case "export":
		err = export(os.Args[2:])
	case "icons":
		err = writeIcons(os.Args[2:])
	case "validate":
		err = validate(os.Args[2:])
	case "version":
		fmt.Println(version)
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `Mini-Link

Usage:
  minilink serve    -config examples/mini-link.yaml [-addr :8080]
  minilink export   -config examples/mini-link.yaml [-out dist]
  minilink icons    [-out docs/icons.html]
  minilink validate -config examples/mini-link.yaml
  minilink version`)
}

func serve(args []string) error {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	addr := fs.String("addr", envDefault("MINI_LINK_ADDR", ":8080"), "HTTP bind address")
	path := fs.String("config", envDefault("MINI_LINK_CONFIG", ""), "config path (.json, .yaml, .yml, or .env)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := config.Load(*path)
	if err != nil {
		return err
	}
	page, err := render.Page(cfg)
	if err != nil {
		return err
	}
	etag := strongETag(page)
	started := time.Now().UTC()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		writeHTML(w, r, page, etag, started, cfg.CacheSeconds)
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})

	log.Printf("mini-link listening on %s", *addr)
	server := &http.Server{
		Addr:              *addr,
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	return server.ListenAndServe()
}

func export(args []string) error {
	fs := flag.NewFlagSet("export", flag.ExitOnError)
	path := fs.String("config", envDefault("MINI_LINK_CONFIG", ""), "config path (.json, .yaml, .yml, or .env)")
	out := fs.String("out", "dist", "output directory")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := config.Load(*path)
	if err != nil {
		return err
	}
	page, err := render.Page(cfg)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(*out, 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(*out, "index.html"), page, 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(*out, "_headers"), []byte(headersFile(cfg.CacheSeconds)), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(*out, "vercel.json"), []byte(vercelConfig(cfg.CacheSeconds)), 0o644); err != nil {
		return err
	}
	fmt.Printf("exported %s\n", filepath.Clean(*out))
	return nil
}

func writeIcons(args []string) error {
	fs := flag.NewFlagSet("icons", flag.ExitOnError)
	out := fs.String("out", "docs/icons.html", "output HTML path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	html, err := render.IconCatalog(icons.All())
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(*out, html, 0o644); err != nil {
		return err
	}
	fmt.Printf("wrote %s\n", filepath.Clean(*out))
	return nil
}

func validate(args []string) error {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	path := fs.String("config", envDefault("MINI_LINK_CONFIG", ""), "config path (.json, .yaml, .yml, or .env)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	cfg, err := config.Load(*path)
	if err != nil {
		return err
	}
	if _, err := render.Page(cfg); err != nil {
		return err
	}
	fmt.Printf("ok: %s with %d links\n", cfg.Name, len(cfg.Links))
	return nil
}

func writeHTML(w http.ResponseWriter, r *http.Request, page []byte, etag string, modified time.Time, ttl int) {
	if match := r.Header.Get("If-None-Match"); match == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	if since := r.Header.Get("If-Modified-Since"); since != "" {
		if t, err := http.ParseTime(since); err == nil && !modified.After(t) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", cacheControl(ttl))
	w.Header().Set("ETag", etag)
	w.Header().Set("Last-Modified", modified.Format(http.TimeFormat))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(page)
}

func strongETag(body []byte) string {
	sum := sha256.Sum256(body)
	return `"` + hex.EncodeToString(sum[:]) + `"`
}

func cacheControl(ttl int) string {
	if ttl <= 0 {
		ttl = 300
	}
	return fmt.Sprintf("public, max-age=%d, s-maxage=%d, stale-while-revalidate=604800", ttl, ttl*288)
}

func headersFile(ttl int) string {
	return "/*\n  Cache-Control: " + cacheControl(ttl) + "\n  X-Content-Type-Options: nosniff\n"
}

func vercelConfig(ttl int) string {
	return fmt.Sprintf(`{
  "headers": [
    {
      "source": "/(.*)",
      "headers": [
        { "key": "Cache-Control", "value": %q },
        { "key": "X-Content-Type-Options", "value": "nosniff" }
      ]
    }
  ]
}
`, cacheControl(ttl))
}

func envDefault(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
