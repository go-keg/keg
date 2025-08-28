package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

func Load[T any](path string, envs ...string) (*T, error) {
	envs = lo.FilterMap(envs, func(item string, index int) (string, bool) {
		if item == "" {
			return "", false
		}
		if !strings.HasPrefix(item, ".env.") {
			item = ".env." + item
		}
		return item, true
	})
	if len(envs) == 0 {
		envs = append(envs, ".env")
	}
	_ = godotenv.Load(envs...)
	readFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(strings.NewReader(replaceEnvVariables(readFile)))
	if err != nil {
		return nil, err
	}
	var cfg T
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoadEnv(filenames ...string) {
	if len(filenames) == 0 {
		filenames = append(filenames, ".env")
	}
	if err := godotenv.Load(filenames...); err != nil {
		fmt.Printf("loading env file error: %s\r\n", err.Error())
	}
}

var reEnv = regexp.MustCompile(`\$\{[\w_-]+}`)

func replaceEnvVariables(text []byte) string {
	var words = make([]string, 0, len(os.Environ())*2)
	for _, env := range os.Environ() {
		envPair := strings.SplitN(env, "=", 2)
		words = append(words, fmt.Sprintf("${%s}", envPair[0]), strings.Trim(strconv.Quote(envPair[1]), "\""))
	}
	str := strings.NewReplacer(words...).Replace(string(text))
	return reEnv.ReplaceAllString(str, "")
}
