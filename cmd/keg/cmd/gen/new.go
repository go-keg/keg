package gen

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"text/template"

	"github.com/go-keg/keg/cmd/keg/config"
	"github.com/go-keg/keg/cmd/keg/utils"
	"github.com/spf13/cobra"
)

var (
	//go:embed templates/**/*.tmpl
	templates embed.FS
)

var Cmd = &cobra.Command{
	Use:     "new",
	Example: "keg new auth-service -d account",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalf("service-name is empty.")
		}
		if dbName == "" {
			dbName = args[0]
		}
		cfg := config.App{
			Name:       config.Name(args[0]),
			DB:         config.Name(dbName),
			Schedule:   false,
			Job:        false,
			UseGRPC:    useGRPC,
			UseGraphQL: useGQL,
		}
		if _, err := os.Stat(fmt.Sprintf("internal/data/%s/ent", cfg.DB.KebabCase())); os.IsNotExist(err) {
			dataPath := fmt.Sprintf("internal/data/%s", cfg.DB.KebabCase())
			log.Fatalf("db-name invalid. err: %s\nsuggest: `$ mkdir -p %s && cd %s && ent new User && go generate ./... && cd ../../../`", err, dataPath, dataPath)
		}
		temp := template.Must(template.New("keg").ParseFS(templates, "templates/**/*.tmpl"))
		// 执行主模板 base.tmpl 并传递数据
		files := map[string]string{
			fmt.Sprintf("cmd/%s/main.go", cfg.Name.KebabCase()):  "main.tmpl",
			fmt.Sprintf("cmd/%s/wire.go", cfg.Name.KebabCase()):  "wire.tmpl",
			fmt.Sprintf("configs/%s.yaml", cfg.Name.KebabCase()): "config.tmpl",
			// app
			fmt.Sprintf("internal/app/%s/biz/biz.go", cfg.Name.KebabCase()):             "biz.tmpl",
			fmt.Sprintf("internal/app/%s/cmd/migrate/migrate.go", cfg.Name.KebabCase()): "cmd.migrate.tmpl",
			fmt.Sprintf("internal/app/%s/conf/conf.go", cfg.Name.KebabCase()):           "conf.tmpl",
			fmt.Sprintf("internal/app/%s/data/data.go", cfg.Name.KebabCase()):           "data.tmpl",
			fmt.Sprintf("internal/app/%s/job/job.go", cfg.Name.KebabCase()):             "job.tmpl",
			fmt.Sprintf("internal/app/%s/schedule/schedule.go", cfg.Name.KebabCase()):   "schedule.tmpl",
			fmt.Sprintf("internal/app/%s/server/http.go", cfg.Name.KebabCase()):         "server.http.tmpl",
			fmt.Sprintf("internal/app/%s/server/server.go", cfg.Name.KebabCase()):       "server.tmpl",
			fmt.Sprintf("internal/app/%s/service/service.go", cfg.Name.KebabCase()):     "service.tmpl",
			// deploy
			fmt.Sprintf("deploy/build/%s/Dockerfile", cfg.Name.KebabCase()): "build.tmpl",
			fmt.Sprintf("deploy/kubernetes/%s.yaml", cfg.Name.KebabCase()):  "kubernetes.tmpl",
		}
		if cfg.UseGRPC {
			files[fmt.Sprintf("internal/app/%s/server/grpc.go", cfg.Name.KebabCase())] = "server.grpc.tmpl"
		}
		if cfg.UseGraphQL {
			files[fmt.Sprintf("internal/app/%s/service/graphql/generate.go", cfg.Name.KebabCase())] = "gql.generate.tmpl"
			files[fmt.Sprintf("internal/app/%s/service/graphql/gqlgen.yml", cfg.Name.KebabCase())] = "gql.gqlgen.tmpl"
			files[fmt.Sprintf("internal/app/%s/service/graphql/graphql.config.yml", cfg.Name.KebabCase())] = "gql.config.tmpl"
			files[fmt.Sprintf("internal/app/%s/service/graphql/resolver.go", cfg.Name.KebabCase())] = "gql.resolver.tmpl"
			files[fmt.Sprintf("internal/data/%s/ent/entc.go", cfg.DB.KebabCase())] = "entc.tmpl"
			files[fmt.Sprintf("internal/app/%s/service/graphql/dataloader/dataloader.go", cfg.Name.KebabCase())] = "gql.dataloader.tmpl"
			files[fmt.Sprintf("internal/app/%s/service/graphql/dataloader/loader.go", cfg.Name.KebabCase())] = "gql.loader.tmpl"
			files[fmt.Sprintf("internal/app/%s/service/graphql/model/extend.go", cfg.Name.KebabCase())] = "gql.model_extend.tmpl"
		}
		for _path, tempName := range files {
			if _, err := os.Stat(_path); os.IsNotExist(err) {
				dir := path.Dir(_path)
				// 获取文件或目录的信息
				if _, err := os.Stat(dir); os.IsNotExist(err) {
					err := os.MkdirAll(dir, 0755)
					if err != nil {
						fmt.Println("无法创建目录:", err)
						return
					}
				} else if err != nil {
					panic(err)
				}
				log.Println("gen", _path, tempName)
				err := utils.WriteFileWithName(temp, cfg, _path, tempName)
				if err != nil {
					panic(err)
				}
			}
		}
		if useGQL {
			src := fmt.Sprintf("./internal/app/%s/service/graphql/", cfg.Name)
			fmt.Printf("exec: go generate %s\n", src)
			cmd := exec.Command("go", "generate", src)
			_, err := cmd.Output()
			if err != nil {
				log.Fatal("Error:", err)
			}
			cmd = exec.Command("go", "generate", fmt.Sprintf("./cmd/%s/", cfg.Name.KebabCase())) // #nosec G204
			_, err = cmd.Output()
			if err != nil {
				log.Fatal("Error:", err)
			}
		}
	},
}

var dbName string
var useGRPC, useGQL bool

func init() {
	Cmd.AddCommand(bizCmd)
	Cmd.Flags().StringVarP(&dbName, "db-name", "d", "", "database name eg: -n account (default: ${ServiceName})")
	Cmd.Flags().BoolVar(&useGRPC, "grpc", false, "use GRPC")
	Cmd.Flags().BoolVar(&useGQL, "gql", false, "use Graphql")
}
