package image

import (
	"log"

	"github.com/go-keg/keg/cmd/keg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &cobra.Command{Use: "image", PersistentPreRun: func(cmd *cobra.Command, args []string) {
	confPath, _ := cmd.Flags().GetString("conf")
	viper.SetConfigFile(confPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("fatal error config file: %v", err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	tag, err = cfg.GetTag()
	if err != nil {
		log.Fatalf("fatal error config file: %v", err)
	}
}}

var cfg config.Config
var tag string

func init() {
	Cmd.AddCommand(tagCmd)
}
