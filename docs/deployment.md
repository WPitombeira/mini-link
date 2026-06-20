# Deployment

Mini-Link can run as a tiny Go HTTP server or export a static site. Static export is the preferred path for edge platforms because the generated page is a single `index.html` with inline CSS and SVG icons.

## Manual Build

```bash
go test ./...
go build -trimpath -ldflags="-s -w" -o bin/minilink ./cmd/minilink
./bin/minilink serve -addr :8080 -config examples/mini-link.yaml
```

Export static files:

```bash
./bin/minilink export -config examples/mini-link.yaml -out dist
```

Static export writes:

- `index.html`
- `robots.txt`
- `sitemap.xml`
- `_headers` for Cloudflare Pages and compatible hosts
- `vercel.json` for Vercel header configuration

## Docker

```bash
docker build -t mini-link .
docker run --rm -p 8080:8080 -v "$PWD/examples/mini-link.yaml:/mini-link.yaml:ro" mini-link
```

## Docker Compose

```bash
docker compose up --build
```

Edit `examples/mini-link.yaml` or mount your own config at `/mini-link.yaml`.

## Cloudflare Pages

Use static export:

```bash
go run ./cmd/minilink export -config examples/mini-link.yaml -out dist
```

In Cloudflare Pages:

- build command: `go run ./cmd/minilink export -config examples/mini-link.yaml -out dist`
- output directory: `dist`

The export writes `dist/_headers` with cache, content type, and security headers Cloudflare Pages understands.

## Cloudflare Workers

Workers do not run Go binaries directly. Deploy Mini-Link as prebuilt static assets with a Worker asset binding or Pages.

With Wrangler static assets:

```bash
go run ./cmd/minilink export -config examples/mini-link.yaml -out dist
```

Create `wrangler.toml`:

```toml
name = "mini-link"
compatibility_date = "2026-06-20"
assets = { directory = "./dist" }
```

Then deploy:

```bash
npx wrangler deploy
```

## Vercel

Use a static project:

```bash
go run ./cmd/minilink export -config examples/mini-link.yaml -out dist
```

Set:

- build command: `go run ./cmd/minilink export -config examples/mini-link.yaml -out dist`
- output directory: `dist`

The export writes `dist/vercel.json` with cache, content type, and security headers.

## Caching

Server mode sends:

- strong `ETag` based on rendered HTML
- `Last-Modified`
- `Cache-Control: public, max-age=300, s-maxage=86400, stale-while-revalidate=604800` by default

Static export writes equivalent platform headers. Increase `cache_seconds` when config changes are infrequent.

## SEO and Security Headers

Mini-Link emits canonical, Open Graph, Twitter Card, and JSON-LD metadata on the page. Server mode and static export also provide `robots.txt` and `sitemap.xml` when `base_url` is configured.

The default security headers are intentionally strict because Mini-Link does not need client JavaScript, frames, forms, camera, microphone, geolocation, or payment APIs.
