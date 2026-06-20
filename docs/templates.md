# Templates

Mini-Link ships multiple static templates. They use the same HTML structure and config schema, so changing the visual style does not change deployment or caching behavior.

## Available Templates

- `classic`: clean, high-contrast professional link page.
- `glass`: Apple Liquid Glass-inspired surface with translucent panels and soft depth. This is still plain CSS, no JavaScript and no external assets.
- `terminal`: developer-focused dark terminal style using system monospace fonts.

## Configuration

```yaml
template: glass
```

```json
{
  "template": "terminal"
}
```

```env
MINI_LINK_TEMPLATE=classic
```

Unknown template names fail validation. Supported values are `classic`, `glass`, and `terminal`.

## Performance Rules

Templates must keep:

- no external CSS, JavaScript, fonts, or images
- stable avatar and link dimensions to avoid layout shift
- visible focus states
- at least 54px link rows for mobile tap targets
- readable contrast for text and accent colors
