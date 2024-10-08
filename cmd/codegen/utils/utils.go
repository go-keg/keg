package utils

import (
	"os"
	"path/filepath"
	"text/template"
)

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
	defer outFile.Close()

	// 执行主模板 base.tmpl 并将输出写入文件
	err = temp.Execute(outFile, data)
	if err != nil {
		panic(err)
	}
	return nil
}
