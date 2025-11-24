package configmap

import (
	"github.com/go-keg/keg/cmd/keg/cmd/k8s/pkg"
	"github.com/spf13/cobra"
)

func init() {
	applyCmd.Flags().StringVarP(&filePath, "filePath", "f", "", "config file path")
}

var applyCmd = &cobra.Command{
	Use: "apply",
	RunE: func(cmd *cobra.Command, args []string) error {
		b := cfg.GetBranch()
		return ApplyConfigMap(cmd.Context(), pkg.NewClientSet(b.KubeConfig), b.Namespace, name, filePath)
	},
}
