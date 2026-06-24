package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/WPitombeira/mini-link/internal/icons"
)

type Config struct {
	Name         string       `json:"name"`
	Title        string       `json:"title"`
	Bio          string       `json:"bio"`
	Avatar       string       `json:"avatar"`
	BaseURL      string       `json:"base_url"`
	Template     string       `json:"template"`
	Accent       string       `json:"accent"`
	Footer       string       `json:"footer"`
	CacheSeconds int          `json:"cache_seconds"`
	CustomIcons  []CustomIcon `json:"custom_icons,omitempty"`
	Links        []Link       `json:"links"`
}

type Link struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Icon     string `json:"icon"`
	IconURL  string `json:"icon_url"`
	Featured bool   `json:"featured"`
	Open     bool   `json:"open"`
	Rel      string `json:"rel"`
	Links    []Link `json:"links,omitempty"`
}

type CustomIcon struct {
	Name    string `json:"name"`
	Label   string `json:"label"`
	Source  string `json:"source"`
	SVG     string `json:"svg"`
	ViewBox string `json:"view_box"`
	Path    string `json:"path"`
	URL     string `json:"url"`
}

var colorPattern = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
var iconNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,63}$`)

func Default() Config {
	return Config{
		Name:         "Mini-Link",
		Title:        "Fast links, one tiny Go binary",
		Bio:          "A blazing-fast LittleLink alternative.",
		Avatar:       "ML",
		Template:     "classic",
		Accent:       "#0f766e",
		Footer:       "Mini-Link is MIT licensed.",
		CacheSeconds: 300,
		Links: []Link{
			{Title: "GitHub", URL: "https://github.com/WPitombeira/mini-link", Icon: "github", Featured: true},
		},
	}
}

func Load(path string) (Config, error) {
	cfg := Default()
	if strings.TrimSpace(path) == "" {
		return loadEnv(cfg, os.Environ())
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	lowerPath := strings.ToLower(path)
	switch {
	case filepath.Ext(lowerPath) == ".json":
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return Config{}, fmt.Errorf("parse json config: %w", err)
		}
	case filepath.Ext(lowerPath) == ".yaml" || filepath.Ext(lowerPath) == ".yml":
		parsed, err := parseYAML(string(raw))
		if err != nil {
			return Config{}, fmt.Errorf("parse yaml config: %w", err)
		}
		cfg = merge(cfg, parsed)
	case filepath.Ext(lowerPath) == ".env" || strings.HasSuffix(lowerPath, ".env.example"):
		env, err := parseDotEnv(string(raw))
		if err != nil {
			return Config{}, fmt.Errorf("parse env config: %w", err)
		}
		cfg, err = loadEnv(cfg, env)
		if err != nil {
			return Config{}, err
		}
	default:
		return Config{}, fmt.Errorf("unsupported config extension %q", filepath.Ext(path))
	}
	return normalize(cfg)
}

func normalize(cfg Config) (Config, error) {
	cfg.Name = strings.TrimSpace(cfg.Name)
	cfg.Title = strings.TrimSpace(cfg.Title)
	cfg.Bio = strings.TrimSpace(cfg.Bio)
	cfg.Avatar = strings.TrimSpace(cfg.Avatar)
	cfg.BaseURL = strings.TrimSpace(cfg.BaseURL)
	cfg.Template = strings.TrimSpace(strings.ToLower(cfg.Template))
	cfg.Accent = strings.TrimSpace(cfg.Accent)
	cfg.Footer = strings.TrimSpace(cfg.Footer)
	if cfg.Name == "" {
		return Config{}, errors.New("name is required")
	}
	if cfg.Title == "" {
		cfg.Title = cfg.Name
	}
	if cfg.Avatar == "" {
		cfg.Avatar = initials(cfg.Name)
	}
	if cfg.Template == "" {
		cfg.Template = "classic"
	}
	switch cfg.Template {
	case "classic", "glass", "terminal":
	default:
		return Config{}, fmt.Errorf("template must be one of classic, glass, or terminal")
	}
	if cfg.BaseURL != "" {
		if err := validateAbsoluteHTTPURL(cfg.BaseURL); err != nil {
			return Config{}, fmt.Errorf("base_url: %w", err)
		}
	}
	if cfg.Accent == "" {
		cfg.Accent = "#0f766e"
	}
	if !colorPattern.MatchString(cfg.Accent) {
		return Config{}, fmt.Errorf("accent must be a hex color like #0f766e")
	}
	if contrastRatio(cfg.Accent, "#ffffff") < 3 {
		return Config{}, fmt.Errorf("accent must have at least 3:1 contrast against white")
	}
	if cfg.CacheSeconds <= 0 {
		cfg.CacheSeconds = 300
	}
	customIcons, err := normalizeCustomIcons(cfg.CustomIcons)
	if err != nil {
		return Config{}, err
	}
	cfg.CustomIcons = customIcons
	customIconNames := map[string]bool{}
	for _, icon := range cfg.CustomIcons {
		customIconNames[icon.Name] = true
	}
	for i := range cfg.Links {
		if err := normalizeLink(&cfg.Links[i], fmt.Sprintf("links[%d]", i), 1, customIconNames); err != nil {
			return Config{}, err
		}
	}
	if len(cfg.Links) == 0 {
		return Config{}, errors.New("at least one link is required")
	}
	return cfg, nil
}

func normalizeCustomIcons(input []CustomIcon) ([]CustomIcon, error) {
	if len(input) == 0 {
		return nil, nil
	}
	seen := map[string]bool{}
	out := make([]CustomIcon, 0, len(input))
	for i, icon := range input {
		where := fmt.Sprintf("custom_icons[%d]", i)
		icon.Name = strings.TrimSpace(strings.ToLower(icon.Name))
		icon.Label = strings.TrimSpace(icon.Label)
		icon.Source = strings.TrimSpace(icon.Source)
		icon.SVG = strings.TrimSpace(icon.SVG)
		icon.ViewBox = strings.TrimSpace(icon.ViewBox)
		icon.Path = strings.TrimSpace(icon.Path)
		icon.URL = strings.TrimSpace(icon.URL)
		if icon.Name == "" {
			return nil, fmt.Errorf("%s.name is required", where)
		}
		if !iconNamePattern.MatchString(icon.Name) {
			return nil, fmt.Errorf("%s.name must use lowercase letters, numbers, dashes, or underscores", where)
		}
		if _, ok := icons.Get(icon.Name); ok {
			return nil, fmt.Errorf("%s.name %q conflicts with built-in icon catalog", where, icon.Name)
		}
		if seen[icon.Name] {
			return nil, fmt.Errorf("%s.name %q is duplicated", where, icon.Name)
		}
		seen[icon.Name] = true
		if icon.Label == "" {
			icon.Label = icon.Name
		}
		forms := 0
		if icon.SVG != "" {
			forms++
			if err := validateInlineSVG(icon.SVG); err != nil {
				return nil, fmt.Errorf("%s.svg: %w", where, err)
			}
		}
		if icon.Path != "" || icon.ViewBox != "" {
			forms++
			if icon.Path == "" || icon.ViewBox == "" {
				return nil, fmt.Errorf("%s path icons require both path and view_box", where)
			}
			if err := validateSVGPart(icon.ViewBox, "view_box"); err != nil {
				return nil, fmt.Errorf("%s.%w", where, err)
			}
			if err := validateSVGPart(icon.Path, "path"); err != nil {
				return nil, fmt.Errorf("%s.%w", where, err)
			}
		}
		if icon.URL != "" {
			forms++
			if err := validateAbsoluteHTTPURL(icon.URL); err != nil {
				return nil, fmt.Errorf("%s.url: %w", where, err)
			}
		}
		if forms != 1 {
			return nil, fmt.Errorf("%s must define exactly one of svg, path plus view_box, or url", where)
		}
		out = append(out, icon)
	}
	return out, nil
}

func normalizeLink(link *Link, path string, depth int, customIconNames map[string]bool) error {
	if depth > 3 {
		return fmt.Errorf("%s exceeds maximum dropdown depth of 3", path)
	}
	link.Title = strings.TrimSpace(link.Title)
	link.URL = strings.TrimSpace(link.URL)
	link.Icon = strings.TrimSpace(link.Icon)
	link.IconURL = strings.TrimSpace(link.IconURL)
	if link.IconURL == "" && isHTTPURL(link.Icon) {
		link.IconURL = link.Icon
		link.Icon = ""
	}
	link.Icon = strings.ToLower(link.Icon)
	link.Rel = strings.TrimSpace(link.Rel)
	if link.Title == "" {
		return fmt.Errorf("%s.title is required", path)
	}
	if len(link.Links) > 0 {
		if link.URL != "" {
			return fmt.Errorf("%s cannot define both url and nested links", path)
		}
		if link.Icon == "" {
			link.Icon = "link"
		}
		if link.IconURL != "" {
			if err := validateAbsoluteHTTPURL(link.IconURL); err != nil {
				return fmt.Errorf("%s.icon_url: %w", path, err)
			}
		} else if !iconExists(link.Icon, customIconNames) {
			return fmt.Errorf("%s.icon %q is not in the SVG catalog or custom_icons", path, link.Icon)
		}
		for i := range link.Links {
			if err := normalizeLink(&link.Links[i], fmt.Sprintf("%s.links[%d]", path, i), depth+1, customIconNames); err != nil {
				return err
			}
		}
		return nil
	}
	if link.URL == "" {
		return fmt.Errorf("%s.url is required", path)
	}
	if err := validateURL(link.URL); err != nil {
		return fmt.Errorf("%s.url: %w", path, err)
	}
	if link.Icon == "" {
		link.Icon = "link"
	}
	if link.IconURL != "" {
		if err := validateAbsoluteHTTPURL(link.IconURL); err != nil {
			return fmt.Errorf("%s.icon_url: %w", path, err)
		}
	} else if !iconExists(link.Icon, customIconNames) {
		return fmt.Errorf("%s.icon %q is not in the SVG catalog or custom_icons", path, link.Icon)
	}
	if link.Rel == "" {
		link.Rel = "me noopener noreferrer"
	} else if strings.HasPrefix(link.URL, "http://") || strings.HasPrefix(link.URL, "https://") {
		link.Rel = ensureRelTokens(link.Rel, "noopener", "noreferrer")
	}
	return nil
}

func iconExists(name string, customIconNames map[string]bool) bool {
	if _, ok := icons.Get(name); ok {
		return true
	}
	return customIconNames[name]
}

func validateInlineSVG(svg string) error {
	lower := strings.ToLower(svg)
	if !strings.HasPrefix(lower, "<svg ") && !strings.HasPrefix(lower, "<svg>") {
		return errors.New("must start with <svg")
	}
	if !strings.Contains(lower, "</svg>") {
		return errors.New("must include closing </svg>")
	}
	blocked := []string{"<script", "<foreignobject", "<iframe", "<object", "<embed", "<image", "<style", " onload=", " onclick=", " onerror=", " href=", " xlink:href=", "javascript:"}
	for _, token := range blocked {
		if strings.Contains(lower, token) {
			return fmt.Errorf("blocked unsafe SVG token %q", token)
		}
	}
	return nil
}

func validateSVGPart(value string, name string) error {
	lower := strings.ToLower(value)
	blocked := []string{"<", ">", `"`, "'", "javascript:", "onload", "onclick", "onerror"}
	for _, token := range blocked {
		if strings.Contains(lower, token) {
			return fmt.Errorf("%s contains unsafe token %q", name, token)
		}
	}
	return nil
}

