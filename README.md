<p align="center">
  <img src="brand/swarm-icons-logo.svg" alt="Swarm Icons logo" width="160" height="160">
</p>

<h1 align="center">Hugo Swarm Icons</h1>

<p align="center">
  <a href="https://frostybee.github.io/hugo-swarm-icons/"><img src="https://img.shields.io/badge/Live%20demo-hugo--swarm--icons-1f6feb" alt="Live demo"></a>
  <a href="https://frostybee.github.io/go-swarm-icons/docs/guides/hugo-integration/"><img src="https://img.shields.io/badge/Documentation-integration%20guide-1f6feb" alt="Documentation"></a>
  <a href="https://github.com/frostybee/hugo-swarm-icons/actions/workflows/ci.yml"><img src="https://github.com/frostybee/hugo-swarm-icons/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
  <img src="https://img.shields.io/badge/Hugo-%E2%89%A50.128.0-ff4088" alt="Hugo Version">
</p>

A Hugo module that inlines SVG icons exported by the [go-swarm-icons](https://github.com/frostybee/go-swarm-icons) CLI and brings the library's fluent manipulation API (resize, rotate, flip, colors, opacity, class, title) to Hugo templates and Markdown.

<p align="center">
  <strong><a href="https://frostybee.github.io/hugo-swarm-icons/">Try the live demo</a></strong>
</p>

Requires Hugo 0.128.0 or later. No icons are bundled; export the sets your site uses with the `swarm-icons` CLI.

## Features

- **Manipulation parameters**: `size`, `width`, `height`, `rotate`, `flip`, `fill`, `stroke`, `stroke-width`, `opacity`, `class`, `style`, and `title`, available from both templates and Markdown.
- **viewBox-aware sizing**: give one dimension and the other derives from the icon's aspect ratio, CSS units preserved, so non-square icons never distort.
- **Composed transforms**: `rotate` and `flip` collapse into a single CSS `transform` declaration and append to any `style` you pass instead of overwriting it.
- **Accessibility defaults**: icons are `aria-hidden` and unfocusable by default; `title` or `aria-label` switches them to `role="img"` with an escaped `<title>` element.
- **Any Iconify set**: export from Tabler, Heroicons, Font Awesome, Material Design, Simple Icons, and 200+ others with the CLI, no Node.js required.
- **Only what you use**: icons are files in `assets/icons/`, so nothing unused ships with the site.
- **Loud on missing icons**: referencing an icon you have not exported fails the build with the exact export command needed to fix it.
- **Tested against the library**: CI byte-compares the template output against a golden file and against the real go-swarm-icons rendering.
- **Pure templates**: no JavaScript, no build step, no runtime dependencies.

## Install

As a Hugo module (recommended):

```bash
hugo mod init github.com/you/your-site   # once, if the site is not a module yet
hugo mod get github.com/frostybee/hugo-swarm-icons
```

```toml
# hugo.toml
[[module.imports]]
path = "github.com/frostybee/hugo-swarm-icons"
```

As a classic theme component (no Go toolchain needed):

```bash
git submodule add https://github.com/frostybee/hugo-swarm-icons themes/hugo-swarm-icons
```

```toml
# hugo.toml
theme = ["your-theme", "hugo-swarm-icons"]
```

## Get icons

Install the CLI and export the icons the site uses into `assets/icons/{prefix}/`, one subdirectory per icon set prefix:

```bash
go install github.com/frostybee/go-swarm-icons/cmd/swarm-icons@latest
swarm-icons json download tabler
swarm-icons icon export tabler home star menu-2 --dest assets/icons/tabler
```

Exported files are sanitized by the library (scripts, event handlers, and external references stripped), so they are safe to inline.

## Use

In templates:

```go-html-template
{{ partial "swarm-icon/icon.html" "tabler:home" }}

{{ partial "swarm-icon/icon.html" (dict
     "icon"   "tabler:home"
     "size"   32
     "rotate" 90
     "flip"   "h"
     "fill"   "currentColor"
     "class"  "nav-icon"
     "title"  "Home") }}
```

In Markdown, via the shortcode:

```text
{{< icon "tabler:home" >}}
{{< icon icon="tabler:arrow-up" size="20" rotate="45" class="nav-icon" title="Up" >}}
```

A missing icon fails the build with the exact `swarm-icons icon export` command needed to fix it.

## Parameters

| Parameter | Example | Behavior |
|---|---|---|
| `icon` | `"tabler:home"` | Required. `prefix:name`, loaded from `assets/icons/{prefix}/{name}.svg`. |
| `size` | `32` | Sets `width` and `height` to the same value. |
| `width` | `"1.5em"` | Sets width; height is derived from the viewBox aspect ratio, unit preserved. `"auto"` uses the viewBox dimensions; `"unset"`/`"none"` removes both. |
| `height` | `"48"` | Same as `width`, driving from the height side. |
| `fill` | `"crimson"` | Sets the root `fill` attribute. |
| `stroke` | `"currentColor"` | Sets the root `stroke` attribute. |
| `stroke-width` | `"1.5"` | Sets the root `stroke-width` attribute (also accepted as `strokewidth`). |
| `opacity` | `0.5` | Sets the root `opacity` attribute. |
| `rotate` | `45` | Appends `rotate(45deg)` to the `transform` declaration in `style`. |
| `flip` | `"h"`, `"v"`, `"both"` | Appends `scaleX(-1)`, `scaleY(-1)`, or `scale(-1, -1)`. Combines with `rotate` into one declaration (rotate first). |
| `class` | `"nav-icon"` | Space-concatenated onto any existing class. |
| `style` | `"color: red"` | Overwrites the style attribute; `rotate`/`flip` then append to it. |
| `title` | `"Home"` | Prepends an escaped `<title>` element and sets `role="img"`. |
| `attrs` | `(dict "id" "x")` | Partial only. Arbitrary extra attributes, highest precedence; empty values are skipped. |

## Accessibility

Icons are decorative by default: the partial adds `aria-hidden="true"` and `focusable="false"` when no labeling is present, matching the library renderer's defaults. Passing `title`, or `aria-label`/`aria-labelledby`/`role` via `attrs`, switches the icon to labeled mode (`role="img"`, no `aria-hidden`).

One deliberate divergence from the Go library: in Go, `Icon.Title()` only adds the `<title>` element and ARIA wiring is a separate renderer step. A template call has no such second step, so here `title` implies `role="img"` directly.

## How it works

Every fluent method in go-swarm-icons operates only on the root `<svg>` element's attributes, never on the icon's inner markup. The partial exploits that: it parses the root tag of the exported file, applies the same attribute algorithms the library uses (`icon_transform.go` semantics for dimensions and transforms, `applyARIA` semantics for accessibility), and re-serializes with ASCII-sorted attribute order, exactly like the library's `renderAttributes`. Output is deterministic, which the golden test relies on.

## Development

The `demo-site/` directory doubles as demo and test fixture, driven by `demo-site/data/test_cases.yaml`. Under `demo-site/assets/icons/`, five SVGs (`home`, `star`, `arrow-up`, `banner`, and lucide's `heart`) are hand-authored golden-test fixtures shaped like real CLI exports (including one non-square viewBox to exercise the aspect-ratio math); do not regenerate them. The rest are real CLI exports used only by the landing page. None of them are a bundled icon set, and none are mounted into consuming sites.

```bash
cd demo-site
hugo server            # browse the demo
```

A built copy is published to GitHub Pages from `main` by `.github/workflows/deploy-demo.yaml`: [frostybee.github.io/hugo-swarm-icons](https://frostybee.github.io/hugo-swarm-icons/).

CI builds the demo site and byte-compares `tests/manipulations/index.html` against `test/golden/manipulations.html`. After an intentional behavior change, regenerate the golden file, inspect the diff, and commit it:

```bash
hugo --source demo-site --destination ../public-test --minify=false
cp public-test/tests/manipulations/index.html test/golden/manipulations.html
git diff test/golden/manipulations.html
```

A second CI job guards against drift from the library itself: `test/parity/` renders the same fixture cases through the real go-swarm-icons API and byte-compares against the golden file, compensating only for the documented `title` divergence. Run it locally with:

```bash
cd test/parity
go test ./...
```

Because both jobs compare against the same golden file, regenerating it after a behavior change re-validates the partial and the library parity together. A parity failure with an unchanged partial means the library's semantics moved and the templates need to follow.

## Scope and limitations

The module ports everything that operates on a rendered icon. The library's runtime machinery has no template equivalent, so the following do not carry over:

- **Aliases and fallback icons**: there is no manager holding alias maps, and no fallback when a lookup fails, so a missing icon fails the build instead. For a static site a hard build-time failure is usually what you want, but it is a behavioral difference.
- **Runtime providers**: icons must be exported to disk before the build; there is no chain provider and no on-demand fetching from the Iconify API. In library terms, the module behaves like a directory provider only.
- **Renderer configuration layers**: the library's five-layer attribute merge (icon, global defaults, per-prefix, per-suffix, caller) reduces to three here: file attributes, then call parameters, then `attrs`.
- **Sprite sheets**: not yet available in template form; see the roadmap. Until then, the pre-build Go program approach in the [go-swarm-icons Hugo guide](https://frostybee.github.io/go-swarm-icons/docs/guides/hugo-integration/) covers it.

## Roadmap

- Sprite-sheet partial (template-side counterpart of `SpriteCollector`)
- Alias support via a site data file

## License

[MIT](LICENSE). Icon files exported from third-party icon sets keep their upstream licenses.
