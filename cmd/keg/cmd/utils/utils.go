package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/go-keg/keg/cmd/keg/cmd/config"
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
	return strings.TrimPrefix(strings.ReplaceAll(string(output), "\n", ""), "refs/tags/"), nil
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
		config.TagPolicyVersion:       version,
		config.TagPolicyBranch:        fmt.Sprintf("latest-%s", branch),
		config.TagPolicyVersionBranch: fmt.Sprintf("%s-%s", version, branch),
	}
	b := cfg.GetBranch(branch)
	return tagTypes[b.TagPolicy], nil
}

func ExecDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Base(cwd), err
}

func WriteFileWithName(temp *template.Template, data any, path string, tempName string) error {
	outFile, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer func(outFile *os.File) {
		_ = outFile.Close()
	}(outFile)

	// 执行主模板 base.tmpl 并将输出写入文件
	err = temp.ExecuteTemplate(outFile, tempName, data)
	if err != nil {
		panic(err)
	}
	return nil
}

func WriteFile(temp *template.Template, data any, path string) error {
	outFile, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer func() { _ = outFile.Close() }()

	// 执行主模板 base.tmpl 并将输出写入文件
	err = temp.Execute(outFile, data)
	if err != nil {
		panic(err)
	}
	return nil
}
