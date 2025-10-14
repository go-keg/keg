package gen

import (
	"fmt"
	"log"
	"text/template"

	"github.com/go-keg/keg/cmd/keg/config"
	"github.com/go-keg/keg/cmd/keg/utils"
	"github.com/spf13/cobra"
)

var skipRepo bool

func init() {
	bizCmd.Flags().BoolVarP(&skipRepo, "skip-repo", "s", false, "skip generate repo")
}

const (
	bizTemp = `package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type {{.PascalCase}}Repo interface {
	Todo(ctx context.Context) (string, error)
}

type {{.PascalCase}}UseCase struct {
	log  *log.Helper
	repo {{.PascalCase}}Repo
}

func New{{.PascalCase}}UseCase(logger log.Logger, repo {{.PascalCase}}Repo) *{{.PascalCase}}UseCase {
	return &{{.PascalCase}}UseCase{
		log:  log.NewHelper(log.With(logger, "module", "usecase/{{.SnakeCase}}")),
		repo: repo,
	}
}

func (r {{.PascalCase}}UseCase) Todo(ctx context.Context) (string, error) {
	return r.repo.Todo(ctx)
}
`
	dataTemp = `package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"{{.GoModule}}/internal/app/{{.ServiceName}}/biz"
)

type {{.CamelCase}}Repo struct {
	log  *log.Helper
}

func New{{.PascalCase}}Repo(logger log.Logger) biz.{{.PascalCase}}Repo {
	return &{{.CamelCase}}Repo{
		log:  log.NewHelper(log.With(logger, "module", "data/{{.SnakeCase}}")),
	}
}

func (r {{.CamelCase}}Repo) Todo(ctx context.Context) (string, error) {
	return "todo...", nil
}
`
)

var bizCmd = &cobra.Command{
	Use:     "biz",
	Example: "keg new biz account_relation",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalf("biz name is empty.")
		}
		names := config.Name(args[0])
		dir, err := utils.ExecDir()
		if err != nil {
			log.Fatal(err)
		}
		biz := template.Must(template.New("biz").Parse(bizTemp))
		data := template.Must(template.New("data").Parse(dataTemp))

		err = utils.WriteFile(biz, map[string]any{
			"PascalCase": names.PascalCase(),
			"SnakeCase":  names.SnakeCase(),
		}, fmt.Sprintf("biz/%s.go", names.SnakeCase()))
		if err != nil {
			log.Fatal(err)
		}
		if !skipRepo {
			rootPath, err := utils.ProjectRootPath()
			if err != nil {
				log.Fatal(err)
			}
			moduleName, err := utils.GoModuleName(rootPath)
			if err != nil {
				log.Fatal(err)
			}
			err = utils.WriteFile(data, map[string]any{
				"GoModule":    moduleName,
				"ServiceName": dir,
				"CamelCase":   names.CamelCase(),
				"PascalCase":  names.PascalCase(),
				"SnakeCase":   names.SnakeCase(),
			}, fmt.Sprintf("data/%s.go", names.SnakeCase()))
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}
