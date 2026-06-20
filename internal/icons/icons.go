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
	"briefcase": {
		Name:   "briefcase",
		Label:  "Briefcase",
		Source: "Heroicons, MIT License",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M9 4a3 3 0 0 0-3 3v1H4a2 2 0 0 0-2 2v8.5A2.5 2.5 0 0 0 4.5 21h15a2.5 2.5 0 0 0 2.5-2.5V10a2 2 0 0 0-2-2h-2V7a3 3 0 0 0-3-3H9Zm7 4H8V7a1 1 0 0 1 1-1h6a1 1 0 0 1 1 1v1Zm-6 5H4v-3h16v3h-6v-1h-4v1Zm0 2v1h4v-1h6v3.5a.5.5 0 0 1-.5.5h-15a.5.5 0 0 1-.5-.5V15h6Z"/></svg>`,
	},
	"calendar": {
		Name:   "calendar",
		Label:  "Calendar",
		Source: "Heroicons, MIT License",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M7 2a1 1 0 0 1 1 1v1h8V3a1 1 0 1 1 2 0v1h1a3 3 0 0 1 3 3v12a3 3 0 0 1-3 3H5a3 3 0 0 1-3-3V7a3 3 0 0 1 3-3h1V3a1 1 0 0 1 1-1Zm13 8H4v9a1 1 0 0 0 1 1h14a1 1 0 0 0 1-1v-9ZM5 6a1 1 0 0 0-1 1v1h16V7a1 1 0 0 0-1-1H5Z"/></svg>`,
	},
	"check": {
		Name:   "check",
		Label:  "Check",
		Source: "Heroicons, MIT License",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M9.2 16.6 4.9 12.3l1.4-1.4 2.9 2.9 8.5-8.5 1.4 1.4-9.9 9.9Z"/></svg>`,
	},
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
	"copy": {
		Name:   "copy",
		Label:  "Copy",
		Source: "Heroicons, MIT License",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M8 7a3 3 0 0 1 3-3h7a3 3 0 0 1 3 3v7a3 3 0 0 1-3 3h-1v1a3 3 0 0 1-3 3H6a3 3 0 0 1-3-3v-8a3 3 0 0 1 3-3h2Zm2 10h4a3 3 0 0 0 3-3V9h1a1 1 0 0 1 1 1v7a1 1 0 0 1-1 1h-7a1 1 0 0 1-1-1Zm-4-8a1 1 0 0 0-1 1v8a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1v-8a1 1 0 0 0-1-1H6Z"/></svg>`,
	},
	"external-link": {
		Name:   "external-link",
		Label:  "External Link",
		Source: "Heroicons, MIT License",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M14 3h7v7h-2V6.4l-8.3 8.3-1.4-1.4L17.6 5H14V3ZM5 5h6v2H5v12h12v-6h2v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V7a2 2 0 0 1 2-2Z"/></svg>`,
	},
	"file-text": {
		Name:   "file-text",
		Label:  "File Text",
		Source: "Heroicons, MIT License",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M6 2h8l6 6v12a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2Zm7 2H6v16h12V9h-5V4Zm2 1.4V7h1.6L15 5.4ZM8 12h8v2H8v-2Zm0 4h8v2H8v-2Zm0-8h3v2H8V8Z"/></svg>`,
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
	"home": {
		Name:   "home",
		Label:  "Home",
		Source: "Heroicons, MIT License",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="m12 3 9 8-1.3 1.5L18 11v9a1 1 0 0 1-1 1h-4v-6h-2v6H7a1 1 0 0 1-1-1v-9l-1.7 1.5L3 11l9-8Zm0 2.7-4 3.6V19h1v-6h6v6h1V9.3l-4-3.6Z"/></svg>`,
	},
	"instagram": {
		Name:   "instagram",
		Label:  "Instagram",
		Source: "Simple Icons, CC0 1.0. Instagram is a trademark of Instagram, LLC.",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M7 2h10a5 5 0 0 1 5 5v10a5 5 0 0 1-5 5H7a5 5 0 0 1-5-5V7a5 5 0 0 1 5-5Zm0 2a3 3 0 0 0-3 3v10a3 3 0 0 0 3 3h10a3 3 0 0 0 3-3V7a3 3 0 0 0-3-3H7Zm5 4a4 4 0 1 1 0 8 4 4 0 0 1 0-8Zm0 2a2 2 0 1 0 0 4 2 2 0 0 0 0-4Zm4.8-3.4a.9.9 0 1 1 0 1.8.9.9 0 0 1 0-1.8Z"/></svg>`,
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
	"mail": {
		Name:   "mail",
		Label:  "Mail",
		Source: "Heroicons, MIT License",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M4 5h16a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V7a2 2 0 0 1 2-2Zm0 3.2V17h16V8.2l-8 5.3-8-5.3Zm1.4-1.2 6.6 4.4L18.6 7H5.4Z"/></svg>`,
	},
	"map-pin": {
		Name:   "map-pin",
		Label:  "Map Pin",
		Source: "Heroicons, MIT License",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 2a7 7 0 0 1 7 7c0 5.2-7 13-7 13S5 14.2 5 9a7 7 0 0 1 7-7Zm0 2a5 5 0 0 0-5 5c0 2.9 3.1 7.4 5 9.9 1.9-2.5 5-7 5-9.9a5 5 0 0 0-5-5Zm0 2.5A2.5 2.5 0 1 1 12 11a2.5 2.5 0 0 1 0-5Z"/></svg>`,
	},
	"moon": {
		Name:   "moon",
		Label:  "Moon",
		Source: "Heroicons, MIT License",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M20.6 14.7A8.5 8.5 0 0 1 9.3 3.4a9 9 0 1 0 11.3 11.3ZM12 21a7 7 0 0 1-4.2-12.6 10.5 10.5 0 0 0 7.8 7.8A7 7 0 0 1 12 21Z"/></svg>`,
	},
	"phone": {
		Name:   "phone",
		Label:  "Phone",
		Source: "Heroicons, MIT License",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M6.6 2.7 10 6.1 7.8 9c.8 1.7 2.5 3.4 4.2 4.2l2.9-2.2 3.4 3.4-1.5 4.4c-.3.8-1 1.3-1.8 1.2C8.6 19.4 4.6 15.4 4 9c-.1-.8.4-1.6 1.2-1.8l1.4-.5Z"/></svg>`,
	},
	"rss": {
		Name:   "rss",
		Label:  "RSS",
		Source: "Heroicons, MIT License",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M5 17a2 2 0 1 1 0 4 2 2 0 0 1 0-4Zm-2-6a10 10 0 0 1 10 10h-3a7 7 0 0 0-7-7v-3Zm0-6a16 16 0 0 1 16 16h-3A13 13 0 0 0 3 8V5Z"/></svg>`,
	},
	"sun": {
		Name:   "sun",
		Label:  "Sun",
		Source: "Heroicons, MIT License",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M11 2h2v3h-2V2Zm0 17h2v3h-2v-3ZM4.2 5.6l1.4-1.4 2.1 2.1-1.4 1.4-2.1-2.1Zm12.1 12.1 1.4-1.4 2.1 2.1-1.4 1.4-2.1-2.1ZM2 11h3v2H2v-2Zm17 0h3v2h-3v-2ZM4.2 18.4l2.1-2.1 1.4 1.4-2.1 2.1-1.4-1.4ZM16.3 6.3l2.1-2.1 1.4 1.4-2.1 2.1-1.4-1.4ZM12 7a5 5 0 1 1 0 10 5 5 0 0 1 0-10Zm0 2a3 3 0 1 0 0 6 3 3 0 0 0 0-6Z"/></svg>`,
	},
	"terminal": {
		Name:   "terminal",
		Label:  "Terminal",
		Source: "Heroicons, MIT License",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M3 4h18a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2Zm0 2v12h18V6H3Zm3 3.4L7.4 8 11 11.6 7.4 15 6 13.6l2.2-2L6 9.4ZM12 14h6v2h-6v-2Z"/></svg>`,
	},
	"user": {
		Name:   "user",
		Label:  "User",
		Source: "Heroicons, MIT License",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 12a5 5 0 1 1 0-10 5 5 0 0 1 0 10Zm0-2a3 3 0 1 0 0-6 3 3 0 0 0 0 6Zm0 3c4.4 0 8 2.7 8 6v2H4v-2c0-3.3 3.6-6 8-6Zm0 2c-3.2 0-6 1.9-6 4h12c0-2.1-2.8-4-6-4Z"/></svg>`,
	},
	"whatsapp": {
		Name:   "whatsapp",
		Label:  "WhatsApp",
		Source: "Simple Icons, CC0 1.0. WhatsApp is a trademark of WhatsApp LLC.",
		SVG:    `<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 2a9.7 9.7 0 0 0-8.4 14.6L2.4 22l5.5-1.2A9.7 9.7 0 1 0 12 2Zm0 2a7.7 7.7 0 1 1-3.4 14.6l-.3-.2-3 .7.7-2.9-.2-.3A7.7 7.7 0 0 1 12 4Zm-3.1 4.2c.2-.4.4-.4.7-.4h.5c.2 0 .4.1.5.4l.7 1.7c.1.2.1.4 0 .6l-.5.7c-.1.1-.2.3 0 .5.4.8 1.5 2.1 2.9 2.6.2.1.4.1.5-.1l.8-1c.2-.2.4-.2.6-.1l1.8.9c.2.1.4.3.3.5-.1.7-.7 1.6-1.5 1.8-.8.2-2.9.1-5.3-2.1-2.1-1.9-3.3-4.4-3.4-5.2-.1-.8.4-1.7.8-2Z"/></svg>`,
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
