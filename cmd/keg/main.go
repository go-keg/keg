package main

import (
	"fmt"

	"github.com/go-keg/keg/cmd/keg/cmd/gen"
	"github.com/go-keg/keg/cmd/keg/cmd/image"
	"github.com/go-keg/keg/cmd/keg/cmd/k8s"
	"github.com/go-keg/keg/cmd/keg/config"
	"github.com/spf13/cobra"
)

const Version = "v0.1.2"

var rootCmd = &cobra.Command{
	Use: "keg",
}

var version = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	rootCmd.AddCommand(image.Cmd, k8s.Cmd, gen.Cmd, version)
	rootCmd.PersistentFlags().StringP("conf", "c", config.FileName, "keg config file")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
