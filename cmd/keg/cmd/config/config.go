package config

import (
	"path"
	"strings"

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
}

type Branch struct {
	Branch     string    `yaml:"branch"`
	Namespace  string    `yaml:"namespace"`
	KubeConfig string    `yaml:"kubeConfig"`
	TagPolicy  TagPolicy `yaml:"tagPolicy"`
}

type App struct {
	Name       string `yaml:"name"`
	DB         string `yaml:"db"`
	Job        bool   `yaml:"job"`
	Schedule   bool   `yaml:"schedule"`
	UseGraphQL bool   `yaml:"useGraphQL"`
	UseGRPC    bool   `yaml:"useGRPC"`
}

func (c Config) GetBranch(branch string) Branch {
	item, ok := lo.Find(c.Branches, func(item Branch) bool {
		return item.Branch == branch
	})
	if !ok {
		item = c.Default
	}
	if strings.HasPrefix(item.KubeConfig, "~/") {
		item.KubeConfig = path.Join(homedir.HomeDir(), item.KubeConfig[2:])
	}
	return item
}
