package d2

import (
	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"context"
	"embed"
	"github.com/pkg/errors"
	"os"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
	"oss.terrastruct.com/d2/lib/log"
	"oss.terrastruct.com/d2/lib/textmeasure"
	"oss.terrastruct.com/util-go/go2"
)

//go:embed templates/*
var d2FS embed.FS

type RenderOptions struct {
	Logger slog.Logger
}

type RenderOption func(*RenderOptions)

func Render(ctx context.Context, template string, opts ...RenderOption) ([]byte, error) {
	options := &RenderOptions{
		Logger: slog.Make(sloghuman.Sink(os.Stderr)).Named("default"),
	}
	for _, opt := range opts {
		opt(options)
	}
	ctx = log.With(ctx, options.Logger)
	ruler, _ := textmeasure.NewRuler()
	layoutResolver := func(engine string) (d2graph.LayoutGraph, error) {
		return d2dagrelayout.DefaultLayout, nil
	}
	renderOpts := &d2svg.RenderOpts{
		Pad:     go2.Pointer(int64(5)),
		ThemeID: &d2themescatalog.GrapeSoda.ID,
	}
	compileOpts := &d2lib.CompileOptions{
		LayoutResolver: layoutResolver,
		Ruler:          ruler,
	}
	data, err := d2FS.ReadFile(template)
	if err != nil {
		return nil, errors.WithMessagef(err, "read template %s failed", template)
	}
	diagram, _, err := d2lib.Compile(ctx, string(data), compileOpts, renderOpts)
	if err != nil {
		return nil, errors.WithMessagef(err, "compile failed")
	}
	out, err := d2svg.Render(diagram, renderOpts)
	if err != nil {
		return nil, errors.WithMessagef(err, "render failed")
	}
	return out, nil
}
