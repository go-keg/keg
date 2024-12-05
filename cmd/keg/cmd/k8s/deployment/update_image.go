package deployment

import (
	"fmt"
	"log"
	"time"

	"github.com/go-keg/keg/cmd/keg/cmd/k8s/pkg"
	"github.com/go-keg/keg/cmd/keg/cmd/utils"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/types"
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

		ts := time.Now().Unix() * 1000
		deploy.Patch(types.StrategicMergePatchType, []byte(fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"redeploy-timestamp":"%d"}},"spec":{"containers":[{"name":"%s","image":"%s"}]}}}}`, ts, deploy.Name, image)))
		fmt.Printf(`updateImage: {"namespace":"%s", "redeploy.timestamp":%d, "name":"%s", "image":"%s"}`, deploy.Namespace, ts, deploy.Name, image)
	},
}
