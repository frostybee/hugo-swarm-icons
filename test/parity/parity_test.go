// Package parity renders the shared fixture cases through the real
// go-swarm-icons library and byte-compares the output against the golden
// file produced by the Hugo partial. It exists to catch drift between the
// library's transform semantics and the template reimplementation.
//
// Single source of truth: demo-site/data/test_cases.yaml drives the demo
// page, the golden test page, and this test. The expected output is
// test/golden/manipulations.html, the same file the Hugo golden CI job
// diffs against.
//
// One deliberate divergence is compensated for below: in the partial,
// `title` implies role="img" and suppresses the decorative ARIA defaults,
// while the library's Icon.Title() only prepends the <title> element (ARIA
// wiring is the renderer's job). See "Accessibility" in the module README.
package parity

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"

	swarmicons "github.com/frostybee/go-swarm-icons"
	"gopkg.in/yaml.v3"
)

type testCase struct {
	ID     string         `yaml:"id"`
	Label  string         `yaml:"label"`
	Params map[string]any `yaml:"params"`
}

var repoRoot = filepath.Join("..", "..")

func loadCases(t *testing.T) []testCase {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repoRoot, "demo-site", "data", "test_cases.yaml"))
	if err != nil {
		t.Fatalf("read test_cases.yaml: %v", err)
	}
	var cases []testCase
	if err := yaml.Unmarshal(data, &cases); err != nil {
		t.Fatalf("parse test_cases.yaml: %v", err)
	}
	if len(cases) == 0 {
		t.Fatal("test_cases.yaml contains no cases")
	}
	return cases
}

var sectionRE = regexp.MustCompile(`<section data-case="([^"]+)">(.*?)</section>`)

func loadGolden(t *testing.T) map[string]string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repoRoot, "test", "golden", "manipulations.html"))
	if err != nil {
		t.Fatalf("read golden file: %v", err)
	}
	golden := make(map[string]string)
	for _, m := range sectionRE.FindAllStringSubmatch(string(data), -1) {
		golden[m[1]] = m[2]
	}
	if len(golden) == 0 {
		t.Fatal("no <section data-case=...> entries found in golden file")
	}
	return golden
}

func buildManager(t *testing.T) *swarmicons.IconManager {
	t.Helper()
	iconsDir := filepath.Join(repoRoot, "demo-site", "assets", "icons")
	entries, err := os.ReadDir(iconsDir)
	if err != nil {
		t.Fatalf("read icons dir: %v", err)
	}
	cfg := swarmicons.NewConfig()
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		p, err := swarmicons.NewDirectoryProvider(filepath.Join(iconsDir, e.Name()), swarmicons.WithRecursive(false))
		if err != nil {
			t.Fatalf("provider for %s: %v", e.Name(), err)
		}
		cfg.AddProvider(e.Name(), p)
	}
	mgr, err := cfg.Build()
	if err != nil {
		t.Fatalf("build manager: %v", err)
	}
	return mgr
}

