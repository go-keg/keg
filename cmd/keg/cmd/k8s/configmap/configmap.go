package configmap

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/go-keg/keg/cmd/keg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var name, filePath string
var cfg config.Config

var Cmd = &cobra.Command{
	Use: "configmap",
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

func init() {
	Cmd.AddCommand(applyCmd)
	Cmd.PersistentFlags().StringVarP(&name, "name", "n", "", "config name")
}

func ApplyConfigMap(ctx context.Context, client *kubernetes.Clientset, namespace, name, filePath string) error {
	_, err := client.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return CreateConfigMap(ctx, client, namespace, name, filePath)
	}
	return UpdateConfigMap(ctx, client, namespace, name, filePath)
}

func CreateConfigMap(ctx context.Context, client *kubernetes.Clientset, namespace, name, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: map[string]string{
			filepath.Base(filePath): string(data),
		},
	}

	_, err = client.CoreV1().ConfigMaps(namespace).Create(ctx, cm, metav1.CreateOptions{})
	return err
}

func UpdateConfigMap(ctx context.Context, client *kubernetes.Clientset, namespace, name, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	cm, err := client.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if cm.Data == nil {
		cm.Data = map[string]string{}
	}
	cm.Data[filepath.Base(filePath)] = string(data)

	_, err = client.CoreV1().ConfigMaps(namespace).Update(ctx, cm, metav1.UpdateOptions{})
	return err
}
