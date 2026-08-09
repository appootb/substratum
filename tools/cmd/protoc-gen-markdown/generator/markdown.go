package generator

import (
	"sort"

	"github.com/appootb/substratum/v2/tools/cmd/protoc-gen-markdown/templates"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

const (
	MarkdownGenerator = "markdown"
)

type Markdown struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
}

func New() pgs.Module {
	return &Markdown{
		ModuleBase: &pgs.ModuleBase{},
	}
}

func (m *Markdown) InitContext(ctx pgs.BuildContext) {
	m.ModuleBase.InitContext(ctx)
	m.ctx = pgsgo.InitContext(ctx.Parameters())
}

func (m *Markdown) Name() string {
	return MarkdownGenerator
}

func (m *Markdown) Execute(targets map[string]pgs.File, _ map[string]pgs.Package) []pgs.Artifact {
	if len(targets) == 0 {
		return nil
	}

	files := sortedFiles(targets)
	tpls := templates.ProtoTemplate(m.Parameters())

	var outDir pgs.FilePath
	var haveOut bool
	for _, f := range files {
		m.Push(f.Name().String())
		for _, tpl := range tpls {
			out := templates.FilePathFor(tpl)(f, m.ctx, tpl)
			if out != nil {
				outDir = out.Dir()
				haveOut = true
				m.AddGeneratorTemplateFile(out.String(), tpl, f)
			}
		}
		m.Pop()
	}

	if !haveOut {
		return m.Artifacts()
	}

	tocTpl := templates.ReadMeTemplate(m.Parameters())
	tocOut := pgs.JoinPaths(outDir.String(), "README.md")
	m.Push("readme")
	m.AddGeneratorTemplateFile(tocOut.String(), tocTpl, files)
	m.Pop()

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
