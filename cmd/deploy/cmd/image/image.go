package image

import (
	"github.com/go-keg/keg/cmd/deploy/cmd/config"
	"github.com/go-keg/keg/cmd/deploy/cmd/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper"
	"log"
)

var Cmd = &cobra.Command{Use: "image"}

var cfg config.Config
var confPath, registry, name, tag string

func init() {
	Cmd.AddCommand(buildCmd, publishCmd, tagCmd)
	Cmd.PersistentFlags().StringVarP(&name, "name", "n", "", "")
	Cmd.PersistentFlags().StringVarP(&confPath, "config", "c", "./y-deploy.yaml", "")

	viper.SetConfigFile(confPath)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalf("fatal error config file: %v", err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	tag, err = utils.GetTag(cfg)
	if err != nil {
		log.Fatalf("fatal error config file: %v", err)
	}
}
