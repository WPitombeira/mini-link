package icons

import (
	"sort"
)

type Icon struct {
	Name   string
	Label  string
	Source string
	SVG    string
}

var catalog = map[string]Icon{
	"code": {
		Name:   "code",
		Label:  "Code",
		Source: "Custom MIT",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="m8.7 16.9-5-4.9 5-4.9 1.4 1.5L6.6 12l3.5 3.4-1.4 1.5Zm6.6 0-1.4-1.5 3.5-3.4-3.5-3.4 1.4-1.5 5 4.9-5 4.9Zm-3.9 2.6-1.9-.6 3.1-14.4 1.9.6-3.1 14.4Z"/></svg>`,
	},
	"contact": {
		Name:   "contact",
		Label:  "Contact",
		Source: "Custom MIT",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M4 5h16a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V7a2 2 0 0 1 2-2Zm0 3.2V17h16V8.2l-8 5.3-8-5.3Zm1.4-1.2 6.6 4.4L18.6 7H5.4Z"/></svg>`,
	},
	"github": {
		Name:   "github",
		Label:  "GitHub",
		Source: "Simple Icons, CC0 1.0. GitHub is a trademark of GitHub, Inc.",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 .3A12 12 0 0 0 8.2 23.7c.6.1.8-.3.8-.6v-2.2c-3.3.7-4-1.4-4-1.4-.5-1.3-1.2-1.7-1.2-1.7-1-.7.1-.7.1-.7 1.1.1 1.7 1.2 1.7 1.2 1 .1.6 2.6 4.2 1.9.1-.7.4-1.2.7-1.5-2.7-.3-5.5-1.3-5.5-5.9 0-1.3.5-2.4 1.2-3.2-.1-.3-.5-1.6.1-3.2 0 0 1-.3 3.3 1.2a11.3 11.3 0 0 1 6 0c2.3-1.5 3.3-1.2 3.3-1.2.6 1.6.2 2.9.1 3.2.8.8 1.2 1.9 1.2 3.2 0 4.6-2.8 5.6-5.5 5.9.4.4.8 1.1.8 2.2v3.2c0 .3.2.7.8.6A12 12 0 0 0 12 .3Z"/></svg>`,
	},
	"globe": {
		Name:   "globe",
		Label:  "Globe",
		Source: "Custom MIT",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20Zm6.9 9h-3.1a14.5 14.5 0 0 0-1.3-5.2A8 8 0 0 1 18.9 11ZM12 4.1c.7 1 1.5 3.3 1.8 6.9h-3.6C10.5 7.4 11.3 5.1 12 4.1ZM4.1 13h3.1c.2 2 .6 3.8 1.3 5.2A8 8 0 0 1 4.1 13Zm3.1-2H4.1a8 8 0 0 1 4.4-5.2A14.5 14.5 0 0 0 7.2 11Zm4.8 8.9c-.7-1-1.5-3.3-1.8-6.9h3.6c-.3 3.6-1.1 5.9-1.8 6.9Zm3.5-1.7c.7-1.4 1.1-3.2 1.3-5.2h3.1a8 8 0 0 1-4.4 5.2Z"/></svg>`,
	},
	"linkedin": {
		Name:   "linkedin",
		Label:  "LinkedIn",
		Source: "Simple Icons, CC0 1.0. LinkedIn is a trademark of LinkedIn Corporation.",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M20.4 20.4h-3.6v-5.6c0-1.3 0-3-1.8-3s-2.1 1.4-2.1 2.9v5.7H9.3V9h3.4v1.6h.1c.5-.9 1.6-1.8 3.3-1.8 3.6 0 4.2 2.3 4.2 5.4v6.2ZM5.3 7.4a2.1 2.1 0 1 1 0-4.2 2.1 2.1 0 0 1 0 4.2Zm1.8 13H3.5V9h3.6v11.4ZM22.2 0H1.8C.8 0 0 .8 0 1.7v20.6c0 .9.8 1.7 1.8 1.7h20.4c1 0 1.8-.8 1.8-1.7V1.7c0-.9-.8-1.7-1.8-1.7Z"/></svg>`,
	},
	"link": {
		Name:   "link",
		Label:  "Link",
		Source: "Custom MIT",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M10.6 13.4a1 1 0 0 1 0-1.4l2.8-2.8a1 1 0 1 1 1.4 1.4L12 13.4a1 1 0 0 1-1.4 0Zm-5.7 5.7a5 5 0 0 1 0-7.1l3-3a5 5 0 0 1 7.8.9l-1.7 1a3 3 0 0 0-4.7-.5l-3 3a3 3 0 1 0 4.3 4.2l1.7-1.7a1 1 0 1 1 1.4 1.5L12 19.1a5 5 0 0 1-7.1 0Zm3.4-5 1.7-1a3 3 0 0 0 4.7.5l3-3a3 3 0 1 0-4.3-4.2l-1.7 1.7a1 1 0 0 1-1.4-1.5L12 4.9a5 5 0 1 1 7.1 7.1l-3 3a5 5 0 0 1-7.8-.9Z"/></svg>`,
	},
	"x": {
		Name:   "x",
		Label:  "X",
		Source: "Simple Icons, CC0 1.0. X is a trademark of X Corp.",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M18.9 1.2h3.7l-8.1 9.3 9.5 12.6h-7.4l-5.8-7.6-6.7 7.6H.4l8.7-9.9L0 1.2h7.6l5.3 7 6-7Zm-1.3 19.7h2L6.5 3.3H4.3l13.3 17.6Z"/></svg>`,
	},
}

func Get(name string) (Icon, bool) {
	icon, ok := catalog[name]
	return icon, ok
}

func All() []Icon {
	out := make([]Icon, 0, len(catalog))
	for _, icon := range catalog {
		out = append(out, icon)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
