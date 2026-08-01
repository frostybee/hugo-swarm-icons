# Changelog

## Unreleased

- Go parity test (`test/parity/`): renders the shared fixture cases through the real go-swarm-icons library (v0.2.0) and byte-compares against the golden file, enforced in CI. Catches drift between the library's transform semantics and the template reimplementation.
- Live demo deployed to GitHub Pages from `main` (`deploy-demo.yaml`): [frostybee.github.io/hugo-swarm-icons](https://frostybee.github.io/hugo-swarm-icons/).

## v0.1.0

Initial release.

- `swarm-icon/icon.html` partial: inlines CLI-exported SVGs with the fluent manipulation parameters (size, width, height, fill, stroke, stroke-width, opacity, rotate, flip, class, style, title, attrs).
- `icon` shortcode for Markdown content, positional and named forms.
- Accessibility defaults matching the library renderer (`aria-hidden`/`focusable` when decorative, `role="img"` when labeled).
- Demo site with a data-driven test-case matrix and a byte-exact golden test enforced in CI on Hugo 0.128.0 and current.
