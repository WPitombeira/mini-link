package render

import (
	"bytes"
	"encoding/json"
	"html"
	"html/template"
	"net/url"
	"strings"

	"github.com/WPitombeira/mini-link/internal/config"
	"github.com/WPitombeira/mini-link/internal/icons"
)

type pageData struct {
	Config      config.Config
	Links       []linkData
	Year        int
	Canonical   string
	Description string
	Template    string
	Schema      template.JS
	Favicons    []FaviconLink
	GoogleAdsID string
}

type FaviconLink struct {
	Rel   string
	Href  string
	Type  string
	Sizes string
}

type Assets struct {
	Favicons []FaviconLink
}

type linkData struct {
	config.Link
	Icon     template.HTML
	IconURL  string
	Host     string
	New      bool
	Style    template.CSS
	Children []linkData
}

func Page(cfg config.Config) ([]byte, error) {
	return PageWithAssets(cfg, DefaultAssets(cfg))
}

func PageWithAssets(cfg config.Config, assets Assets) ([]byte, error) {
	links := buildLinks(cfg.Links, customIconMap(cfg.CustomIcons))

	data := pageData{
		Config:      cfg,
		Links:       links,
		Year:        2026,
		Canonical:   canonicalURL(cfg),
		Description: description(cfg),
		Template:    "theme-" + cfg.Template,
		Schema:      template.JS(schemaJSON(cfg, links)),
		Favicons:    assets.Favicons,
		GoogleAdsID: cfg.Tracking.GoogleAdsID,
	}

	var buf bytes.Buffer
	err := pageTemplate.Execute(&buf, data)
	return buf.Bytes(), err
}

func DefaultAssets(cfg config.Config) Assets {
	source := cfg.Favicon.SourceURL
	if source == "" {
		source = cfg.AvatarURL
	}
	if source != "" {
		return Assets{Favicons: []FaviconLink{{Rel: "icon", Href: source}}}
	}
	return Assets{Favicons: []FaviconLink{{Rel: "icon", Href: "/favicon.svg", Type: "image/svg+xml"}}}
}

func buildLinks(links []config.Link, customIcons map[string]customIcon) []linkData {
	out := make([]linkData, 0, len(links))
	for _, link := range links {
		icon, iconURL := resolveIcon(link, customIcons)
		out = append(out, linkData{
			Link:     link,
			Icon:     icon,
			IconURL:  iconURL,
			Host:     host(link.URL),
			New:      opensNewTab(link.URL),
			Style:    linkStyle(link),
			Children: buildLinks(link.Links, customIcons),
		})
	}
	return out
}

func linkStyle(link config.Link) template.CSS {
	if link.Color == "" {
		return ""
	}
	return template.CSS("--link-color:" + link.Color)
}

type customIcon struct {
	SVG string
	URL string
}

func customIconMap(items []config.CustomIcon) map[string]customIcon {
	out := make(map[string]customIcon, len(items))
	for _, item := range items {
		out[item.Name] = customIcon{SVG: customIconSVG(item), URL: item.URL}
	}
	return out
}

func customIconSVG(item config.CustomIcon) string {
	if item.SVG != "" {
		return item.SVG
	}
	if item.Path == "" {
		return ""
	}
	return `<svg viewBox="` + html.EscapeString(item.ViewBox) + `" aria-hidden="true"><path d="` + html.EscapeString(item.Path) + `"/></svg>`
}

func resolveIcon(link config.Link, customIcons map[string]customIcon) (template.HTML, string) {
	if link.IconURL != "" {
		return "", link.IconURL
	}
	if item, ok := customIcons[link.Icon]; ok {
		if item.URL != "" {
			return "", item.URL
		}
		return template.HTML(item.SVG), ""
	}
	icon, ok := icons.Get(link.Icon)
	if !ok {
		icon, _ = icons.Get("link")
	}
	return template.HTML(icon.SVG), ""
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

func canonicalURL(cfg config.Config) string {
	if cfg.BaseURL == "" {
		return ""
	}
	parsed, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return ""
	}
	parsed.RawQuery = ""
	parsed.Fragment = ""
	if parsed.Path == "" {
		parsed.Path = "/"
	}
	return parsed.String()
}

func description(cfg config.Config) string {
	if cfg.Bio != "" {
		return cfg.Bio
	}
	return cfg.Title
}

