package actions

import (
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
)

var r *render.Engine

var assetBox = packr.NewBox("../assets")

func init() {
	r = render.New(render.Options{
		AssetsBox: assetBox,
	})
}
