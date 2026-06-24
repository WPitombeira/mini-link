# Agent Runbook

Use this runbook when an agent is changing or deploying Mini-Link.

## Validate

```bash
go test ./...
go run ./cmd/minilink validate -config examples/mini-link.yaml
go run ./cmd/minilink validate -config examples/mini-link.json
go run ./cmd/minilink validate -config examples/mini-link.env.example
go run ./cmd/minilink validate -config examples/dropdowns.yaml
go run ./cmd/minilink validate -config examples/dropdowns.json
go run ./cmd/minilink validate -config examples/custom-icons.yaml
go run ./cmd/minilink validate -config examples/custom-icons.json
go run ./cmd/minilink validate -config examples/minimal.yaml
go run ./cmd/minilink icons -out docs/icons.html
go run ./cmd/minilink export -config examples/mini-link.yaml -out dist
```

## Browser QA

```bash
go run ./cmd/minilink serve -addr :8080 -config examples/mini-link.yaml
```

Then check:

- `/` renders the profile and all links.
- `/favicon.svg` or `/favicon.png` returns the generated favicon.
- `/site.webmanifest` returns valid JSON and references existing icon files.
- a featured link has distinct styling.
- dropdown groups open and close with mouse and keyboard.
- SVG icons render inline.
- custom inline icons render without extra requests.
- external icon, avatar, and favicon URLs render only when configured and CSP includes `https:` in `img-src`.
- link clicks navigate to the expected targets.
- refreshing with `If-None-Match` returns `304`.
- mobile width has no overflow.
- `/robots.txt`, `/sitemap.xml`, and `/llms.txt` render when `base_url` is configured.
- Lighthouse Performance, Accessibility, Best Practices, and SEO score 100 on the served page.

## Release Checklist

- `LICENSE` remains MIT.
- `README.md`, `docs/configuration.md`, `docs/deployment.md`, and `docs/icons.md` match the current CLI.
- `docs/icons.html` is regenerated after icon changes.
- `go test ./...` passes.
- `go vet ./...` passes.
- `go build -trimpath -ldflags="-s -w" -o bin/minilink ./cmd/minilink` succeeds.
