package main

import (
	"github.com/go-keg/keg/cmd/deploy/cmd/image"
	"github.com/go-keg/keg/cmd/deploy/cmd/k8s"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "deploy",
}

func init() {
	rootCmd.AddCommand(image.Cmd, k8s.Cmd)
	rootCmd.PersistentFlags().String("conf", "./.deploy.yaml", "")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
