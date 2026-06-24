# Cloudflare Deployment

Mini-Link is a good fit for Cloudflare because the exported site is static: `index.html`, `robots.txt`, `sitemap.xml`, `llms.txt`, and platform header files. Use Pages when you want the simplest Git-backed static deployment. Use Workers Static Assets when you want a Worker project and Wrangler-driven deploys.

## Prepare the Site

Set `base_url` to the production URL before exporting. This keeps canonical, sitemap, and social metadata correct.

```bash
go test ./...
go run ./cmd/minilink validate -config examples/mini-link.yaml
go run ./cmd/minilink export -config examples/mini-link.yaml -out dist
```

Check the exported files:

```bash
ls dist
```

Expected files:

- `index.html`
- `robots.txt`
- `sitemap.xml`
- `llms.txt`
- `_headers`
- `vercel.json`

## Option 1: Cloudflare Pages

Use this path for a static link-in-bio site backed by Git.

Cloudflare Pages settings:

- Framework preset: None
- Build command: `go run ./cmd/minilink export -config examples/mini-link.yaml -out dist`
- Build output directory: `dist`

The exported `_headers` file configures cache and security headers for Pages.

Manual deploy with Wrangler:

```bash
go run ./cmd/minilink export -config examples/mini-link.yaml -out dist
npx wrangler pages deploy dist --project-name mini-link
```

## Option 2: Workers Static Assets

Use this path when you prefer a Worker project or want to add Worker logic later.

Copy the example:

```bash
cp examples/cloudflare/wrangler.toml wrangler.toml
```

Deploy:

```bash
go run ./cmd/minilink export -config examples/mini-link.yaml -out dist
npx wrangler deploy
```

The important Wrangler setting is:

```toml
[assets]
directory = "./dist"
```

Workers Static Assets should be used instead of legacy Workers Sites for new Worker-hosted static projects.

## Custom Domain

After deployment, attach the production domain in Cloudflare and update your Mini-Link config:

```yaml
base_url: https://links.example.com
```

Then export and deploy again so `sitemap.xml`, canonical URL, Open Graph URL, and `robots.txt` all point to the final domain.

## Verify After Deploy

```bash
curl -I https://links.example.com
curl https://links.example.com/robots.txt
curl https://links.example.com/sitemap.xml
curl https://links.example.com/llms.txt
```

Confirm:

- `Cache-Control` includes `s-maxage` and `stale-while-revalidate`.
- `Content-Security-Policy` is present.
- `robots.txt` references the production sitemap URL.
- `sitemap.xml` contains the production `base_url`.
- `llms.txt` lists the production URL and public HTTP links.

## Cache Behavior

Mini-Link exports CDN-friendly headers:

```text
Cache-Control: public, max-age=300, s-maxage=86400, stale-while-revalidate=604800
```

Use a higher `cache_seconds` value if your links rarely change. Keep it lower while iterating on design or profile content.