func ensureRelTokens(rel string, tokens ...string) string {
	seen := map[string]bool{}
	parts := strings.Fields(rel)
	for _, part := range parts {
		seen[part] = true
	}
	for _, token := range tokens {
		if !seen[token] {
			parts = append(parts, token)
		}
	}
	return strings.Join(parts, " ")
}

func contrastRatio(a string, b string) float64 {
	la := relativeLuminance(a)
	lb := relativeLuminance(b)
	light := math.Max(la, lb)
	dark := math.Min(la, lb)
	return (light + 0.05) / (dark + 0.05)
}

func relativeLuminance(hex string) float64 {
	r := linearRGB(hexByte(hex[1:3]))
	g := linearRGB(hexByte(hex[3:5]))
	b := linearRGB(hexByte(hex[5:7]))
	return 0.2126*r + 0.7152*g + 0.0722*b
}

func hexByte(value string) float64 {
	n, _ := strconv.ParseUint(value, 16, 8)
	return float64(n) / 255
}

func linearRGB(value float64) float64 {
	if value <= 0.03928 {
		return value / 12.92
	}
	return math.Pow((value+0.055)/1.055, 2.4)
}

func validateURL(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return err
	}
	switch parsed.Scheme {
	case "http", "https":
		if parsed.Host == "" {
			return errors.New("http links require a host")
		}
	case "mailto":
		if parsed.Opaque == "" && parsed.Path == "" {
			return errors.New("mailto links require an address")
		}
	case "tel":
		if parsed.Opaque == "" && parsed.Path == "" {
			return errors.New("tel links require a number")
		}
	default:
		return fmt.Errorf("unsupported scheme %q", parsed.Scheme)
	}
	return nil
}

