package eigenkey

import "github.com/ccmonky/typemap"

func init() {
	typemap.MustRegisterType[KeyPostFunc]()
	typemap.MustRegisterType[HTTPRequestEigenkeyGen]()
	typemap.MustRegisterType[*HTTPRequestEigenkeyExtractor](typemap.WithDependencies([]string{
		typemap.GetTypeIdString[KeyPostFunc](),
		typemap.GetTypeIdString[HTTPRequestEigenkeyGen](),
	}))
}
