package utils

import (
	"fmt"
	"github.com/go-keg/keg/cmd/deploy/cmd/config"
	"io"
	"os/exec"
	"strings"
)

func GetBranch() (string, error) {
	cmd1 := exec.Command("git", "branch")
	cmd2 := exec.Command("sed", "-n", "/\\* /s///p")
	in, _ := cmd2.StdinPipe()
	out, _ := cmd1.StdoutPipe()
	go func() {
		defer func() {
			_ = out.Close()
			_ = in.Close()
		}()
		_, _ = io.Copy(in, out)
	}()
	err := cmd1.Run()
	if err != nil {
		return "", err
	}
	output, err := cmd2.Output()
	if err != nil {
		return "", err
	}
	branch := string(output)
	branch = strings.ReplaceAll(branch, "/", "-")
	branch = strings.ReplaceAll(branch, "#", "")
	return strings.Trim(branch, "\n"), nil
}

func GetVersion() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--always")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.Trim(strings.ReplaceAll(string(output), "\n", ""), "refs/tags/"), nil
}

func GetTag(cfg config.Config) (string, error) {
	version, err := GetVersion()
	if err != nil {
		return "", err
	}
	branch, err := GetBranch()
	if err != nil {
		return "", err
	}
	tagTypes := map[config.TagPolicy]string{
		config.TagPolicyVersion:       fmt.Sprintf(":%s", version),
		config.TagPolicyBranch:        fmt.Sprintf(":latest-%s", branch),
		config.TagPolicyVersionBranch: fmt.Sprintf(":%s-%s", version, branch),
	}
	b := cfg.GetBranch(branch)
	return tagTypes[b.TagPolicy], nil
}
