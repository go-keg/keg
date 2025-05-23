package utils

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/manifoldco/promptui"
)

func TestGetBranch(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"", "dev", false},
		{"", "main", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBranch()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBranch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBranch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

var gitVersion = regexp.MustCompile(`^v(\d+\.\d+\.\d+)-(\d+)-g([a-f0-9]+)$`)

func TestGetVersion(t *testing.T) {
	got, err := GetVersion()
	if err != nil {
		t.Errorf("GetVersion() error = %v", err)
		return
	}
	if !gitVersion.MatchString(got) {
		t.Errorf("GetVersion() got = %v", got)
	}
}

func TestExecDir(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"", "utils", false},
		{"", "utils-X", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExecDir()
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExecDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestName(t *testing.T) {
	items := []string{"Vim", "Emacs", "Sublime", "VSCode", "Atom"}
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label:    "What's your text editor",
			Items:    items,
			AddLabel: "Other",
		}

		index, result, err = prompt.Run()

		if index == -1 {
			items = append(items, result)
		}
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %s\n", result)
}

func TestProjectRootPath(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"", "/Users/eiixy/workspace/eiixy/go-keg/keg", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProjectRootPath()
			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectRootPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ProjectRootPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}
