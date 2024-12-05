package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-keg/keg/cmd/keg/internal"
	"github.com/go-keg/keg/cmd/keg/utils"
	"github.com/samber/lo"
	"k8s.io/client-go/util/homedir"
)

// TagPolicy 镜像 tag 生成策略
type TagPolicy string

const (
	TagPolicyVersion       TagPolicy = "version"        // xxx-service:v1.0.0
	TagPolicyBranch        TagPolicy = "branch"         // xxx-service:latest-qa
	TagPolicyVersionBranch TagPolicy = "version-branch" // xxx-service:v1.0.0-6-gdaeffa2-qa
)

type Config struct {
	GoModule      string   `yaml:"go_module"`
	ImageRegistry string   `yaml:"image_registry"`
	Branches      []Branch `yaml:"branches"`
	Default       Branch   `yaml:"default"`
	Apps          []App    `yaml:"apps"`
	K8s           K8s      `yaml:"k8s"`
}

func (r Config) App() App {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	rootPath, err := utils.ProjectRootPath()
	if err != nil {
		log.Fatal(err)
	}

	// 检查文件是否存在
	_, err = os.Stat(filepath.Join(rootPath, internal.ConfigFile))
	if err == nil {
		if !strings.Contains(currentDir, filepath.Join(rootPath, "internal/app/")) {
			log.Fatal("当前所在目录不在 internal/app/{service-name} 中")
		}
		rootPath := strings.TrimPrefix(currentDir, filepath.Join(rootPath, "internal/app/"))
		fmt.Println("rootPath", rootPath)
		arr := strings.SplitN(rootPath, "/", 3)
		if len(arr) >= 2 {
			serviceName := arr[1]
			fmt.Println("serviceName:", serviceName)
			for _, app := range r.Apps {
				if string(app.Name) == serviceName {
					return app
				}
			}
			log.Fatalf("当前服务不在 %s 中\n", internal.ConfigFile)
		}
	}
	return App{}
}

type Branch struct {
	Branch     string    `yaml:"branch"`
	Namespace  string    `yaml:"namespace"`
	KubeConfig string    `yaml:"kubeConfig"`
	TagPolicy  TagPolicy `yaml:"tagPolicy"`
}

type App struct {
	Name       Name `yaml:"name"`
	DB         Name `yaml:"db"`
	Job        bool `yaml:"job"`
	Schedule   bool `yaml:"schedule"`
	UseGraphQL bool `yaml:"use_graphql"`
	UseGRPC    bool `yaml:"use_grpc"`
}

func (r App) GoModule() string {
	rootPath, err := utils.ProjectRootPath()
	if err != nil {
		panic(err)
	}
	moduleName, err := utils.GoModuleName(rootPath)
	if err != nil {
		panic(err)
	}
	return moduleName
}

type Name string

func (r Name) PascalCase() string {
	return lo.PascalCase(string(r))
}

func (r Name) CamelCase() string {
	return lo.CamelCase(string(r))
}

func (r Name) KebabCase() string {
	return lo.KebabCase(string(r))
}

func (r Name) SnakeCase() string {
	return lo.SnakeCase(string(r))
}

type K8s struct {
	ImageRegistry     string `yaml:"image_registry"`
	Namespace         string `yaml:"namespace"`
	ImageVersion      string `yaml:"image_version"`
	DefaultReplicas   int    `yaml:"default_replicas"`
	ImagePullPolicy   string `yaml:"image_pull_policy"`
	ResourcesRequests struct {
		Cpu    string `yaml:"cpu"`
		Memory string `yaml:"memory"`
	} `yaml:"resources_requests"`
}

func (r Config) GetBranch() Branch {
	branch, err := utils.GetBranch()
	if err != nil {
		log.Fatalf("get branch err: %v", err)
	}
	item, ok := lo.Find(r.Branches, func(item Branch) bool {
		return item.Branch == branch
	})
	if !ok {
		item = r.Default
	}
	if strings.HasPrefix(item.KubeConfig, "~/") {
		item.KubeConfig = path.Join(homedir.HomeDir(), item.KubeConfig[2:])
	}
	return item
}

func (r Config) GetTag() (string, error) {
	version, err := utils.GetVersion()
	if err != nil {
		return "", err
	}
	branch, err := utils.GetBranch()
	if err != nil {
		return "", err
	}
	tagTypes := map[TagPolicy]string{
		TagPolicyVersion:       version,
		TagPolicyBranch:        fmt.Sprintf("latest-%s", branch),
		TagPolicyVersionBranch: fmt.Sprintf("%s-%s", version, branch),
	}
	b := r.GetBranch()
	return tagTypes[b.TagPolicy], nil
}
