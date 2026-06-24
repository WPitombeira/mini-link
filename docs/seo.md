# SEO and AI Readiness

Mini-Link is a single-page link profile, so SEO output focuses on crawlability, canonical metadata, structured data, fast rendering, and AI-readable discovery files.

## Generated Files

Server mode and static export provide:

- `/`: rendered HTML page
- `/robots.txt`: allows crawling and points to the sitemap when `base_url` is set
- `/sitemap.xml`: one canonical URL for the profile
- `/llms.txt`: Markdown summary for AI systems and agents

Static export writes the same files into `dist/`.

## HTML Metadata

The page includes:

- `<title>`
- meta description
- canonical link when `base_url` is set
- robots meta tag
- Open Graph title, description, URL, and profile type
- Twitter summary metadata
- JSON-LD `Person` schema with public HTTP links in `sameAs`

## Performance Notes

Built-in icons and inline custom icons render inside the HTML and require no extra browser request. External icon URLs use `<img>` and can reduce performance because each icon depends on another network request, cache policy, and remote host latency.

Mini-Link keeps JavaScript disabled by default. Dropdowns use native `<details>` and `<summary>`.

## Validation Workflow

```bash
go test ./...
go vet ./...
go build -trimpath -ldflags="-s -w" -o bin/minilink ./cmd/minilink
./bin/minilink validate -config examples/mini-link.yaml
./bin/minilink export -config examples/mini-link.yaml -out dist
./bin/minilink serve -addr :8080 -config examples/mini-link.yaml
```

Then run Lighthouse against `http://127.0.0.1:8080/` and verify:

- Performance: 100
- Accessibility: 100
- Best Practices: 100
- SEO: 100

`llms.txt` is an emerging AI/GEO convention, not a guaranteed Google ranking factor. It is included because it is low-cost, public-data-only, and useful for AI crawler orientation.
