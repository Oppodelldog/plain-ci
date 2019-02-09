package assets

import (
	"github.com/Oppodelldog/plainci/config"
	"github.com/go-playground/statics/static"
)

var Templates *static.Files
var Images *static.Files

func init() {
	var err error

	cfg := &static.Config{
		UseStaticFiles: config.UseStaticFiles,
		FallbackToDisk: false,
		AbsPkgPath:     config.AbsoluteAssetsPath,
	}

	Templates, err = newStaticTemplates(cfg)
	if err != nil {
		panic(err)
	}

	Images, err = newStaticImages(cfg)
	if err != nil {
		panic(err)
	}
}
