package plugins

import (
	_ "embed"
	"strings"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
	"github.com/go-openapi/inflect"
	"github.com/samber/lo"
)

//go:embed dataloader.tmpl
var dataloaderTemplate string

func NewDataloader(filename, entPath string, models []string) plugin.Plugin {
	return &Dataloader{
		filename: filename,
		models:   models,
		entPath:  entPath,
	}
}

type Dataloader struct {
	filename string
	models   []string
	entPath  string
}

var _ plugin.CodeGenerator = &Dataloader{}

func (m *Dataloader) Name() string {
	return "dataloader"
}

func (m *Dataloader) GenerateCode(data *codegen.Data) error {
	inf := inflect.NewDefaultRuleset()
	build := &DataloaderBuild{
		Models: lo.Map(m.models, func(item string, index int) ModelName {
			return ModelName{
				ToLower:    strings.ToLower(item),
				PascalCase: lo.PascalCase(item),
				CamelCase:  lo.CamelCase(item),
				Pluralize:  lo.PascalCase(inf.Pluralize(item)),
			}
		}),
		EntPath: m.entPath,
	}
	packageName := data.Config.Resolver.Package
	if s := strings.Split(m.filename, "/"); len(s) > 1 {
		packageName = s[len(s)-2]
	}
	return templates.Render(templates.Options{
		PackageName:     packageName,
		Filename:        m.filename,
		Data:            build,
		Packages:        data.Config.Packages,
		Template:        dataloaderTemplate,
		GeneratedHeader: true,
	})
}

type ModelName struct {
	ToLower    string
	PascalCase string
	CamelCase  string
	Pluralize  string
}
type DataloaderBuild struct {
	codegen.Data

	ExecPackageName     string
	ResolverPackageName string
	Models              []ModelName
	EntPath             string
}
