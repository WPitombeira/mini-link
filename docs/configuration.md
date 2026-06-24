# Configuration

Mini-Link can load the same profile from JSON, a small YAML subset, or environment variables.

Common fields:

- `name`
- `title`
- `bio`
- `avatar`
- `base_url`
- `template`: `classic`, `glass`, or `terminal`
- `accent`
- `footer`
- `cache_seconds`
- `links`

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
- nested `links:` under a link item for dropdown groups
- strings, booleans, and integer `cache_seconds`
- comments with `#` outside quoted strings

Advanced YAML features such as anchors, multi-line strings, nested objects outside link groups, and custom tags are not supported.

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

## Dropdowns

A link can either point directly to a `url`, or it can be a dropdown group with nested `links`. Dropdowns render with native HTML `<details>` and `<summary>`, so they work without JavaScript.

```yaml
links:
  - title: Projects
    icon: briefcase
    open: true
    links:
      - title: Mini-Link
        url: https://github.com/WPitombeira/mini-link
        icon: github
      - title: Portfolio
        url: https://example.com
        icon: globe
```

Dropdown rules:

- dropdown groups use `title`, optional `icon`, optional `featured`, optional `open`, and nested `links`
- direct links use `title`, `url`, optional `icon`, optional `featured`, and optional `rel`
- a single item cannot define both `url` and nested `links`
- nesting is limited to three levels to keep the page readable

Env files can use `MINI_LINK_LINKS_JSON` for dropdowns because numbered env variables are intentionally flat.
