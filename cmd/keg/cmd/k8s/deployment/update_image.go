package deployment

import (
	"fmt"
	"log"

	"github.com/go-keg/keg/cmd/keg/cmd/k8s/pkg"
	"github.com/go-keg/keg/cmd/keg/cmd/utils"
	"github.com/spf13/cobra"
)

var updateImageCmd = &cobra.Command{
	Use:     "update-image",
	Example: "keg-deploy k8s deployment update-image -n account-interface -c ./y-deploy.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		branch, err := utils.GetBranch()
		if err != nil {
			log.Fatalf("get branch err: %v", err)
		}
		b := cfg.GetBranch(branch)
		deploy := NewDeployments(b.Namespace, name, pkg.NewClientSet(b.KubeConfig))
		tag, err := utils.GetTag(cfg)
		if err != nil {
			log.Fatalf("get tag err: %v", err)
		}
		image := fmt.Sprintf("%s/%s", cfg.ImageRegistry, name+tag)
		deploy.UpdateImage(image)
	},
}
