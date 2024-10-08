package deployment

import (
	"fmt"
	"github.com/go-keg/keg/cmd/deploy/cmd/k8s/pkg"
	"github.com/go-keg/keg/cmd/deploy/cmd/utils"
	"github.com/spf13/cobra"
	"log"
)

var updateImageCmd = &cobra.Command{
	Use:     "update-image",
	Example: "y-deploy k8s deployment update-image -n account-interface -c ./y-deploy.yaml",
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
		image := fmt.Sprintf("%s/%s", cfg.Registry, name+tag)
		deploy.UpdateImage(image)
	},
}
