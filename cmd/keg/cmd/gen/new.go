package gen

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"github.com/go-keg/keg/cmd/keg/cmd/utils"
	"github.com/spf13/cobra"
)

type Config struct {
	Service    utils.Names
	DB         utils.Names
	UseGraphQL bool
	UseGRPC    bool
}

var (
	//go:embed template/*.tmpl
	_templates embed.FS
)

var Cmd = &cobra.Command{
	Use:     "new",
	Example: "codegen new account-example -d account",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalf("service-name is empty.")
		}
		if dbName == "" {
			dbName = args[0]
		}
		cfg := Config{
			Service:    utils.NewNames(args[0]),
			DB:         utils.NewNames(dbName),
			UseGRPC:    useGRPC,
			UseGraphQL: useGQL,
		}
		if _, err := os.Stat(fmt.Sprintf("internal/data/%s/ent", cfg.DB.KebabCase)); os.IsNotExist(err) {
			dataPath := fmt.Sprintf("internal/data/%s", cfg.DB.KebabCase)
			log.Fatalf("db-name invalid. err: %s\nsuggest: `$ mkdir -p %s && cd %s && ent new User && go generate ./... && cd ../../../`", err, dataPath, dataPath)
		}
		temp := template.Must(template.New("example").Funcs(template.FuncMap{
			"toUpper": strings.ToUpper,
		}).ParseFS(_templates, "template/*.tmpl"))
		// 执行主模板 base.tmpl 并传递数据
		files := map[string]string{
			fmt.Sprintf("cmd/%s/main.go", cfg.Service.KebabCase):                         "main.tmpl",
			fmt.Sprintf("cmd/%s/wire.go", cfg.Service.KebabCase):                         "wire.tmpl",
			fmt.Sprintf("configs/%s.yaml", cfg.Service.KebabCase):                        "config.tmpl",
			fmt.Sprintf("internal/app/%s/biz/biz.go", cfg.Service.KebabCase):             "biz.tmpl",
			fmt.Sprintf("internal/app/%s/cmd/migrate/migrate.go", cfg.Service.KebabCase): "cmd.migrate.tmpl",
			fmt.Sprintf("internal/app/%s/conf/conf.go", cfg.Service.KebabCase):           "conf.tmpl",
			fmt.Sprintf("internal/app/%s/data/data.go", cfg.Service.KebabCase):           "data.tmpl",
			fmt.Sprintf("internal/app/%s/job/job.go", cfg.Service.KebabCase):             "job.tmpl",
			fmt.Sprintf("internal/app/%s/schedule/schedule.go", cfg.Service.KebabCase):   "schedule.tmpl",
			fmt.Sprintf("internal/app/%s/server/http.go", cfg.Service.KebabCase):         "server.http.tmpl",
			fmt.Sprintf("internal/app/%s/server/server.go", cfg.Service.KebabCase):       "server.tmpl",
			fmt.Sprintf("internal/app/%s/service/service.go", cfg.Service.KebabCase):     "service.tmpl",
		}
		if cfg.UseGRPC {
			files[fmt.Sprintf("internal/app/%s/server/grpc.go", cfg.Service.KebabCase)] = "server.grpc.tmpl"
		}
		if cfg.UseGraphQL {
			files[fmt.Sprintf("internal/app/%s/service/graphql/generate.go", cfg.Service.KebabCase)] = "gql.generate.tmpl"
			files[fmt.Sprintf("internal/app/%s/service/graphql/gqlgen.yml", cfg.Service.KebabCase)] = "gql.gqlgen.tmpl"
			files[fmt.Sprintf("internal/app/%s/service/graphql/resolver.go", cfg.Service.KebabCase)] = "gql.resolver.tmpl"
			files[fmt.Sprintf("internal/data/%s/ent/entc.go", cfg.DB.KebabCase)] = "gql.entc.tmpl"
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
				log.Println("gen", _path)
				err := utils.WriteFileWithName(temp, cfg, _path, tempName)
				if err != nil {
					panic(err)
				}
			}
		}
		if useGQL {
			src := fmt.Sprintf("./internal/app/%s/service/graphql/", cfg.Service.KebabCase)
			fmt.Printf("exec: go generate %s\n", src)
			cmd := exec.Command("go", "generate", src)
			_, err := cmd.Output()
			if err != nil {
				log.Fatal("Error:", err)
			}
			cmd = exec.Command("go", "generate", fmt.Sprintf("./cmd/%s/", cfg.Service.KebabCase)) // #nosec G204
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
