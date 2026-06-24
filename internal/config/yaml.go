package config

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func parseYAML(raw string) (Config, error) {
	var cfg Config
	var links []Link
	var stack []linkFrame
	inLinks := false

	scanner := bufio.NewScanner(strings.NewReader(raw))
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := stripYAMLComment(scanner.Text())
		if strings.TrimSpace(line) == "" {
			continue
		}
		trimmed := strings.TrimSpace(line)
		indent := len(line) - len(strings.TrimLeft(line, " "))

		if indent == 0 {
			key, value, ok := strings.Cut(trimmed, ":")
			if !ok {
				return Config{}, fmt.Errorf("line %d: expected key: value", lineNo)
			}
			key = strings.TrimSpace(key)
			value = strings.TrimSpace(value)
			if key == "links" {
				inLinks = true
				continue
			}
			inLinks = false
			if err := setConfigScalar(&cfg, key, value); err != nil {
				return Config{}, fmt.Errorf("line %d: %w", lineNo, err)
			}
			continue
		}

		if !inLinks {
			return Config{}, fmt.Errorf("line %d: nested values are only supported under links", lineNo)
		}
		if strings.HasPrefix(trimmed, "- ") {
			parent := currentParent(stack, indent)
			if parent == nil {
				links = append(links, Link{})
				stack = pushLinkFrame(stack, indent, &links[len(links)-1])
			} else {
				parent.Links = append(parent.Links, Link{})
				stack = pushLinkFrame(stack, indent, &parent.Links[len(parent.Links)-1])
			}
			current := stack[len(stack)-1].link
			rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
			if rest == "" {
				continue
			}
			key, value, ok := strings.Cut(rest, ":")
			if !ok {
				return Config{}, fmt.Errorf("line %d: expected - key: value", lineNo)
			}
			if err := setLinkScalar(current, strings.TrimSpace(key), strings.TrimSpace(value)); err != nil {
				return Config{}, fmt.Errorf("line %d: %w", lineNo, err)
			}
			continue
		}
		current := currentParent(stack, indent)
		if current == nil {
			return Config{}, fmt.Errorf("line %d: link item must start with -", lineNo)
		}
		key, value, ok := strings.Cut(trimmed, ":")
		if !ok {
			return Config{}, fmt.Errorf("line %d: expected key: value", lineNo)
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "links" && value == "" {
			continue
		}
		if err := setLinkScalar(current, key, value); err != nil {
			return Config{}, fmt.Errorf("line %d: %w", lineNo, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return Config{}, err
	}
	if links != nil {
		cfg.Links = links
	}
	return cfg, nil
}

type linkFrame struct {
	indent int
	link   *Link
}

func currentParent(stack []linkFrame, indent int) *Link {
	for i := len(stack) - 1; i >= 0; i-- {
		if stack[i].indent < indent {
			return stack[i].link
		}
	}
	return nil
}

func pushLinkFrame(stack []linkFrame, indent int, link *Link) []linkFrame {
	for len(stack) > 0 && stack[len(stack)-1].indent >= indent {
		stack = stack[:len(stack)-1]
	}
	return append(stack, linkFrame{indent: indent, link: link})
}

func setConfigScalar(cfg *Config, key string, value string) error {
	switch key {
	case "name":
		cfg.Name = unquote(value)
	case "title":
		cfg.Title = unquote(value)
	case "bio":
		cfg.Bio = unquote(value)
	case "avatar":
		cfg.Avatar = unquote(value)
	case "base_url":
		cfg.BaseURL = unquote(value)
	case "template":
		cfg.Template = unquote(value)
	case "accent":
		cfg.Accent = unquote(value)
	case "footer":
		cfg.Footer = unquote(value)
	case "cache_seconds":
		ttl, err := strconv.Atoi(unquote(value))
		if err != nil {
			return err
		}
		cfg.CacheSeconds = ttl
	default:
		return fmt.Errorf("unknown key %q", key)
	}
	return nil
}

func setLinkScalar(link *Link, key string, value string) error {
	switch key {
	case "title":
		link.Title = unquote(value)
	case "url":
		link.URL = unquote(value)
	case "icon":
		link.Icon = unquote(value)
	case "rel":
		link.Rel = unquote(value)
	case "featured":
		parsed, err := strconv.ParseBool(unquote(value))
		if err != nil {
			return err
		}
		link.Featured = parsed
	case "open":
		parsed, err := strconv.ParseBool(unquote(value))
		if err != nil {
			return err
		}
		link.Open = parsed
	default:
		return fmt.Errorf("unknown link key %q", key)
	}
	return nil
}

func stripYAMLComment(line string) string {
	var quote rune
	for i, r := range line {
		switch r {
		case '\'', '"':
			if quote == 0 {
				quote = r
			} else if quote == r {
				quote = 0
			}
		case '#':
			if quote == 0 {
				return strings.TrimRight(line[:i], " ")
			}
		}
	}
	return line
}

func unquote(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
			return value[1 : len(value)-1]
		}
	}
	return value
}
