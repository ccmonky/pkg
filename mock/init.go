package mock

import (
	"context"

	"github.com/ccmonky/typemap"
)

func init() {
	typemap.MustRegisterType[Matcher]()
	typemap.MustRegisterType[ResponseMocker]()

	generators := []ResponseMocker{
		new(TransparentResponseMocker),
		new(ResponseMockerFromURL),
		new(ResponseMockerBuilder),
	}
	for _, gen := range generators {
		typemap.MustRegister[ResponseMocker](context.Background(), gen.ID(), gen)
	}
}
