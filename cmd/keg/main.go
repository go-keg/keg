package main

import (
	"fmt"

	"github.com/go-keg/keg/cmd/keg/cmd/gen"
	"github.com/go-keg/keg/cmd/keg/cmd/image"
	initCmd "github.com/go-keg/keg/cmd/keg/cmd/init"
	"github.com/go-keg/keg/cmd/keg/cmd/k8s"
	"github.com/go-keg/keg/cmd/keg/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "keg",
}

var version = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(internal.Version)
	},
}

func init() {
	rootCmd.AddCommand(image.Cmd, k8s.Cmd, gen.Cmd, initCmd.Cmd, version)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
