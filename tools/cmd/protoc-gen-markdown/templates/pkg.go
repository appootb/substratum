package templates

import (
	"text/template"

	"github.com/appootb/substratum/v2/tools/cmd/protoc-gen-markdown/templates/markdown"
	"github.com/appootb/substratum/v2/tools/cmd/protoc-gen-markdown/templates/readme"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

type RegisterFn func(tpl *template.Template, params pgs.Parameters)

func ProtoTemplate(params pgs.Parameters) []*template.Template {
	return []*template.Template{
		makeTemplate("md", markdown.Register, params),
	}
}

func ReadMeTemplate(params pgs.Parameters) *template.Template {
	return makeTemplate("readme", readme.Register, params)
}

func FilePathFor(tpl *template.Template) func(f pgs.File, ctx pgsgo.Context, tpl *template.Template) *pgs.FilePath {
	return func(f pgs.File, ctx pgsgo.Context, _ *template.Template) *pgs.FilePath {
		out := ctx.OutputPath(f)
		if out == "" {
			return nil
		}
		out = out.SetExt("." + tpl.Name())
		return &out
	}
}

func makeTemplate(ext string, fn RegisterFn, params pgs.Parameters) *template.Template {
	tpl := template.New(ext)
	fn(tpl, params)
	return tpl
}
