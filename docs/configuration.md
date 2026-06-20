# Configuration

Mini-Link can load the same profile from JSON, a small YAML subset, or environment variables.

## JSON

```bash
minilink serve -config examples/mini-link.json
```

JSON is the most exact format and supports every field directly.

## YAML

```bash
minilink serve -config examples/mini-link.yaml
```

Mini-Link does not depend on a YAML package, so YAML support intentionally covers the simple profile shape used by this project:

- top-level scalar keys
- `links:` as a list of maps
- strings, booleans, and integer `cache_seconds`
- comments with `#` outside quoted strings

Advanced YAML features such as anchors, multi-line strings, nested objects outside `links`, and custom tags are not supported.

## Env

```bash
minilink serve -config examples/mini-link.env.example
```

You can also run with only process environment variables:

```bash
MINI_LINK_NAME="WPitombeira" \
MINI_LINK_BIO="Builder, AI systems, automation" \
MINI_LINK_LINK_1_TITLE="Website" \
MINI_LINK_LINK_1_URL="https://wpitombeira.com.br" \
go run ./cmd/minilink serve
```

Link variables use numbered groups:

```text
MINI_LINK_LINK_1_TITLE
MINI_LINK_LINK_1_URL
MINI_LINK_LINK_1_ICON
MINI_LINK_LINK_1_FEATURED
MINI_LINK_LINK_1_REL
```

For larger env-only deployments, use `MINI_LINK_LINKS_JSON` with a JSON array of links.
