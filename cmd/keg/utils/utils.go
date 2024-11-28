package utils

import (
	"bufio"
	"fmt"
	"github.com/go-keg/keg/cmd/keg/internal"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
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

func GoModuleName(path string) (string, error) {
	goModulePath := filepath.Join(path, "go.mod")
	_, err := os.Stat(goModulePath)
	if os.IsNotExist(err) {
		return "", nil
	}
	file, err := os.Open(goModulePath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1], nil
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("module name not found in go.mod")
}

func ProjectRootPath() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := currentDir
	for {
		_, err = os.Stat(filepath.Join(dir, internal.ConfigFile))
		if err == nil {
			return dir, nil
		}
		parentDir := filepath.Dir(dir)
		if parentDir == dir { // 到达根目录时停止
			return "", fmt.Errorf("file not found: %s", internal.ConfigFile)
		}
		dir = parentDir
	}
}
