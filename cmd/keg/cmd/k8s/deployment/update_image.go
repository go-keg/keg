package deployment

import (
	"fmt"
	"time"

	"github.com/go-keg/keg/cmd/keg/cmd/k8s/pkg"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/types"
)

var tag string

func init() {
	updateImageCmd.Flags().StringVarP(&tag, "tag", "t", "latest", "Tag of the deployment to update")
}

var updateImageCmd = &cobra.Command{
	Use:     "update-image",
	Example: "keg k8s deployment update-image -n account-interface",
	Run: func(cmd *cobra.Command, args []string) {
		b := cfg.GetBranch()
		deploy := NewDeployments(b.Namespace, name, pkg.NewClientSet(b.KubeConfig))
		image := fmt.Sprintf("%s/%s:%s", cfg.ImageRegistry, name, tag)
		ts := time.Now().Unix() * 1000
		deploy.Patch(types.StrategicMergePatchType, []byte(fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"redeploy-timestamp":"%d"}},"spec":{"containers":[{"name":"%s","image":"%s"}]}}}}`, ts, deploy.Name, image)))
		fmt.Printf("updateImage namespace: %s redeploy.timestamp: %d name: %s image: %s\n", deploy.Namespace, ts, deploy.Name, image)
	},
}
