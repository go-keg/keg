package init

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-keg/keg/cmd/keg/config"
	"github.com/go-keg/keg/cmd/keg/internal"
	"github.com/go-keg/keg/cmd/keg/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	yaml3 "gopkg.in/yaml.v3"
)

var Cmd = &cobra.Command{
	Use:     "init",
	Example: "keg init",
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg config.Config
		var err error
		currentDir, err := os.Getwd()
		if err != nil {
			return err
		}
		configPath := filepath.Join(currentDir, internal.ConfigFile)
		f, err := os.Stat(configPath)
		fmt.Println("configPath", configPath, f, err)

		if os.IsNotExist(err) {
			cfg.GoModule, err = utils.GoModuleName(currentDir)
			if err != nil {
				return err
			}
			if cfg.GoModule == "" {
				pn := promptui.Prompt{
					Label:   "GoModule",
					Default: "",
				}
				cfg.GoModule, err = pn.Run()
				if err != nil {
					return err
				}
				pn = promptui.Prompt{
					Label:   "DefaultService",
					Default: "account",
				}
			}

			bytes, err := yaml3.Marshal(&cfg)
			if err != nil {
				return err
			}
			outFile, err := os.Create(configPath)
			if err != nil {
				panic(err)
			}
			defer func() { _ = outFile.Close() }()
			_, err = outFile.Write(bytes)
			if err != nil {
				return err
			}
		}
		return err
	},
}