func validateAbsoluteHTTPURL(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return err
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("must use http or https")
	}
	if parsed.Host == "" {
		return errors.New("requires a host")
	}
	return nil
}

func isHTTPURL(raw string) bool {
	return strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://")
}

func merge(base Config, next Config) Config {
	if next.Name != "" {
		base.Name = next.Name
	}
	if next.Title != "" {
		base.Title = next.Title
	}
	if next.Bio != "" {
		base.Bio = next.Bio
	}
	if next.Avatar != "" {
		base.Avatar = next.Avatar
	}
	if next.BaseURL != "" {
		base.BaseURL = next.BaseURL
	}
	if next.Template != "" {
		base.Template = next.Template
	}
	if next.Accent != "" {
		base.Accent = next.Accent
	}
	if next.Footer != "" {
		base.Footer = next.Footer
	}
	if next.CacheSeconds != 0 {
		base.CacheSeconds = next.CacheSeconds
	}
	if next.CustomIcons != nil {
		base.CustomIcons = next.CustomIcons
	}
	if next.Links != nil {
		base.Links = next.Links
	}
	return base
}

func loadEnv(base Config, entries []string) (Config, error) {
	values := map[string]string{}
	for _, entry := range entries {
		key, value, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}
		values[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}

	cfg := base
	cfg.Name = first(values, "MINI_LINK_NAME", cfg.Name)
	cfg.Title = first(values, "MINI_LINK_TITLE", cfg.Title)
	cfg.Bio = first(values, "MINI_LINK_BIO", cfg.Bio)
	cfg.Avatar = first(values, "MINI_LINK_AVATAR", cfg.Avatar)
	cfg.BaseURL = first(values, "MINI_LINK_BASE_URL", cfg.BaseURL)
	cfg.Template = first(values, "MINI_LINK_TEMPLATE", cfg.Template)
	cfg.Accent = first(values, "MINI_LINK_ACCENT", cfg.Accent)
	cfg.Footer = first(values, "MINI_LINK_FOOTER", cfg.Footer)
	if raw := values["MINI_LINK_CACHE_SECONDS"]; raw != "" {
		ttl, err := strconv.Atoi(raw)
		if err != nil {
			return Config{}, fmt.Errorf("MINI_LINK_CACHE_SECONDS: %w", err)
		}
		cfg.CacheSeconds = ttl
	}
	if raw := values["MINI_LINK_CUSTOM_ICONS_JSON"]; raw != "" {
		var customIcons []CustomIcon
		if err := json.Unmarshal([]byte(raw), &customIcons); err != nil {
			return Config{}, fmt.Errorf("MINI_LINK_CUSTOM_ICONS_JSON: %w", err)
		}
		cfg.CustomIcons = customIcons
	}
	if raw := values["MINI_LINK_LINKS_JSON"]; raw != "" {
		var links []Link
		if err := json.Unmarshal([]byte(raw), &links); err != nil {
			return Config{}, fmt.Errorf("MINI_LINK_LINKS_JSON: %w", err)
		}
		cfg.Links = links
	} else if links := numberedLinks(values); len(links) > 0 {
		cfg.Links = links
	}
	return normalize(cfg)
}

