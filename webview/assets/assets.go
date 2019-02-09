package assets

import (
	"github.com/Oppodelldog/plainci/config"
	"github.com/go-playground/statics/static"
)

var Templates *static.Files
var Images *static.Files

func init() {
	var err error

	config := &static.Config{
		UseStaticFiles: config.UseStaticFiles,
		FallbackToDisk: false,
		AbsPkgPath:     config.AbsoluteAssetsPath,
	}

	Templates, err = newStaticTemplates(config)
	if err != nil {
		panic(err)
	}

	Images, err = newStaticImages(config)
	if err != nil {
		panic(err)
	}
}