// str renders a YAML scalar the way the partial's printf "%v" does.
func str(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case int:
		return strconv.Itoa(x)
	case float64:
		return strconv.FormatFloat(x, 'g', -1, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func toFloat(t *testing.T, v any) float64 {
	t.Helper()
	switch x := v.(type) {
	case int:
		return float64(x)
	case float64:
		return x
	case string:
		f, err := strconv.ParseFloat(x, 64)
		if err != nil {
			t.Fatalf("not a number: %v", v)
		}
		return f
	default:
		t.Fatalf("not a number: %v (%T)", v, v)
		return 0
	}
}

func toInt(t *testing.T, v any) int {
	t.Helper()
	switch x := v.(type) {
	case int:
		return x
	case float64:
		return int(x)
	default:
		t.Fatalf("not an int: %v (%T)", v, v)
		return 0
	}
}

// renderWithLibrary applies the case's params through the real library API in
// the same order the partial applies them: dimensions, fill/stroke/
// stroke-width/opacity, class, style, rotate, flip, title, attrs.
func renderWithLibrary(t *testing.T, mgr *swarmicons.IconManager, params map[string]any) string {
	t.Helper()
	ref, _ := params["icon"].(string)
	if ref == "" {
		t.Fatal("case has no icon param")
	}

	// Passthrough attrs go through Get as caller attributes so the renderer's
	// ARIA logic sees aria-label/role, exactly as it would in an application.
	var callerAttrs map[string]string
	if raw, ok := params["attrs"].(map[string]any); ok {
		callerAttrs = make(map[string]string, len(raw))
		for k, v := range raw {
			callerAttrs[k] = str(v)
		}
	}

	var icon *swarmicons.Icon
	var err error
	if callerAttrs != nil {
		icon, err = mgr.Get(ref, callerAttrs)
	} else {
		icon, err = mgr.Get(ref)
	}
	if err != nil {
		t.Fatalf("Get(%q): %v", ref, err)
	}

	// Dimensions: size wins outright; both width and height given means
	// literal writes with no derivation (the partial mirrors setDimension
	// only when a single side drives).
	if v, ok := params["size"]; ok {
		icon = icon.Size(toInt(t, v))
	} else {
		w, wok := params["width"]
		h, hok := params["height"]
		switch {
		case wok && hok:
			icon = icon.Attr(map[string]string{"width": str(w), "height": str(h)})
		case wok:
			icon = icon.Width(str(w))
		case hok:
			icon = icon.Height(str(h))
		}
	}

	if v, ok := params["fill"]; ok {
		icon = icon.Fill(str(v))
	}
	if v, ok := params["stroke"]; ok {
		icon = icon.Stroke(str(v))
	}
	if v, ok := params["stroke-width"]; ok {
		icon = icon.StrokeWidth(str(v))
	}
	if v, ok := params["opacity"]; ok {
		icon = icon.Opacity(toFloat(t, v))
	}
	if v, ok := params["class"]; ok {
		icon = icon.Class(str(v))
	}
	// Style overwrites like Attr; rotate/flip then append to its transform.
	if v, ok := params["style"]; ok {
		icon = icon.Attr(map[string]string{"style": str(v)})
	}
	if v, ok := params["rotate"]; ok {
		icon = icon.Rotate(toFloat(t, v))
	}
	if v, ok := params["flip"]; ok {
		icon = icon.Flip(str(v))
	}

	_, hasTitle := params["title"]
	if hasTitle {
		icon = icon.Title(str(params["title"]))
	}

	// Documented divergence: the partial's title implies role="img" and
	// suppresses the decorative defaults; the raw library Title() leaves
	// ARIA untouched. Rebuild with the adjusted attribute set (New copies
	// the map and skips sanitization, which already ran on load).
	if hasTitle {
		attrs := icon.Attributes()
		if attrs["role"] == "" {
			attrs["role"] = "img"
		}
		delete(attrs, "aria-hidden")
		delete(attrs, "focusable")
		icon = swarmicons.New(icon.Content(), attrs)
	}

	return icon.ToHTML()
}

func TestPartialMatchesLibrary(t *testing.T) {
	cases := loadCases(t)
	golden := loadGolden(t)
	mgr := buildManager(t)

	seen := make(map[string]bool, len(cases))
	for _, tc := range cases {
		seen[tc.ID] = true
		t.Run(tc.ID, func(t *testing.T) {
			want, ok := golden[tc.ID]
			if !ok {
				t.Fatalf("case %q has no entry in the golden file; regenerate it (see README Development)", tc.ID)
			}
			got := renderWithLibrary(t, mgr, tc.Params)
			if got != want {
				t.Errorf("library output diverges from Hugo partial\nlibrary: %s\npartial: %s", got, want)
			}
		})
	}

	for id := range golden {
		if !seen[id] {
			t.Errorf("golden file has case %q that is missing from test_cases.yaml", id)
		}
	}
}
