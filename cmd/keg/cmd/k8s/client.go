package k8s

import (
	"github.com/go-keg/keg/cmd/keg/cmd/k8s/configmap"
	"github.com/go-keg/keg/cmd/keg/cmd/k8s/deployment"
	"github.com/go-keg/keg/cmd/keg/cmd/k8s/gen"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "k8s",
}

func init() {
	Cmd.AddCommand(deployment.Cmd, gen.Cmd, configmap.Cmd)
}
