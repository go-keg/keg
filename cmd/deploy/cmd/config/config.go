package config

import (
	"k8s.io/client-go/util/homedir"
	"path"
	"strings"
)

// TagPolicy 镜像 tag 生成策略
type TagPolicy string

const (
	TagPolicyVersion       TagPolicy = "version"        // xxx-service:v1.0.0
	TagPolicyBranch        TagPolicy = "branch"         // xxx-service:latest-qa
	TagPolicyVersionBranch TagPolicy = "version-branch" // xxx-service:v1.0.0-6-gdaeffa2-qa
)

type Config struct {
	Registry string    `yaml:"registry"`
	Branches []Branch  `yaml:"branches"`
	Services []Service `yaml:"services"`
	Default  Branch    `yaml:"default"`
}

type Service struct {
	Name    string
	Changes []string
	Scripts []string
}

type Branch struct {
	Branch     string
	Namespace  string
	KubeConfig string
	TagPolicy  TagPolicy
}

func (c Config) GetBranch(branch string) Branch {
	b := c.Default
	for _, item := range c.Branches {
		if item.Branch == branch {
			b = item
		}
	}
	if strings.HasPrefix(b.KubeConfig, "~/") {
		b.KubeConfig = path.Join(homedir.HomeDir(), b.KubeConfig[2:])
	}
	return b
}
