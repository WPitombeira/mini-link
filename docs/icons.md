# SVG Icons

Mini-Link ships precompiled inline SVG icons. They add no browser requests and require no package installation.

Generate the browser-viewable catalog:

```bash
go run ./cmd/minilink icons -out docs/icons.html
open docs/icons.html
```

Current icon keys:

- `code`
- `contact`
- `github`
- `globe`
- `linkedin`
- `link`
- `x`

Use an icon key in JSON, YAML, or env config:

```yaml
links:
  - title: GitHub
    url: https://github.com/WPitombeira
    icon: github
```

Brand icons are sourced from Simple Icons under CC0 1.0. Brand names and logos may still be trademarks of their owners. Custom utility icons in this repository are MIT licensed.
