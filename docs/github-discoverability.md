# GitHub Discoverability

Use this checklist to make Mini-Link easier to find from GitHub search, Google, social shares, and profile visitors.

## Repository Settings

Set these in GitHub's repository **About** panel.

Description:

```text
Blazing-fast, self-hosted LittleLink alternative written in Go. No JS by default, no external runtime libs, Cloudflare-ready.
```

Website:

```text
https://wpitombeira.com.br
```

Topics:

```text
go
link-in-bio
littlelink
cloudflare-workers
static-site-generator
self-hosted
no-javascript
seo
svg-icons
mit-license
```

Social preview:

- Upload [github-social-preview.png](github-social-preview.png) in repository **Settings** -> **Social preview**.
- Keep [github-social-preview.svg](github-social-preview.svg) as the editable source.

## README Checklist

The README should keep these elements above the fold:

- project name
- social preview or screenshot
- one-sentence positioning
- live demo link
- quick start
- "Why Mini-Link" bullets
- deploy links
- suggested metadata/topics

Searchers should understand in a few seconds that Mini-Link is:

- a LittleLink alternative
- written in Go
- self-hosted
- fast and cache-friendly
- usable on Cloudflare Workers Static Assets and Pages
- MIT licensed

## Profile README Snippet

Use this in a GitHub profile README or pinned-project section:

```markdown
### Mini-Link

[Mini-Link](https://github.com/WPitombeira/mini-link) is my blazing-fast, self-hosted LittleLink alternative written in Go. It runs without JavaScript by default, exports SEO-ready static assets, and deploys cleanly to Cloudflare Workers Static Assets.

Live: https://wpitombeira.com.br
```

## Launch Copy

Short:

```text
Mini-Link is a blazing-fast LittleLink alternative written in Go: no JS by default, no external runtime libs, SEO-ready static output, and Cloudflare Workers Static Assets deployment.
```

Long:

```text
I built Mini-Link as a small, fast, self-hosted replacement for LittleLink-style link-in-bio pages. It is written in Go, uses only the standard library, supports JSON/YAML/env config, renders inline SVG icons, generates favicon/robots/sitemap/llms.txt output, and deploys as either a Go server or static assets on Cloudflare, Vercel, Docker, or any plain host.
```

## Maintenance Signals

These help humans and search systems trust the repository:

- keep examples runnable
- keep `docs/icons.html` generated after icon changes
- cut releases when config or deployment behavior changes
- add screenshots after visual template changes
- keep the live demo current
- keep deployment docs tested against the latest Wrangler/Cloudflare flow
