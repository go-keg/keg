package config

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/go-keg/keg/cmd/keg/utils"
	"github.com/samber/lo"
	"k8s.io/client-go/util/homedir"
)

const FileName = "keg.yaml"

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

func (r Name) UpperCase() string {
	return strings.ToUpper(lo.SnakeCase(string(r)))
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
	switch r.GetBranch().TagPolicy {
	case TagPolicyVersion:
		return version, nil
	case TagPolicyBranch:
		return fmt.Sprintf("latest-%s", branch), nil
	case TagPolicyVersionBranch:
		return fmt.Sprintf("%s-%s", version, branch), nil
	default:
		return fmt.Sprintf("%s-%s", version, branch), nil
	}
}
