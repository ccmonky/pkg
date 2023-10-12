package d2

import (
	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"context"
	"github.com/pkg/errors"
	"os"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2elklayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
	"oss.terrastruct.com/d2/lib/log"
	"oss.terrastruct.com/d2/lib/textmeasure"
	"oss.terrastruct.com/util-go/go2"
)

var ArchTemplate = `user: 客户端 {
  tooltip: API网关
}
app: xxx {
  link: https://baidu.com
  cluster-hangzhou: {
    node1
    node2
    redis: redis-xxx {
      shape: cylinder
      link: https://baidu.com
    }
    rdb: rdb-xxx {
      shape: cylinder
      link: https://baidu.com
    }
    node1 -> redis
    node2 -> redis
    redis -> rdb
  }
  cluster-beijing: {
    node1
    node2
    redis: redis-xxx {
      shape: cylinder
      link: https://baidu.com
    }
    rdb: rdb-xxx {
      shape: cylinder
      link: https://baidu.com
    }
    node1 -> redis
    node2 -> redis
    redis -> rdb
  }
  cluster-shenzhen: {
    node1
    node2
    redis: redis-xxx {
      shape: cylinder
      link: https://baidu.com
    }
    rdb: rdb-xxx {
      shape: cylinder
      link: https://baidu.com
    }
    node1 -> redis
    node2 -> redis
    redis -> rdb
  }
}
user -> app: http
`

type RenderOptions struct {
	Logger slog.Logger
}

type RenderOption func(*RenderOptions)

func Render(ctx context.Context, opts ...RenderOption) ([]byte, error) {
	options := &RenderOptions{
		Logger: slog.Make(sloghuman.Sink(os.Stderr)).Named("default"),
	}
	for _, opt := range opts {
		opt(options)
	}
	ctx = log.With(ctx, options.Logger)
	ruler, _ := textmeasure.NewRuler()
	layoutResolver := func(engine string) (d2graph.LayoutGraph, error) {
		return d2elklayout.DefaultLayout, nil
	}
	renderOpts := &d2svg.RenderOpts{
		Pad:     go2.Pointer(int64(5)),
		ThemeID: &d2themescatalog.GrapeSoda.ID,
	}
	compileOpts := &d2lib.CompileOptions{
		LayoutResolver: layoutResolver,
		Ruler:          ruler,
	}
	diagram, _, err := d2lib.Compile(ctx, ArchTemplate, compileOpts, renderOpts)
	if err != nil {
		return nil, errors.WithMessagef(err, "compile failed")
	}
	out, err := d2svg.Render(diagram, renderOpts)
	if err != nil {
		return nil, errors.WithMessagef(err, "render failed")
	}
	return out, nil
}
