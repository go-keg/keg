package gen

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var namespace, dir, outDir string

func init() {
	configCmd.Flags().StringVarP(&namespace, "namespace", "n", "local", "")
	configCmd.Flags().StringVarP(&dir, "dir", "d", "./deploy/kubernetes", "")
	configCmd.Flags().StringVarP(&outDir, "outDir", "o", "./deploy/kubernetes/output/", "")
}

var configCmd = &cobra.Command{
	Use:     "config",
	Example: "keg k8s gen config -n dev",
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := os.OpenFile(".env.k8s", os.O_RDONLY, 0666)
		envs, err := godotenv.Parse(file)
		if err != nil {
			panic(err)
		}
		service := &Temp{
			Path:   "./deploy/kubernetes",
			Envs:   envs,
			Output: outDir,
		}

		service.Read(service.Path)
		log.Println("configure files generated.")
	},
}

type Temp struct {
	Path   string
	Output string
	Envs   map[string]string
}

func (c Temp) Read(dir string) {
	var items = make([]string, 0, len(c.Envs)+2)
	items = append(items, "${NAMESPACE}", namespace)
	for key, val := range c.Envs {
		items = append(items, "${"+key+"}", val)
	}
	replacer := strings.NewReplacer(items...)

	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, fileInfo := range files {
		if !fileInfo.IsDir() && strings.HasSuffix(fileInfo.Name(), ".yaml") {
			content, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, fileInfo.Name()))
			if err != nil {
				panic(err)
			}
			content = []byte(replacer.Replace(string(content)))
			c.Write(fileInfo.Name(), content)
		}
	}
}

func (c Temp) Write(filename string, content []byte) {
	filePath := path.Join(c.Output, namespace+"_"+filename)
	log.Println("generate", filePath)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	_, err = f.Write(content)
	if err != nil {
		panic(err)
	}
}
