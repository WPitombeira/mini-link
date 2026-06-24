package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/WPitombeira/mini-link/internal/config"
	"github.com/WPitombeira/mini-link/internal/render"
)

const faviconMaxBytes = 5 << 20

type assetFile struct {
	Name        string
	ContentType string
	Body        []byte
}

type faviconResult struct {
	Assets   []assetFile
	Favicons []render.FaviconLink
}

func buildFavicon(ctx context.Context, cfg config.Config) (faviconResult, error) {
	source := cfg.Favicon.SourceURL
	sourcePath := cfg.Favicon.SourcePath
	if source == "" && sourcePath == "" {
		source = cfg.AvatarURL
	}
	if source == "" && sourcePath == "" {
		svg := initialsFaviconSVG(cfg)
		return faviconResult{
			Assets: []assetFile{{Name: "favicon.svg", ContentType: "image/svg+xml", Body: []byte(svg)}},
			Favicons: []render.FaviconLink{
				{Rel: "icon", Href: "/favicon.svg", Type: "image/svg+xml"},
			},
		}, nil
	}

	body, contentType, err := readFaviconSource(ctx, source, sourcePath)
	if err != nil {
		return faviconResult{}, err
	}
	if strings.Contains(contentType, "svg") || looksLikeSVG(body) {
		asset := assetFile{Name: "favicon.svg", ContentType: "image/svg+xml", Body: body}
		return faviconResult{
			Assets: []assetFile{asset},
			Favicons: []render.FaviconLink{
				{Rel: "icon", Href: "/favicon.svg", Type: "image/svg+xml"},
			},
		}, nil
	}

	img, _, err := image.Decode(bytes.NewReader(body))
	if err != nil {
		return faviconResult{}, fmt.Errorf("decode favicon image: %w", err)
	}
	faviconPNG, err := encodePNG(squareResize(img, 32))
	if err != nil {
		return faviconResult{}, err
	}
	applePNG, err := encodePNG(squareResize(img, 180))
	if err != nil {
		return faviconResult{}, err
	}
	return faviconResult{
		Assets: []assetFile{
			{Name: "favicon.png", ContentType: "image/png", Body: faviconPNG},
			{Name: "apple-touch-icon.png", ContentType: "image/png", Body: applePNG},
		},
		Favicons: []render.FaviconLink{
			{Rel: "icon", Href: "/favicon.png", Type: "image/png", Sizes: "32x32"},
			{Rel: "apple-touch-icon", Href: "/apple-touch-icon.png", Sizes: "180x180"},
		},
	}, nil
}

func readFaviconSource(ctx context.Context, sourceURL string, sourcePath string) ([]byte, string, error) {
	if sourcePath != "" {
		body, err := os.ReadFile(filepath.Clean(sourcePath))
		if err != nil {
			return nil, "", err
		}
		return body, contentTypeFromName(sourcePath), nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, "", fmt.Errorf("download favicon source: status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, faviconMaxBytes+1))
	if err != nil {
		return nil, "", err
	}
	if len(body) > faviconMaxBytes {
		return nil, "", fmt.Errorf("favicon source exceeds %d bytes", faviconMaxBytes)
	}
	return body, resp.Header.Get("Content-Type"), nil
}

func contentTypeFromName(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".svg":
		return "image/svg+xml"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	default:
		return ""
	}
}

func looksLikeSVG(body []byte) bool {
	return strings.HasPrefix(strings.TrimSpace(strings.ToLower(string(body[:min(len(body), 256)]))), "<svg")
}

func squareResize(src image.Image, size int) image.Image {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	side := min(w, h)
	x0 := b.Min.X + (w-side)/2
	y0 := b.Min.Y + (h-side)/2
	cropped := image.NewRGBA(image.Rect(0, 0, side, side))
	draw.Draw(cropped, cropped.Bounds(), src, image.Point{X: x0, Y: y0}, draw.Src)

	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			sx := x * side / size
			sy := y * side / size
			dst.Set(x, y, cropped.At(sx, sy))
		}
	}
	return dst
}

func encodePNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func initialsFaviconSVG(cfg config.Config) string {
	label := cfg.Avatar
	if label == "" {
		label = cfg.Name
	}
	if len([]rune(label)) > 2 {
		label = string([]rune(label)[:2])
	}
	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64"><rect width="64" height="64" rx="14" fill="%s"/><text x="32" y="39" text-anchor="middle" font-family="Arial, sans-serif" font-size="24" font-weight="700" fill="#fff">%s</text></svg>`, xmlEscape(cfg.Accent), xmlEscape(label))
}

