package render

import (
	"bytes"
	"html/template"
	"net/url"
	"strings"

	"github.com/WPitombeira/mini-link/internal/config"
	"github.com/WPitombeira/mini-link/internal/icons"
)

type pageData struct {
	Config config.Config
	Links  []linkData
	Year   int
}

type linkData struct {
	config.Link
	Icon template.HTML
	Host string
	New  bool
}

func Page(cfg config.Config) ([]byte, error) {
	links := make([]linkData, 0, len(cfg.Links))
	for _, link := range cfg.Links {
		icon, ok := icons.Get(link.Icon)
		if !ok {
			icon, _ = icons.Get("link")
		}
		links = append(links, linkData{Link: link, Icon: template.HTML(icon.SVG), Host: host(link.URL), New: opensNewTab(link.URL)})
	}

	var buf bytes.Buffer
	err := pageTemplate.Execute(&buf, pageData{Config: cfg, Links: links, Year: 2026})
	return buf.Bytes(), err
}

func opensNewTab(raw string) bool {
	return strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://")
}

func host(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Hostname() == "" {
		return raw
	}
	return strings.TrimPrefix(parsed.Hostname(), "www.")
}

var pageTemplate = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<meta name="color-scheme" content="light">
<meta name="description" content="{{.Config.Bio}}">
<title>{{.Config.Title}}</title>
<style>
:root{--accent:{{.Config.Accent}};--bg:#fff;--text:#111827;--muted:#5b6472;--line:#d8dee8;--soft:#f8fafc;--shadow:0 12px 32px rgba(17,24,39,.08)}
*{box-sizing:border-box}
html{font-family:ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;background:var(--bg);color:var(--text);text-rendering:optimizeLegibility}
body{margin:0;min-height:100vh;display:grid;place-items:center;padding:32px 18px;line-height:1.5}
main{width:min(100%,520px);min-height:min(820px,calc(100vh - 64px));display:grid;grid-template-rows:auto 1fr auto;border:1px solid var(--line);border-radius:8px;padding:28px 28px 34px}
.topbar{display:flex;align-items:center;justify-content:space-between;font-weight:760;font-size:15px;letter-spacing:0;color:#0b0f17}
.brand{display:inline-flex;align-items:center;gap:10px}
.brand svg{width:24px;height:24px;fill:currentColor}
.menu{font-size:24px;line-height:1;color:#0b0f17}
.profile{display:grid;align-content:center;gap:24px}
.identity{text-align:center;display:grid;justify-items:center;gap:14px}
.avatar{width:96px;height:96px;border:2px solid var(--text);border-radius:999px;display:grid;place-items:center;background:#fff;color:var(--text);font-weight:780;font-size:36px;letter-spacing:0}
h1{font-size:32px;line-height:1.08;margin:0 0 5px;font-weight:780;letter-spacing:0}
.bio{margin:0;color:var(--muted);font-size:15px}
.links{display:grid;gap:10px}
.link{min-height:54px;display:grid;grid-template-columns:24px 1fr 24px;align-items:center;gap:14px;padding:13px 15px;border:1px solid var(--text);border-radius:5px;color:var(--text);text-decoration:none;background:#fff;transition:transform .16s ease,border-color .16s ease,box-shadow .16s ease}
.link:hover{border-color:var(--accent);box-shadow:0 10px 24px rgba(17,24,39,.08);transform:translateY(-1px)}
.link:focus-visible{outline:3px solid color-mix(in srgb,var(--accent) 32%,transparent);outline-offset:3px}
.link.featured{min-height:76px;border-color:var(--accent);color:#166534;background:#fbfffd}
.link:not(.featured) .label{text-align:center}
.link:not(.featured) .host{display:none}
.link:not(.featured) .arrow{visibility:hidden}
.icon{width:22px;height:22px;color:currentColor;display:grid;place-items:center}
.icon svg{width:22px;height:22px;display:block;fill:currentColor}
.label{display:grid;gap:1px;min-width:0}
.title{font-size:15px;font-weight:720;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.host{font-size:12px;color:var(--muted);white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.featured .host{color:#1f6f46}
.arrow{font-size:26px;color:currentColor;line-height:1}
footer{text-align:center;color:var(--muted);font-size:12px}
footer a{color:var(--text);text-decoration-color:var(--line);text-underline-offset:3px}
@media (max-width:600px){body{padding:0;place-items:start center}main{min-height:100vh;border:0;border-radius:0;padding:30px 18px}.profile{align-content:start;padding-top:52px}.avatar{width:86px;height:86px;font-size:32px}h1{font-size:30px}.bio{font-size:14px}.link{min-height:52px;padding:12px 13px}.link.featured{min-height:72px}}
@media (prefers-reduced-motion:reduce){.link{transition:none}.link:hover{transform:none}}
</style>
</head>
<body>
<main>
<header class="topbar">
<span class="brand"><span class="icon">{{with index .Links 0}}{{.Icon}}{{end}}</span> Mini-Link</span>
<span class="menu" aria-hidden="true">≡</span>
</header>
<section class="profile">
<section class="identity" aria-label="Profile">
<div class="avatar">{{.Config.Avatar}}</div>
<div>
<h1>{{.Config.Name}}</h1>
<p class="bio">{{.Config.Bio}}</p>
</div>
</section>
<nav class="links" aria-label="Links">
{{range .Links}}<a class="link{{if .Featured}} featured{{end}}" href="{{.URL}}" rel="{{.Rel}}"{{if .New}} target="_blank"{{end}}>
<span class="icon">{{.Icon}}</span>
<span class="label"><span class="title">{{.Title}}</span><span class="host">{{.Host}}</span></span>
<span class="arrow" aria-hidden="true">›</span>
</a>
{{end}}</nav>
</section>
<footer>{{.Config.Footer}} <a href="https://github.com/WPitombeira/mini-link">Mini-Link</a></footer>
</main>
</body>
</html>`))
