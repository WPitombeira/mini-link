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
- `favicon.svg` or generated PNG favicon files
- `site.webmanifest`
- `robots.txt`
- `sitemap.xml`
- `llms.txt`
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

The export writes `dist/_headers` with cache, content type, and security headers Cloudflare Pages understands. See [cloudflare.md](cloudflare.md) for the full tutorial.

## Cloudflare Workers

Workers do not run Go binaries directly. Deploy Mini-Link as prebuilt static assets with Workers Static Assets or Pages.

With Wrangler static assets:

```bash
go run ./cmd/minilink export -config examples/mini-link.yaml -out dist
```

Create `wrangler.toml`:

```toml
name = "mini-link"
compatibility_date = "2026-06-24"

[assets]
directory = "./dist"
```

Then deploy:

```bash
npx wrangler deploy
```

See [cloudflare.md](cloudflare.md) for the full tutorial, including custom domains and post-deploy checks.

## Vercel

Use a static project:

```bash
go run ./cmd/minilink export -config examples/mini-link.yaml -out dist
```

Set:

- build command: `go run ./cmd/minilink export -config examples/mini-link.yaml -out dist`
- output directory: `dist`

The export writes `dist/vercel.json` with cache, content type, and security headers.

## Favicon Assets

Export mode resolves favicon assets before writing `index.html`.

- `favicon.source_path` reads a local image from the build machine.
- `favicon.source_url` downloads a remote image during the build.
- if neither is configured, `avatar_url` is used as the favicon source.
- if no image source exists, Mini-Link writes an initials-based `favicon.svg`.

PNG, JPEG, and GIF sources are center-cropped into `favicon.png` and `apple-touch-icon.png`. SVG sources are copied as `favicon.svg`.

Optional upload is configured with `asset_upload`. It supports S3-compatible APIs, including Cloudflare R2:

```yaml
asset_upload:
  provider: r2
  endpoint: https://<account-id>.r2.cloudflarestorage.com
  bucket: mini-link-assets
  access_key_id: <R2_ACCESS_KEY_ID>
  secret_access_key: <R2_SECRET_ACCESS_KEY>
  public_base_url: https://assets.example.com
  prefix: mini-link
```

YAML values are literal; Mini-Link does not expand `${...}` placeholders. For real secrets, prefer `.env` config or CI-provided environment variables instead of committed YAML.

## Caching

Server mode sends:

- strong `ETag` based on rendered HTML
- `Last-Modified`
- `Cache-Control: public, max-age=300, s-maxage=86400, stale-while-revalidate=604800` by default

Static export writes equivalent platform headers. Increase `cache_seconds` when config changes are infrequent.

## SEO and Security Headers

Mini-Link emits canonical, Open Graph, Twitter Card, and JSON-LD metadata on the page. Server mode and static export also provide `robots.txt`, `sitemap.xml`, and `llms.txt` when `base_url` is configured.

The default security headers are intentionally strict because Mini-Link does not need client JavaScript, frames, forms, camera, microphone, geolocation, or payment APIs. If external icon, avatar, or favicon URLs are configured, Mini-Link adds `img-src 'self' https: data:` so those images can load.
