package generator

import (
	"sort"

	"github.com/appootb/substratum/v2/tools/cmd/protoc-gen-validate/templates"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

const (
	ValidateGenerator = "validate"
	LangParam         = "lang"
)

type Validate struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
}

func New() pgs.Module {
	return &Validate{
		ModuleBase: &pgs.ModuleBase{},
	}
}

func (m *Validate) InitContext(ctx pgs.BuildContext) {
	m.ModuleBase.InitContext(ctx)
	m.ctx = pgsgo.InitContext(ctx.Parameters())
}

func (m *Validate) Name() string {
	return ValidateGenerator
}

func (m *Validate) Execute(targets map[string]pgs.File, _ map[string]pgs.Package) []pgs.Artifact {
	lang := m.Parameters().Str(LangParam)
	m.Assert(lang != "", "`lang` parameter must be set")

	tpls := templates.Template(m.Parameters())[lang]
	m.Assert(tpls != nil, "could not find templates for `lang`: ", lang)

	for _, f := range sortedFiles(targets) {
		m.Push(f.Name().String())
		for _, msg := range f.AllMessages() {
			m.CheckRules(msg)
		}
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
