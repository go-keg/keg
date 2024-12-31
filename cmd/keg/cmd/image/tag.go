package image

import (
	"fmt"

	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use: "tag",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(tag)
	},
}
