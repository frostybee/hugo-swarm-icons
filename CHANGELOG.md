# Changelog

## v0.1.0

Initial release.

- `swarm-icon/icon.html` partial: inlines CLI-exported SVGs with the fluent manipulation parameters (size, width, height, fill, stroke, stroke-width, opacity, rotate, flip, class, style, title, attrs).
- `icon` shortcode for Markdown content, positional and named forms.
- Accessibility defaults matching the library renderer (`aria-hidden`/`focusable` when decorative, `role="img"` when labeled).
- Example site with a data-driven test-case matrix and a byte-exact golden test enforced in CI on Hugo 0.128.0 and current.
