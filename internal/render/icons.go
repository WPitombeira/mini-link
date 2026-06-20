package render

import (
	"bytes"
	"html/template"

	"github.com/WPitombeira/mini-link/internal/icons"
)

type iconCatalogItem struct {
	Name   string
	Label  string
	Source string
	SVG    template.HTML
}

func IconCatalog(items []icons.Icon) ([]byte, error) {
	view := make([]iconCatalogItem, 0, len(items))
	for _, item := range items {
		view = append(view, iconCatalogItem{
			Name:   item.Name,
			Label:  item.Label,
			Source: item.Source,
			SVG:    template.HTML(item.SVG),
		})
	}
	var buf bytes.Buffer
	err := iconTemplate.Execute(&buf, view)
	return buf.Bytes(), err
}

var iconTemplate = template.Must(template.New("icons").Parse(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Mini-Link SVG Icon Catalog</title>
<style>
html{font-family:ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;color:#111827;background:#fff}
body{margin:0;padding:32px}
main{max-width:980px;margin:auto}
h1{font-size:30px;margin:0 0 8px}
p{color:#5b6472;line-height:1.55}
.grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(180px,1fr));gap:12px;margin-top:24px}
.card{border:1px solid #d8dee8;border-radius:8px;padding:16px;display:grid;gap:10px}
.icon{width:32px;height:32px;color:#0f766e}
.icon svg{width:32px;height:32px;fill:currentColor}
code{background:#f3f5f8;border:1px solid #e5e9f0;border-radius:6px;padding:2px 5px}
small{color:#6b7280;line-height:1.4}
</style>
</head>
<body>
<main>
<h1>Mini-Link SVG Icon Catalog</h1>
<p>Use the icon key in any JSON, YAML, or env config. SVGs are embedded at build time and rendered inline, so no browser request is needed for icons.</p>
<section class="grid">
{{range .}}<article class="card">
<div class="icon">{{.SVG}}</div>
<strong>{{.Label}}</strong>
<code>{{.Name}}</code>
<small>{{.Source}}</small>
</article>
{{end}}</section>
</main>
</body>
</html>
`))
