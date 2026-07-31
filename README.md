# hugo-swarm-icons

A Hugo module that inlines SVG icons exported by the [go-swarm-icons](https://github.com/frostybee/go-swarm-icons) CLI and brings the library's fluent manipulation API (resize, rotate, flip, colors, opacity, class, title) to Hugo templates and Markdown.

Requires Hugo 0.128.0 or later. No icons are bundled; export the sets your site uses with the `swarm-icons` CLI.

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

The `demo-site/` directory doubles as demo and test fixture, driven by `demo-site/data/test_cases.yaml`:

```bash
cd demo-site
hugo server            # browse the demo
```

CI builds the demo site and byte-compares `tests/manipulations/index.html` against `test/golden/manipulations.html`. After an intentional behavior change, regenerate the golden file, inspect the diff, and commit it:

```bash
hugo --source demo-site --destination ../public-test --minify=false
cp public-test/tests/manipulations/index.html test/golden/manipulations.html
git diff test/golden/manipulations.html
```

## Roadmap

- Sprite-sheet partial (template-side counterpart of `SpriteCollector`)
- Go parity test rendering fixtures through the real library API

## License

[MIT](LICENSE). Icon files exported from third-party icon sets keep their upstream licenses.