func schemaJSON(cfg config.Config, links []linkData) string {
	sameAs := externalURLs(links)
	payload := map[string]any{
		"@context":    "https://schema.org",
		"@type":       "Person",
		"name":        cfg.Name,
		"description": description(cfg),
		"url":         canonicalURL(cfg),
		"sameAs":      sameAs,
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}
	return string(raw)
}

func externalURLs(links []linkData) []string {
	var urls []string
	seen := map[string]bool{}
	collectExternalURLs(links, seen, &urls)
	return urls
}

func collectExternalURLs(links []linkData, seen map[string]bool, urls *[]string) {
	for _, link := range links {
		if len(link.Children) > 0 {
			collectExternalURLs(link.Children, seen, urls)
			continue
		}
		if strings.HasPrefix(link.URL, "https://") || strings.HasPrefix(link.URL, "http://") {
			if !seen[link.URL] {
				seen[link.URL] = true
				*urls = append(*urls, link.URL)
			}
		}
	}
}

var pageTemplate = template.Must(template.New("page").Parse(`{{define "linkItem"}}{{if .Children}}<details class="dropdown{{if .Featured}} featured{{end}}"{{if .Open}} open{{end}}>
<summary>
<span class="icon">{{if .IconURL}}<img src="{{.IconURL}}" alt="" loading="eager" decoding="async">{{else}}{{.Icon}}{{end}}</span>
<span class="label"><span class="title">{{.Title}}</span><span class="host">{{len .Children}} links</span></span>
<span class="chevron" aria-hidden="true">⌄</span>
</summary>
<div class="dropdown-links">
{{range .Children}}{{template "linkItem" .}}{{end}}</div>
</details>
{{else}}<a class="link{{if .Featured}} featured{{end}}" href="{{.URL}}" rel="{{.Rel}}"{{if .New}} target="_blank"{{end}}{{if .Style}} style="{{.Style}}"{{end}}>
<span class="icon">{{if .IconURL}}<img src="{{.IconURL}}" alt="" loading="eager" decoding="async">{{else}}{{.Icon}}{{end}}</span>
<span class="label"><span class="title">{{.Title}}</span><span class="host">{{.Host}}</span></span>
<span class="arrow" aria-hidden="true">›</span>
</a>
{{end}}{{end}}<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<meta name="color-scheme" content="light">
<meta name="description" content="{{.Description}}">
<meta name="robots" content="index,follow">
{{if .Canonical}}<link rel="canonical" href="{{.Canonical}}">{{end}}
<meta property="og:type" content="profile">
<meta property="og:title" content="{{.Config.Title}}">
<meta property="og:description" content="{{.Description}}">
{{if .Canonical}}<meta property="og:url" content="{{.Canonical}}">{{end}}
<meta name="twitter:card" content="summary">
<meta name="twitter:title" content="{{.Config.Title}}">
<meta name="twitter:description" content="{{.Description}}">
{{range .Favicons}}<link rel="{{.Rel}}" href="{{.Href}}"{{if .Type}} type="{{.Type}}"{{end}}{{if .Sizes}} sizes="{{.Sizes}}"{{end}}>
{{end}}<link rel="manifest" href="/site.webmanifest">
<title>{{.Config.Title}}</title>
<script type="application/ld+json">{{.Schema}}</script>
{{if .GoogleAdsID}}<script async src="https://www.googletagmanager.com/gtag/js?id={{.GoogleAdsID}}"></script>
<script>window.dataLayer=window.dataLayer||[];function gtag(){dataLayer.push(arguments)}gtag("js",new Date);gtag("config","{{.GoogleAdsID}}");</script>
{{end}}
<style>
:root{--accent:{{.Config.Accent}};--bg:#fff;--text:#111827;--muted:#5b6472;--line:#d8dee8;--soft:#f8fafc;--shadow:0 16px 44px rgba(17,24,39,.10);--radius:18px}
*{box-sizing:border-box}
html{font-family:ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;background:var(--bg);color:var(--text);text-rendering:optimizeLegibility}
body{margin:0;min-height:100vh;display:grid;place-items:center;padding:32px 18px;line-height:1.5}
body.theme-glass{--bg:#f5f7fb;--text:#0b1220;--muted:#475569;--line:rgba(148,163,184,.42);background:radial-gradient(circle at 50% -12%,rgba(255,255,255,.94),rgba(245,247,251,.88) 36%,#eef2f7 100%)}
body.theme-terminal{--bg:#07110e;--text:#d8ffe8;--muted:#91c7a6;--line:#1f4f38;--soft:#0b1b15;--shadow:0 18px 52px rgba(0,0,0,.34);background:#07110e;color:var(--text)}
main{width:min(100%,540px);min-height:min(820px,calc(100vh - 64px));display:grid;grid-template-rows:auto 1fr auto;border:1px solid var(--line);border-radius:var(--radius);padding:28px 28px 34px;background:#fff}
.theme-glass main{border-color:rgba(255,255,255,.68);background:linear-gradient(145deg,rgba(255,255,255,.78),rgba(255,255,255,.46));box-shadow:var(--shadow);backdrop-filter:blur(18px) saturate(1.25)}
.theme-terminal main{background:linear-gradient(180deg,#0a1712,#07110e);border-color:#23533d;border-radius:12px}
.topbar{display:flex;align-items:center;justify-content:space-between;font-weight:760;font-size:15px;letter-spacing:0;color:#0b0f17}
.theme-terminal .topbar{font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;color:#8ff0b2}
.brand{display:inline-flex;align-items:center;gap:10px}
.brand svg{width:24px;height:24px;fill:currentColor}
.menu{font-size:24px;line-height:1;color:#0b0f17}
.theme-terminal .menu{color:#8ff0b2}
.profile{display:grid;align-content:center;gap:24px}
.identity{text-align:center;display:grid;justify-items:center;gap:22px}
.avatar{width:96px;height:120px;border:2px solid var(--text);border-radius:999px;display:grid;place-items:center;background:#fff;color:var(--text);font-weight:780;font-size:36px;letter-spacing:0;overflow:hidden}
.avatar img{width:100%;height:100%;border-radius:inherit;object-fit:cover;object-position:center 38%;display:block}
.theme-glass .avatar{border-color:rgba(255,255,255,.82);background:linear-gradient(145deg,rgba(255,255,255,.92),rgba(255,255,255,.58));box-shadow:inset 0 1px 0 rgba(255,255,255,.9),0 20px 44px rgba(15,23,42,.13)}
.theme-terminal .avatar{border-color:#44d17c;background:#07110e;color:#8ff0b2;border-radius:14px;font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace}
h1{font-size:32px;line-height:1.08;margin:0 0 5px;font-weight:780;letter-spacing:0}
.theme-terminal h1{font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;color:#d8ffe8}
.bio{margin:0;color:var(--muted);font-size:15px}
.links{display:grid;gap:10px}
.link,.dropdown summary{--button-bg:color-mix(in srgb,var(--link-color,var(--text)) 14%,#fff);--button-border:color-mix(in srgb,var(--link-color,var(--text)) 34%,var(--line));min-height:56px;display:grid;grid-template-columns:24px 1fr 24px;align-items:center;gap:14px;padding:13px 15px;border:1px solid var(--button-border);border-radius:10px;color:var(--text);text-decoration:none;background:var(--button-bg);transition:transform .16s ease,border-color .16s ease,box-shadow .16s ease}
.dropdown summary{cursor:pointer;list-style:none}
.dropdown summary::-webkit-details-marker{display:none}
.link:hover,.dropdown summary:hover{border-color:var(--link-color,var(--accent));box-shadow:0 10px 24px rgba(17,24,39,.08);transform:translateY(-1px)}
.link:focus-visible,.dropdown summary:focus-visible{outline:3px solid color-mix(in srgb,var(--link-color,var(--accent)) 32%,transparent);outline-offset:3px}
.link.featured{min-height:76px;border-color:var(--link-color,var(--accent));background:color-mix(in srgb,var(--link-color,var(--accent)) 18%,#fff)}
.dropdown.featured summary{min-height:76px;border-color:var(--link-color,var(--accent));background:color-mix(in srgb,var(--link-color,var(--accent)) 18%,#fff)}
.theme-glass .link,.theme-glass .dropdown summary{border-color:color-mix(in srgb,var(--link-color,#fff) 28%,rgba(255,255,255,.74));background:linear-gradient(145deg,color-mix(in srgb,var(--link-color,#fff) 18%,rgba(255,255,255,.82)),color-mix(in srgb,var(--link-color,#fff) 10%,rgba(255,255,255,.58)));box-shadow:inset 0 1px 0 rgba(255,255,255,.9)}
.theme-glass .link.featured{background:linear-gradient(145deg,color-mix(in srgb,var(--link-color,var(--accent)) 24%,rgba(255,255,255,.86)),color-mix(in srgb,var(--link-color,var(--accent)) 14%,rgba(255,255,255,.62)))}
.theme-glass .dropdown.featured summary{background:linear-gradient(145deg,color-mix(in srgb,var(--link-color,var(--accent)) 24%,rgba(255,255,255,.86)),color-mix(in srgb,var(--link-color,var(--accent)) 14%,rgba(255,255,255,.62)))}
.theme-terminal .link,.theme-terminal .dropdown summary{font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;border-color:#24543d;background:#0b1b15;color:#d8ffe8;border-radius:8px}
.theme-terminal .link.featured{background:#0f2a1f;color:#8ff0b2;border-color:#44d17c}
.theme-terminal .dropdown.featured summary{background:#0f2a1f;color:#8ff0b2;border-color:#44d17c}
.link:not(.featured) .label,.dropdown:not(.featured) summary .label{text-align:center}
.link:not(.featured) .host{display:none}
.link:not(.featured) .arrow{visibility:hidden}
.dropdown:not(.featured) summary .host{display:none}
.icon{width:22px;height:22px;color:var(--link-color,currentColor);display:grid;place-items:center}
.icon svg{width:22px;height:22px;display:block;fill:currentColor}
.icon img{width:22px;height:22px;display:block;object-fit:contain}
.label{display:grid;gap:1px;min-width:0}
.title{font-size:15px;font-weight:720;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.host{font-size:12px;color:var(--muted);white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.featured .host{color:#1f6f46}
.arrow{font-size:26px;color:currentColor;line-height:1}
.chevron{font-size:20px;color:currentColor;line-height:1;text-align:center;transition:transform .16s ease}
.dropdown[open] .chevron{transform:rotate(180deg)}
.dropdown-links{display:grid;gap:8px;margin:8px 0 0 18px;padding-left:12px;border-left:1px solid var(--line)}
.dropdown-links .link{min-height:52px}
.theme-terminal .dropdown-links{border-left-color:#24543d}
footer{text-align:center;color:var(--muted);font-size:12px}
footer a{min-height:48px;display:inline-flex;align-items:center;color:var(--text);text-decoration-color:var(--line);text-underline-offset:3px}
.theme-terminal footer a{color:#d8ffe8}
@media (max-width:600px){body{padding:0;place-items:start center}main{min-height:100vh;border:0;border-radius:0;padding:30px 18px}.profile{align-content:start;padding-top:46px}.avatar{width:84px;height:106px;font-size:32px}h1{font-size:30px}.bio{font-size:14px}.link,.dropdown summary{min-height:54px;padding:12px 13px}.link.featured,.dropdown.featured summary{min-height:72px}.dropdown-links{margin-left:10px;padding-left:10px}.theme-glass main{box-shadow:none}.theme-terminal main{border-radius:0}}
@media (prefers-reduced-motion:reduce){.link,.dropdown summary,.chevron{transition:none}.link:hover,.dropdown summary:hover{transform:none}}
</style>
</head>
<body class="{{.Template}}">
<main>
<header class="topbar">
<span class="brand"><span class="icon">{{with index .Links 0}}{{if .IconURL}}<img src="{{.IconURL}}" alt="" loading="eager" decoding="async">{{else}}{{.Icon}}{{end}}{{end}}</span> Mini-Link</span>
<span class="menu" aria-hidden="true">≡</span>
</header>
<section class="profile">
<section class="identity" aria-label="Profile">
<div class="avatar" aria-hidden="true">{{if .Config.AvatarURL}}<img src="{{.Config.AvatarURL}}" alt="" loading="eager" decoding="async">{{else}}{{.Config.Avatar}}{{end}}</div>
<div>
<h1>{{.Config.Name}}</h1>
<p class="bio">{{.Config.Bio}}</p>
</div>
</section>
<nav class="links" aria-label="Links">
{{range .Links}}{{template "linkItem" .}}{{end}}</nav>
</section>
<footer>{{.Config.Footer}} <a href="https://github.com/WPitombeira/mini-link">Mini-Link</a></footer>
</main>
</body>
</html>`))
