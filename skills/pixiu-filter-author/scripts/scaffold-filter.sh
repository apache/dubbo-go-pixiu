#!/usr/bin/env bash
#
# scaffold-filter.sh — generate a minimal pixiu HTTP filter package.
#
# Usage:
#   scaffold-filter.sh <name> [--proxy]
#
# Without --proxy: creates pkg/filter/<name>/<name>.go (single file).
# With --proxy:    creates pkg/filter/http/<name>/{plugin.go,filter.go,config.go}.
#
# Run from the repo root, or pass PIXIU_ROOT=<path>.
#
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "usage: $0 <name> [--proxy]" >&2
  exit 1
fi

NAME="$1"
KIND_PREFIX="dgp.filter.http"
PROXY=0
if [[ "${2:-}" == "--proxy" ]]; then
  PROXY=1
fi

if [[ ! "$NAME" =~ ^[a-z][a-z0-9_]*$ ]]; then
  echo "error: name must be lowercase_snake_case" >&2
  exit 1
fi

ROOT="${PIXIU_ROOT:-$(pwd)}"
if [[ ! -f "$ROOT/go.mod" ]]; then
  echo "error: $ROOT does not look like the pixiu repo root (no go.mod). Pass PIXIU_ROOT=<path>." >&2
  exit 1
fi

if [[ $PROXY -eq 1 ]]; then
  DIR="$ROOT/pkg/filter/http/$NAME"
else
  DIR="$ROOT/pkg/filter/$NAME"
fi

if [[ -e "$DIR" ]]; then
  echo "error: $DIR already exists; aborting to avoid overwrite" >&2
  exit 1
fi

mkdir -p "$DIR"

KIND="${KIND_PREFIX}.${NAME}"
PKG="$NAME"

if [[ $PROXY -eq 1 ]]; then
  # Split-file layout.
  cat > "$DIR/plugin.go" <<EOF
package ${PKG}

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

// Kind is the unique identifier for this filter in yaml.
const Kind = "${KIND}"

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type Plugin struct{}

func (p *Plugin) Kind() string { return Kind }

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &Config{}}, nil
}

type FilterFactory struct {
	cfg *Config
}

func (f *FilterFactory) Config() any   { return f.cfg }
func (f *FilterFactory) Apply() error  { return nil }

func (f *FilterFactory) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {
	inst := &Filter{cfg: f.cfg.DeepCopy()}
	chain.AppendDecodeFilters(inst)
	return nil
}
EOF

  cat > "$DIR/filter.go" <<EOF
package ${PKG}

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

type Filter struct {
	cfg *Config
}

func (f *Filter) Decode(ctx *http.HttpContext) filter.FilterStatus {
	// TODO: implement request-phase behavior.
	return filter.Continue
}
EOF

  cat > "$DIR/config.go" <<EOF
package ${PKG}

// Config holds the yaml-bound configuration for the ${NAME} filter.
type Config struct {
	// TODO: add yaml-tagged fields, e.g.:
	// Enabled bool \`yaml:"enabled" json:"enabled" mapstructure:"enabled"\`
}

// DeepCopy returns an independent per-request copy of Config.
func (c *Config) DeepCopy() *Config {
	if c == nil {
		return nil
	}
	cp := *c
	// TODO: explicitly clone any slices/maps/pointers added to Config.
	return &cp
}
EOF

else
  # Single-file layout.
  cat > "$DIR/${NAME}.go" <<EOF
package ${PKG}

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

// Kind is the unique identifier for this filter in yaml.
const Kind = "${KIND}"

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	Plugin        struct{}
	FilterFactory struct{ cfg *Config }
	Filter        struct{ cfg *Config }
	Config        struct {
		// TODO: add yaml-tagged fields.
	}
)

func (p *Plugin) Kind() string { return Kind }
func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &Config{}}, nil
}

func (f *FilterFactory) Config() any  { return f.cfg }
func (f *FilterFactory) Apply() error { return nil }
func (f *FilterFactory) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {
	inst := &Filter{cfg: f.cfg.DeepCopy()}
	chain.AppendDecodeFilters(inst)
	return nil
}

func (f *Filter) Decode(ctx *http.HttpContext) filter.FilterStatus {
	// TODO: implement.
	return filter.Continue
}

// DeepCopy returns an independent per-request copy of Config.
func (c *Config) DeepCopy() *Config {
	if c == nil {
		return nil
	}
	cp := *c
	// TODO: explicitly clone any slices/maps/pointers added to Config.
	return &cp
}
EOF
fi

echo "scaffolded $DIR"
echo ""
echo "Next steps:"
echo "  1. Fill in Config fields + Apply() validation"
echo "  2. Implement Decode (and/or add Encode + chain.AppendEncodeFilters)"
echo "  3. Add blank import to pkg/pluginregistry/registry.go"
echo "  4. Add config block under http_filters in configs/conf.yaml"
echo "  5. (Optional) Write ${NAME}_test.go — only if your filter has branching logic worth testing; many existing filters (cors, csrf, jwt, httpproxy, ...) ship without tests."
