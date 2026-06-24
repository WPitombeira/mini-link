package main

import (
	"context"
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
	serveAssets := render.DefaultAssets(cfg)
	servedAssetFiles := []assetFile{{Name: "favicon.svg", ContentType: "image/svg+xml", Body: []byte(initialsFaviconSVG(cfg))}}
	if cfg.Favicon.SourcePath != "" {
		favicon, err := buildFavicon(context.Background(), cfg)
		if err != nil {
			return err
		}
		serveAssets = render.Assets{Favicons: favicon.Favicons}
		servedAssetFiles = favicon.Assets
	}
	page, err := render.PageWithAssets(cfg, serveAssets)
	if err != nil {
		return err
	}
	etag := strongETag(page)
	started := time.Now().UTC()
	allowExternalImages := config.HasExternalMedia(cfg)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		writeHTML(w, r, page, etag, started, cfg.CacheSeconds, allowExternalImages)
	})
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		writeText(w, []byte(robotsTXT(cfg)), "text/plain; charset=utf-8", cfg.CacheSeconds, allowExternalImages)
	})
	mux.HandleFunc("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		writeText(w, []byte(sitemapXML(cfg, started)), "application/xml; charset=utf-8", cfg.CacheSeconds, allowExternalImages)
	})
	mux.HandleFunc("/llms.txt", func(w http.ResponseWriter, r *http.Request) {
		writeText(w, []byte(llmsTXT(cfg)), "text/markdown; charset=utf-8", cfg.CacheSeconds, allowExternalImages)
	})
	for _, asset := range servedAssetFiles {
		asset := asset
		mux.HandleFunc("/"+asset.Name, func(w http.ResponseWriter, r *http.Request) {
			writeText(w, asset.Body, asset.ContentType, cfg.CacheSeconds, allowExternalImages)
		})
	}
	mux.HandleFunc("/site.webmanifest", func(w http.ResponseWriter, r *http.Request) {
		writeText(w, []byte(siteManifest(cfg, serveAssets.Favicons)), "application/manifest+json; charset=utf-8", cfg.CacheSeconds, allowExternalImages)
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
	favicon, err := buildFavicon(context.Background(), cfg)
	if err != nil {
		return err
	}
	faviconLinks := favicon.Favicons
	if cfg.AssetUpload.Provider != "" {
		uploadedLinks, err := uploadAssets(context.Background(), cfg.AssetUpload, favicon.Assets)
		if err != nil {
			return err
		}
		faviconLinks = uploadedLinks
	}
	page, err := render.PageWithAssets(cfg, render.Assets{Favicons: faviconLinks})
	if err != nil {
		return err
	}
	if err := os.MkdirAll(*out, 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(*out, "index.html"), page, 0o644); err != nil {
		return err
	}
	if err := writeAssets(*out, favicon.Assets); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(*out, "site.webmanifest"), []byte(siteManifest(cfg, faviconLinks)), 0o644); err != nil {
		return err
	}
	allowExternalImages := config.HasExternalMedia(cfg) || cfg.AssetUpload.Provider != ""
	if err := os.WriteFile(filepath.Join(*out, "_headers"), []byte(headersFile(cfg.CacheSeconds, allowExternalImages)), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(*out, "vercel.json"), []byte(vercelConfig(cfg.CacheSeconds, allowExternalImages)), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(*out, "robots.txt"), []byte(robotsTXT(cfg)), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(*out, "sitemap.xml"), []byte(sitemapXML(cfg, time.Now().UTC())), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(*out, "llms.txt"), []byte(llmsTXT(cfg)), 0o644); err != nil {
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
	fmt.Printf("ok: %s with %d links\n", cfg.Name, config.CountLinks(cfg.Links))
	return nil
}

func writeHTML(w http.ResponseWriter, r *http.Request, page []byte, etag string, modified time.Time, ttl int, allowExternalImages bool) {
	w.Header().Set("Cache-Control", cacheControl(ttl))
	w.Header().Set("ETag", etag)
	w.Header().Set("Last-Modified", modified.Format(http.TimeFormat))
	securityHeaders(w.Header(), allowExternalImages)
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
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(page)
}

func writeText(w http.ResponseWriter, body []byte, contentType string, ttl int, allowExternalImages bool) {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", cacheControl(ttl))
	securityHeaders(w.Header(), allowExternalImages)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
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

func headersFile(ttl int, allowExternalImages bool) string {
	return "/*\n" +
		"  Cache-Control: " + cacheControl(ttl) + "\n" +
		"  Content-Security-Policy: " + contentSecurityPolicy(allowExternalImages) + "\n" +
		"  X-Content-Type-Options: nosniff\n" +
		"  Referrer-Policy: strict-origin-when-cross-origin\n" +
		"  Permissions-Policy: geolocation=(), microphone=(), camera=(), payment=()\n"
}

func vercelConfig(ttl int, allowExternalImages bool) string {
	return fmt.Sprintf(`{
  "headers": [
    {
      "source": "/(.*)",
      "headers": [
        { "key": "Cache-Control", "value": %q },
        { "key": "Content-Security-Policy", "value": %q },
        { "key": "X-Content-Type-Options", "value": "nosniff" },
        { "key": "Referrer-Policy", "value": "strict-origin-when-cross-origin" },
        { "key": "Permissions-Policy", "value": "geolocation=(), microphone=(), camera=(), payment=()" }
      ]
    }
  ]
}
`, cacheControl(ttl), contentSecurityPolicy(allowExternalImages))
}

func securityHeaders(header http.Header, allowExternalImages bool) {
	header.Set("Content-Security-Policy", contentSecurityPolicy(allowExternalImages))
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	header.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=()")
}

func contentSecurityPolicy(allowExternalImages bool) string {
	policy := "default-src 'none'; connect-src 'self'; style-src 'unsafe-inline'; script-src 'none'; img-src 'self' data:; manifest-src 'self'; base-uri 'none'; form-action 'none'; frame-ancestors 'none'"
	if allowExternalImages {
		policy = strings.Replace(policy, "img-src 'self' data:", "img-src 'self' https: data:", 1)
	}
	return policy
}

func robotsTXT(cfg config.Config) string {
	base := strings.TrimRight(cfg.BaseURL, "/")
	if base == "" {
		return "User-agent: *\nAllow: /\n"
	}
	return "User-agent: *\nAllow: /\nSitemap: " + base + "/sitemap.xml\n"
}

func sitemapXML(cfg config.Config, modified time.Time) string {
	loc := strings.TrimRight(cfg.BaseURL, "/")
	if loc == "" {
		loc = "/"
	}
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>%s</loc>
    <lastmod>%s</lastmod>
    <changefreq>weekly</changefreq>
    <priority>1.0</priority>
  </url>
</urlset>
`, xmlEscape(loc), modified.Format("2006-01-02"))
}

func xmlEscape(value string) string {
	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&apos;")
	return replacer.Replace(value)
}

func llmsTXT(cfg config.Config) string {
	var b strings.Builder
	base := strings.TrimRight(cfg.BaseURL, "/")
	if base == "" {
		base = "/"
	}
	fmt.Fprintf(&b, "# %s\n\n", cfg.Title)
	if cfg.Bio != "" {
		fmt.Fprintf(&b, "> %s\n\n", cfg.Bio)
	}
	fmt.Fprintf(&b, "Mini-Link profile for %s.\n\n", cfg.Name)
	fmt.Fprintf(&b, "## Primary URL\n\n- [%s](%s)\n\n", cfg.Name, base)
	links := publicLinks(cfg.Links)
	if len(links) > 0 {
		b.WriteString("## Public Links\n\n")
		for _, link := range links {
			fmt.Fprintf(&b, "- [%s](%s)\n", link.Title, link.URL)
		}
		b.WriteString("\n")
	}
	b.WriteString("## Machine Readable Files\n\n")
	fmt.Fprintf(&b, "- [Sitemap](%s/sitemap.xml)\n", strings.TrimRight(base, "/"))
	fmt.Fprintf(&b, "- [Robots](%s/robots.txt)\n", strings.TrimRight(base, "/"))
	return b.String()
}

func publicLinks(links []config.Link) []config.Link {
	var out []config.Link
	for _, link := range links {
		if len(link.Links) > 0 {
			out = append(out, publicLinks(link.Links)...)
			continue
		}
		if strings.HasPrefix(link.URL, "http://") || strings.HasPrefix(link.URL, "https://") {
			out = append(out, link)
		}
	}
	return out
}

func envDefault(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
