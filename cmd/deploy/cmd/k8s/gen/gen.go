package gen

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "gen",
}

func init() {
	Cmd.AddCommand(configCmd)
}
