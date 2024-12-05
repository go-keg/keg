package deployment

import (
	"context"
	"log"

	"github.com/go-keg/keg/cmd/keg/cmd/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

var Cmd = &cobra.Command{
	Use: "deployment",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		confPath, _ := cmd.Flags().GetString("conf")
		viper.SetConfigFile(confPath)
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("fatal error config file: %v", err)
		}
		err = viper.Unmarshal(&cfg)
		if err != nil {
			panic(err)
		}
	},
}

var cfg config.Config
var name string

func init() {
	Cmd.AddCommand(updateImageCmd)
	Cmd.PersistentFlags().StringVarP(&name, "name", "n", "", "")
}

type Deployments struct {
	Namespace string
	Name      string
	clientset *kubernetes.Clientset
}

func NewDeployments(namespace string, name string, clientset *kubernetes.Clientset) *Deployments {
	return &Deployments{Namespace: namespace, Name: name, clientset: clientset}
}

func (d Deployments) Patch(patchType types.PatchType, data []byte) {
	ctx := context.Background()
	_, err := d.clientset.AppsV1().Deployments(d.Namespace).Patch(ctx, d.Name, patchType, data, v1.PatchOptions{})
	if err != nil {
		panic(err)
	}
}
