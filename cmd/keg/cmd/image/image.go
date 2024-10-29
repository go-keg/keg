package image

import (
	"log"

	"github.com/go-keg/keg/cmd/keg/cmd/config"
	"github.com/go-keg/keg/cmd/keg/cmd/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper"
)

var Cmd = &cobra.Command{Use: "image"}

var cfg config.Config
var confPath, name, tag string

func init() {
	Cmd.AddCommand(tagCmd)
	Cmd.PersistentFlags().StringVarP(&name, "name", "n", "", "")
	Cmd.PersistentFlags().StringVarP(&confPath, "config", "c", "./.deploy.yaml", "")

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
