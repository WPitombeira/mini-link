# Mini-Link

Mini-Link is a LittleLink-style link-in-bio page written in Go. It is designed to be small, fast, cache-friendly, and easy to host anywhere.

It uses only the Go standard library. The rendered page has inline CSS and inline SVG icons, so the browser can load it with a single HTML request when you use built-in or inline custom icons.

## Features

- Go HTTP server and static exporter
- JSON, small YAML subset, and env configuration
- strong ETag, Last-Modified, and CDN-friendly Cache-Control headers
- precompiled inline SVG icon catalog
- user-defined custom icons with inline SVG/path support
- avatar image support plus export-time favicon generation
- selectable `classic`, `glass`, and `terminal` templates
- native dropdown groups without JavaScript
- SEO-ready canonical, social metadata, JSON-LD, robots, sitemap, web manifest, favicon, and `llms.txt` output
- Docker, Docker Compose, Cloudflare Pages, Cloudflare Workers/static assets, Vercel, and manual deployment docs
- MIT license

## Quick Start

```bash
go test ./...
go run ./cmd/minilink serve -config examples/mini-link.yaml
```

Open `http://localhost:8080`.

Build a tiny binary:

```bash
go build -trimpath -ldflags="-s -w" -o bin/minilink ./cmd/minilink
```

Export static files:

```bash
go run ./cmd/minilink export -config examples/mini-link.yaml -out dist
```

## Configuration

Use the format that fits your deployment:

- JSON: `examples/mini-link.json`
- YAML: `examples/mini-link.yaml`
- env file: `examples/mini-link.env.example`

See [docs/configuration.md](docs/configuration.md).

Extra examples:

- `examples/minimal.yaml`
- `examples/dropdowns.yaml`
- `examples/dropdowns.json`
- `examples/custom-icons.yaml`
- `examples/custom-icons.json`
- `examples/cloudflare/wrangler.toml`
- `examples/cloudflare/pages.toml`

## Templates

Set `template` to `classic`, `glass`, or `terminal`. See [docs/templates.md](docs/templates.md).

## SVG Icons

Generate the browser-viewable icon guide:

```bash
go run ./cmd/minilink icons -out docs/icons.html
```

Then open `docs/icons.html`. See [docs/icons.md](docs/icons.md).

External icon URLs are supported with `icon_url`, but they add browser requests and require a looser image CSP. Use inline custom icons when performance is the priority.

## Favicons and Avatars

Set `avatar_url` to show a profile image. Static export can also download `avatar_url`, `favicon.source_url`, or `favicon.source_path`, then write optimized favicon files into `dist/`.

If no favicon source is configured, Mini-Link creates a small SVG favicon from the profile initials. Optional S3-compatible upload is available for S3 and Cloudflare R2; for Cloudflare Pages and Workers Static Assets, the default recommendation is to keep generated assets in `dist/` and let Cloudflare serve/cache them with the site.

## Deployment

See [docs/deployment.md](docs/deployment.md) for Docker, Docker Compose, manual Go builds, Cloudflare Pages, Cloudflare Workers/static assets, and Vercel. See [docs/cloudflare.md](docs/cloudflare.md) for a Cloudflare-specific tutorial.

## SEO

See [docs/seo.md](docs/seo.md) for generated SEO files, AI-readiness output, and the Lighthouse validation workflow.

## License

MIT. See [LICENSE](LICENSE).
