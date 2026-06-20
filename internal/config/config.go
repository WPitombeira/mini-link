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
	Name         string `json:"name"`
	Title        string `json:"title"`
	Bio          string `json:"bio"`
	Avatar       string `json:"avatar"`
	BaseURL      string `json:"base_url"`
	Template     string `json:"template"`
	Accent       string `json:"accent"`
	Footer       string `json:"footer"`
	CacheSeconds int    `json:"cache_seconds"`
	Links        []Link `json:"links"`
}

type Link struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Icon     string `json:"icon"`
	Featured bool   `json:"featured"`
	Rel      string `json:"rel"`
}

var colorPattern = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

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
	for i := range cfg.Links {
		cfg.Links[i].Title = strings.TrimSpace(cfg.Links[i].Title)
		cfg.Links[i].URL = strings.TrimSpace(cfg.Links[i].URL)
		cfg.Links[i].Icon = strings.TrimSpace(strings.ToLower(cfg.Links[i].Icon))
		cfg.Links[i].Rel = strings.TrimSpace(cfg.Links[i].Rel)
		if cfg.Links[i].Title == "" {
			return Config{}, fmt.Errorf("links[%d].title is required", i)
		}
		if cfg.Links[i].URL == "" {
			return Config{}, fmt.Errorf("links[%d].url is required", i)
		}
		if err := validateURL(cfg.Links[i].URL); err != nil {
			return Config{}, fmt.Errorf("links[%d].url: %w", i, err)
		}
		if cfg.Links[i].Icon == "" {
			cfg.Links[i].Icon = "link"
		}
		if _, ok := icons.Get(cfg.Links[i].Icon); !ok {
			return Config{}, fmt.Errorf("links[%d].icon %q is not in the SVG catalog", i, cfg.Links[i].Icon)
		}
		if cfg.Links[i].Rel == "" {
			cfg.Links[i].Rel = "me noopener noreferrer"
		} else if strings.HasPrefix(cfg.Links[i].URL, "http://") || strings.HasPrefix(cfg.Links[i].URL, "https://") {
			cfg.Links[i].Rel = ensureRelTokens(cfg.Links[i].Rel, "noopener", "noreferrer")
		}
	}
	if len(cfg.Links) == 0 {
		return Config{}, errors.New("at least one link is required")
	}
	return cfg, nil
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
		links = append(links, Link{
			Title:    title,
			URL:      url,
			Icon:     values[prefix+"ICON"],
			Rel:      values[prefix+"REL"],
			Featured: featured,
		})
	}
	return links
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
