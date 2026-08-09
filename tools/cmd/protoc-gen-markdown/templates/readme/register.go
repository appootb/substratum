package readme

import (
	"text/template"

	"github.com/appootb/substratum/v2/tools/cmd/protoc-gen-markdown/templates/shared"
	pgs "github.com/lyft/protoc-gen-star/v2"
)

func Register(tpl *template.Template, params pgs.Parameters) {
	shared.Register(tpl, params)
	template.Must(tpl.Parse(fileTpl))
}
