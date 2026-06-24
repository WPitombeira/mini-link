# SVG Icons

Mini-Link ships precompiled inline SVG icons. They add no browser requests and require no package installation. Users can also define their own icons.

Generate the browser-viewable catalog:

```bash
go run ./cmd/minilink icons -out docs/icons.html
open docs/icons.html
```

Current icon keys:

- `briefcase`
- `calendar`
- `check`
- `code`
- `contact`
- `copy`
- `external-link`
- `file-text`
- `github`
- `globe`
- `home`
- `instagram`
- `linkedin`
- `link`
- `mail`
- `map-pin`
- `moon`
- `phone`
- `rss`
- `sun`
- `terminal`
- `user`
- `whatsapp`
- `x`

Use an icon key in JSON, YAML, or env config:

```yaml
links:
  - title: GitHub
    url: https://github.com/WPitombeira
    icon: github
```

## Custom Icons

Use `custom_icons` when you want to ship your own SVG with the generated page.

Path-based icons are the safest format:

```yaml
custom_icons:
  - name: sparkle
    label: Sparkle
    source: User provided, MIT compatible
    view_box: "0 0 24 24"
    path: "M12 2 15 9 22 12 15 15 12 22 9 15 2 12 9 9Z"
links:
  - title: Inline custom icon
    url: https://example.com
    icon: sparkle
```

Full inline SVG is also supported when it passes Mini-Link's safety checks:

```yaml
custom_icons:
  - name: badge
    svg: '<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 2 22 12 12 22 2 12Z"/></svg>'
```

External SVG URLs are supported, but they cost an extra browser request per icon and require Mini-Link to allow remote images in the Content Security Policy. That can reduce the performance advantage of the one-request default page.

```yaml
custom_icons:
  - name: remote_logo
    url: https://cdn.example.com/logo.svg
links:
  - title: Remote icon by key
    url: https://example.com
    icon: remote_logo
  - title: Remote icon per link
    url: https://example.com/direct
    icon_url: https://cdn.example.com/direct-icon.svg
```

Generic icons are sourced from Heroicons under the MIT License where noted. Brand icons are sourced from Simple Icons under CC0 1.0. Brand names and logos may still be trademarks of their owners. Custom utility icons in this repository are MIT licensed. See [../THIRD_PARTY_NOTICES.md](../THIRD_PARTY_NOTICES.md).
