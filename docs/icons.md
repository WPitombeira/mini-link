# SVG Icons

Mini-Link ships precompiled inline SVG icons. They add no browser requests and require no package installation.

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

Generic icons are sourced from Heroicons under the MIT License where noted. Brand icons are sourced from Simple Icons under CC0 1.0. Brand names and logos may still be trademarks of their owners. Custom utility icons in this repository are MIT licensed. See [../THIRD_PARTY_NOTICES.md](../THIRD_PARTY_NOTICES.md).
