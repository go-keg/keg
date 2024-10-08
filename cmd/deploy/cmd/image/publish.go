package image

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os/exec"
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("%s/%s", cfg.Registry, name+tag)
		fmt.Println("exec:", "docker", "tag", name+tag, url)
		err := exec.Command("docker", "tag", name+tag, url).Run()
		if err != nil {
			log.Fatalf("exec docker tag error: %v", err)
		}
		fmt.Println("exec:", "docker", "push", url)
		c := exec.Command("docker", "push", url)
		stdout, err := c.StdoutPipe()
		if err != nil {
			panic(err)
		}

		err = c.Start()
		if err != nil {
			log.Fatalf("exec docker push error: %v", err)
		}

		reader := bufio.NewReader(stdout)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				break
			}
			fmt.Print(readString)

		}
		err = c.Wait()
	},
}
