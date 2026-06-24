# Configuration

Mini-Link can load the same profile from JSON, a small YAML subset, or environment variables.

Common fields:

- `name`
- `title`
- `bio`
- `avatar`
- `avatar_url`
- `base_url`
- `template`: `classic`, `glass`, or `terminal`
- `accent`
- `footer`
- `cache_seconds`
- `favicon`
- `asset_upload`
- `custom_icons`
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

Custom icons in env files use JSON:

```text
MINI_LINK_CUSTOM_ICONS_JSON=[{"name":"sparkle","view_box":"0 0 24 24","path":"M12 2 15 9 22 12 15 15 12 22 9 15 2 12 9 9Z"}]
```

Favicon and upload settings can also be set with env variables:

```text
MINI_LINK_AVATAR_URL=https://cdn.example.com/avatar.png
MINI_LINK_FAVICON_SOURCE_URL=https://cdn.example.com/avatar.png
MINI_LINK_ASSET_UPLOAD_PROVIDER=r2
MINI_LINK_ASSET_UPLOAD_ENDPOINT=https://account-id.r2.cloudflarestorage.com
MINI_LINK_ASSET_UPLOAD_BUCKET=mini-link
MINI_LINK_ASSET_UPLOAD_ACCESS_KEY_ID=...
MINI_LINK_ASSET_UPLOAD_SECRET_ACCESS_KEY=...
MINI_LINK_ASSET_UPLOAD_PUBLIC_BASE_URL=https://assets.example.com
MINI_LINK_ASSET_UPLOAD_PREFIX=profiles/wp
```

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

## Custom Icons

Custom icons can be defined once and reused by `icon` key:

```yaml
custom_icons:
  - name: sparkle
    label: Sparkle
    view_box: "0 0 24 24"
    path: "M12 2 15 9 22 12 15 15 12 22 9 15 2 12 9 9Z"
links:
  - title: Custom
    url: https://example.com
    icon: sparkle
```

Supported custom icon forms:

- `view_box` plus `path`: preferred, safest, rendered as inline SVG
- `svg`: full inline SVG, accepted only when it passes safety checks
- `url`: external SVG/image URL, rendered as an `<img>`

Links can also define a one-off external icon:

```yaml
links:
  - title: CDN icon
    url: https://example.com
    icon_url: https://cdn.example.com/icon.svg
```

External icon URLs add browser requests and require a looser `img-src` Content Security Policy. They are supported for flexibility, but built-in icons and inline custom icons preserve Mini-Link's fastest one-request behavior.

## Favicon

Mini-Link always emits favicon tags and a web manifest.

Default behavior:

- if `favicon.source_url` or `favicon.source_path` is configured during export, Mini-Link reads that image and generates optimized favicon assets
- if no favicon source is configured but `avatar_url` exists, export uses `avatar_url` as the favicon source
- if neither is configured, export creates a small `favicon.svg` from the profile initials

```yaml
avatar_url: https://cdn.example.com/avatar.png
favicon:
  source_url: https://cdn.example.com/avatar.png
```

For local build assets:

```yaml
favicon:
  source_path: ./assets/avatar.png
```

Supported image input formats are PNG, JPEG, GIF, and SVG. Raster images are center-cropped and exported as:

- `favicon.png` at 32x32
- `apple-touch-icon.png` at 180x180

SVG sources are exported as `favicon.svg`.

## Asset Upload

Mini-Link can optionally upload generated favicon assets to an S3-compatible bucket during `export`. Cloudflare R2 uses the same S3-compatible flow.

```yaml
asset_upload:
  provider: r2
  endpoint: https://account-id.r2.cloudflarestorage.com
  bucket: mini-link
  access_key_id: ${MINI_LINK_ASSET_UPLOAD_ACCESS_KEY_ID}
  secret_access_key: ${MINI_LINK_ASSET_UPLOAD_SECRET_ACCESS_KEY}
  public_base_url: https://assets.example.com
  prefix: profiles/wp
```

For Cloudflare Workers Static Assets and Pages, the default recommendation is not to upload favicons to R2. Let Mini-Link write them into `dist/` and let Cloudflare serve/cache them as regular static assets. Use R2/S3 when you need assets shared across multiple deployments, runtime user uploads, or a long-lived bucket independent from the Mini-Link build artifact.
