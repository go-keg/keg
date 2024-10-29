package main

import (
	"github.com/go-keg/keg/cmd/keg/cmd/gen"
	"github.com/go-keg/keg/cmd/keg/cmd/image"
	"github.com/go-keg/keg/cmd/keg/cmd/k8s"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "keg",
}

func init() {
	rootCmd.AddCommand(image.Cmd, k8s.Cmd, gen.Cmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
