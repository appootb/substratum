package generator

import (
	"sort"

	"github.com/appootb/substratum/v2/proto/go/api"
	"github.com/appootb/substratum/v2/tools/cmd/protoc-gen-ootb/templates"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

const (
	OutOfTheBoxGenerator = "ootb"
	LangParam            = "lang"
)

type OutOfTheBox struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
}

func New() pgs.Module {
	return &OutOfTheBox{
		ModuleBase: &pgs.ModuleBase{},
	}
}

func (m *OutOfTheBox) InitContext(ctx pgs.BuildContext) {
	m.ModuleBase.InitContext(ctx)
	m.ctx = pgsgo.InitContext(ctx.Parameters())
}

func (m *OutOfTheBox) Name() string {
	return OutOfTheBoxGenerator
}

func (m *OutOfTheBox) Execute(targets map[string]pgs.File, _ map[string]pgs.Package) []pgs.Artifact {
	if len(targets) == 0 {
		return nil
	}

	lang := m.Parameters().Str(LangParam)
	m.Assert(lang != "", "`lang` parameter must be set")

	tpls := templates.Template(m.Parameters())[lang]
	m.Assert(tpls != nil, "could not find templates for `lang`: ", lang)

	for _, f := range sortedFiles(targets) {
		m.Push(f.Name().String())
		m.checkWebsocket(f)

		for _, tpl := range tpls {
			out := templates.FilePathFor(tpl)(f, m.ctx, tpl)
			// A nil path means no output should be generated for this file.
			if out != nil {
				m.AddGeneratorTemplateFile(out.String(), tpl, f)
			}
		}

		m.Pop()
	}

	return m.Artifacts()
}

func (m *OutOfTheBox) checkWebsocket(f pgs.File) {
	for _, svc := range f.Services() {
		for _, method := range svc.Methods() {
			var rule *api.WebsocketRule
			ok, err := method.Extension(api.E_Websocket, &rule)
			m.CheckErr(err, "unable to read websocket extension from method ", method.FullyQualifiedName())
			if !ok {
				continue
			}
			m.Assert(rule != nil && rule.GetUrl() != "",
				"websocket url is required for method ", method.FullyQualifiedName())
			m.Assert(method.ClientStreaming() || method.ServerStreaming(),
				"websocket option requires a streaming RPC: ", method.FullyQualifiedName())
		}
	}
}

func sortedFiles(targets map[string]pgs.File) []pgs.File {
	files := make([]pgs.File, 0, len(targets))
	for _, f := range targets {
		files = append(files, f)
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].InputPath().String() < files[j].InputPath().String()
	})
	return files
}