func numberedLinks(values map[string]string) []Link {
	var links []Link
	for i := 1; ; i++ {
		prefix := fmt.Sprintf("MINI_LINK_LINK_%d_", i)
		title := values[prefix+"TITLE"]
		url := values[prefix+"URL"]
		if title == "" && url == "" {
			break
		}
		featured, _ := strconv.ParseBool(values[prefix+"FEATURED"])
		open, _ := strconv.ParseBool(values[prefix+"OPEN"])
		links = append(links, Link{
			Title:    title,
			URL:      url,
			Icon:     first(values, prefix+"ICON", ""),
			IconURL:  first(values, prefix+"ICON_URL", ""),
			Rel:      values[prefix+"REL"],
			Featured: featured,
			Open:     open,
		})
	}
	return links
}

func CountLinks(links []Link) int {
	count := 0
	for _, link := range links {
		if len(link.Links) > 0 {
			count += CountLinks(link.Links)
			continue
		}
		count++
	}
	return count
}

func HasExternalIcons(cfg Config) bool {
	for _, icon := range cfg.CustomIcons {
		if icon.URL != "" {
			return true
		}
	}
	return linksHaveExternalIcons(cfg.Links)
}

func linksHaveExternalIcons(links []Link) bool {
	for _, link := range links {
		if link.IconURL != "" || isHTTPURL(link.Icon) {
			return true
		}
		if linksHaveExternalIcons(link.Links) {
			return true
		}
	}
	return false
}

func parseDotEnv(raw string) ([]string, error) {
	var out []string
	scanner := bufio.NewScanner(strings.NewReader(raw))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return nil, fmt.Errorf("invalid env line %q", line)
		}
		out = append(out, strings.TrimSpace(key)+"="+unquote(strings.TrimSpace(value)))
	}
	return out, scanner.Err()
}

func first(values map[string]string, key string, fallback string) string {
	if value := strings.TrimSpace(values[key]); value != "" {
		return value
	}
	return fallback
}

func initials(name string) string {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "ML"
	}
	var b strings.Builder
	for _, part := range parts {
		b.WriteString(strings.ToUpper(part[:1]))
		if b.Len() >= 2 {
			break
		}
	}
	return b.String()
}
