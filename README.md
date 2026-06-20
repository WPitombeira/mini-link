# Mini-Link

Mini-Link is a LittleLink-style link-in-bio page written in Go. It is designed to be small, fast, cache-friendly, and easy to host anywhere.

It uses only the Go standard library. The rendered page has inline CSS and inline SVG icons, so the browser can load it with a single HTML request.

## Features

- Go HTTP server and static exporter
- JSON, small YAML subset, and env configuration
- strong ETag, Last-Modified, and CDN-friendly Cache-Control headers
- precompiled inline SVG icon catalog
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

## SVG Icons

Generate the browser-viewable icon guide:

```bash
go run ./cmd/minilink icons -out docs/icons.html
```

Then open `docs/icons.html`. See [docs/icons.md](docs/icons.md).

## Deployment

See [docs/deployment.md](docs/deployment.md) for Docker, Docker Compose, manual Go builds, Cloudflare Pages, Cloudflare Workers/static assets, and Vercel.

## License

MIT. See [LICENSE](LICENSE).
