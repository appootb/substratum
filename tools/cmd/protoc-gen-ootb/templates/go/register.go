package golang

import (
	"text/template"

	"github.com/appootb/substratum/v2/tools/cmd/protoc-gen-ootb/templates/shared"
	pgs "github.com/lyft/protoc-gen-star/v2"
)

func Register(tpl *template.Template, params pgs.Parameters) {
	shared.Register(tpl, params)
	//
	template.Must(tpl.Parse(fileTpl))
	template.Must(tpl.New("service").Parse(serviceTpl))
	template.Must(tpl.New("wrapper").Parse(wrapperTpl))
	template.Must(tpl.New("method").Parse(methodTpl))
	template.Must(tpl.New("streaming").Parse(webStreamTpl))
	template.Must(tpl.New("clientStream").Parse(clientStreamTpl))
	template.Must(tpl.New("serverStream").Parse(serverStreamTpl))
}
