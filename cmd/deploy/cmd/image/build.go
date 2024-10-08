package image

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("exec:", "docker", "build", "-t", name+tag, "-f", dockerfile, ".")
		err := exec.Command("docker", "build", "-t", name+tag, "-f", dockerfile, ".").Run()
		if err != nil {
			log.Fatalf("exec docker build error: %v", err)
		}
	},
}

var dockerfile string

func init() {
	buildCmd.Flags().StringVarP(&dockerfile, "dockerfile", "f", "", "")
}
