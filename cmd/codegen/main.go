package main

import (
	"github.com/go-keg/keg/cmd/codegen/new"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "codegen",
	Short: "Monorepo CodeGen",
}

func init() {
	rootCmd.AddCommand(new.Cmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
