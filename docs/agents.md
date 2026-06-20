# Agent Runbook

Use this runbook when an agent is changing or deploying Mini-Link.

## Validate

```bash
go test ./...
go run ./cmd/minilink validate -config examples/mini-link.yaml
go run ./cmd/minilink validate -config examples/mini-link.json
go run ./cmd/minilink validate -config examples/mini-link.env.example
go run ./cmd/minilink icons -out docs/icons.html
go run ./cmd/minilink export -config examples/mini-link.yaml -out dist
```

## Browser QA

```bash
go run ./cmd/minilink serve -addr :8080 -config examples/mini-link.yaml
```

Then check:

- `/` renders the profile and all links.
- a featured link has distinct styling.
- SVG icons render inline.
- link clicks navigate to the expected targets.
- refreshing with `If-None-Match` returns `304`.
- mobile width has no overflow.
- `/robots.txt` and `/sitemap.xml` render when `base_url` is configured.
- Lighthouse Performance, Accessibility, Best Practices, and SEO score 100 on the served page.

## Release Checklist

- `LICENSE` remains MIT.
- `README.md`, `docs/configuration.md`, `docs/deployment.md`, and `docs/icons.md` match the current CLI.
- `docs/icons.html` is regenerated after icon changes.
- `go test ./...` passes.
- `go vet ./...` passes.
- `go build -trimpath -ldflags="-s -w" -o bin/minilink ./cmd/minilink` succeeds.