func writeAssets(out string, assets []assetFile) error {
	for _, asset := range assets {
		target := filepath.Join(out, asset.Name)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(target, asset.Body, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func buildStaticAssets(cfg config.Config) ([]assetFile, error) {
	assets := make([]assetFile, 0, len(cfg.StaticAssets))
	for _, item := range cfg.StaticAssets {
		body, err := os.ReadFile(filepath.Clean(item.SourcePath))
		if err != nil {
			return nil, err
		}
		contentType := item.ContentType
		if contentType == "" {
			contentType = contentTypeFromName(item.OutputPath)
		}
		if contentType == "" {
			contentType = contentTypeFromName(item.SourcePath)
		}
		if contentType == "" {
			contentType = http.DetectContentType(body)
		}
		assets = append(assets, assetFile{Name: item.OutputPath, ContentType: contentType, Body: body})
	}
	return assets, nil
}

func uploadAssets(ctx context.Context, upload config.AssetUpload, assets []assetFile) ([]render.FaviconLink, error) {
	if upload.Provider == "" {
		return nil, nil
	}
	links := make([]render.FaviconLink, 0, len(assets))
	for _, asset := range assets {
		key := path.Join(upload.Prefix, asset.Name)
		if err := putS3Object(ctx, upload, key, asset); err != nil {
			return nil, err
		}
		href := strings.TrimRight(upload.PublicBaseURL, "/") + "/" + key
		switch asset.Name {
		case "apple-touch-icon.png":
			links = append(links, render.FaviconLink{Rel: "apple-touch-icon", Href: href, Sizes: "180x180"})
		case "favicon.png":
			links = append(links, render.FaviconLink{Rel: "icon", Href: href, Type: "image/png", Sizes: "32x32"})
		case "favicon.svg":
			links = append(links, render.FaviconLink{Rel: "icon", Href: href, Type: "image/svg+xml"})
		}
	}
	return links, nil
}

func putS3Object(ctx context.Context, upload config.AssetUpload, key string, asset assetFile) error {
	endpoint, err := url.Parse(upload.Endpoint)
	if err != nil {
		return err
	}
	endpoint.Path = path.Join(endpoint.Path, upload.Bucket, key)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint.String(), bytes.NewReader(asset.Body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", asset.ContentType)
	req.Header.Set("Cache-Control", "public, max-age=31536000, immutable")
	signS3Request(req, upload, asset.Body)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("upload %s: status %d: %s", asset.Name, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

func signS3Request(req *http.Request, upload config.AssetUpload, body []byte) {
	now := time.Now().UTC()
	date := now.Format("20060102")
	amzDate := now.Format("20060102T150405Z")
	payloadHash := sha256Hex(body)
	req.Header.Set("X-Amz-Date", amzDate)
	req.Header.Set("X-Amz-Content-Sha256", payloadHash)
	req.Host = req.URL.Host

	signedHeaders := []string{"cache-control", "content-type", "host", "x-amz-content-sha256", "x-amz-date"}
	canonicalHeaders := "cache-control:" + req.Header.Get("Cache-Control") + "\n" +
		"content-type:" + req.Header.Get("Content-Type") + "\n" +
		"host:" + req.Host + "\n" +
		"x-amz-content-sha256:" + payloadHash + "\n" +
		"x-amz-date:" + amzDate + "\n"
	sort.Strings(signedHeaders)
	signed := strings.Join(signedHeaders, ";")
	canonicalRequest := strings.Join([]string{
		req.Method,
		req.URL.EscapedPath(),
		req.URL.RawQuery,
		canonicalHeaders,
		signed,
		payloadHash,
	}, "\n")
	scope := date + "/" + upload.Region + "/s3/aws4_request"
	stringToSign := "AWS4-HMAC-SHA256\n" + amzDate + "\n" + scope + "\n" + sha256Hex([]byte(canonicalRequest))
	signingKey := s3SigningKey(upload.SecretAccessKey, date, upload.Region)
	signature := hex.EncodeToString(hmacSHA256(signingKey, stringToSign))
	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential="+upload.AccessKeyID+"/"+scope+", SignedHeaders="+signed+", Signature="+signature)
}

func s3SigningKey(secret string, date string, region string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secret), date)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, "s3")
	return hmacSHA256(kService, "aws4_request")
}

func hmacSHA256(key []byte, value string) []byte {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte(value))
	return mac.Sum(nil)
}

func sha256Hex(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

func siteManifest(cfg config.Config, favicons []render.FaviconLink) string {
	var icons []string
	for _, favicon := range favicons {
		if favicon.Rel != "icon" && favicon.Rel != "apple-touch-icon" {
			continue
		}
		sizes := favicon.Sizes
		if sizes == "" && favicon.Type == "image/svg+xml" {
			sizes = "any"
		}
		icon := fmt.Sprintf(`{"src":%q,"sizes":%q`, favicon.Href, sizes)
		if favicon.Type != "" {
			icon += fmt.Sprintf(`,"type":%q`, favicon.Type)
		}
		icon += "}"
		icons = append(icons, icon)
	}
	return fmt.Sprintf(`{"name":%q,"short_name":%q,"icons":[%s],"start_url":"/","display":"standalone","background_color":"#ffffff","theme_color":%q}
`, cfg.Title, cfg.Name, strings.Join(icons, ","), cfg.Accent)
}

func init() {
	image.RegisterFormat("jpeg", "\xff\xd8", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "\x89PNG\r\n\x1a\n", png.Decode, png.DecodeConfig)
	image.RegisterFormat("gif", "GIF8?a", gif.Decode, gif.DecodeConfig)
}
