package deployment

import (
	"context"
	"fmt"
	"github.com/go-keg/keg/cmd/deploy/cmd/config"
	"github.com/go-keg/keg/cmd/deploy/cmd/k8s/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"log"
	"time"
)

var Cmd = &cobra.Command{
	Use: "deployment",
}

var cfg config.Config
var name string

func init() {
	Cmd.AddCommand(updateImageCmd)
	Cmd.PersistentFlags().StringVarP(&name, "name", "n", "", "")
	confPath, _ := Cmd.Flags().GetString("conf")
	viper.SetConfigFile(confPath)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalf("fatal error config file: %v", err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
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
	patch, err := d.clientset.AppsV1().Deployments(d.Namespace).Patch(ctx, d.Name, patchType, data, v1.PatchOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("patch:", patch.Namespace, patch.Name)
	fmt.Println("data:", string(data))
}

func (d Deployments) UpdateImage(image string) {
	ops := pkg.Ops{
		{Op: "replace", Path: "/spec/template/spec/containers/0/image", Value: image},
		{Op: "replace", Path: "/spec/template/metadata/annotations/redeploy-timestamp", Value: fmt.Sprintf("%d", time.Now().Unix()*1000)},
	}
	d.Patch(types.JSONPatchType, ops.JSON())
}
