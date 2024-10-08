package utils

import "testing"

func TestExecDir(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"", "utils", false},
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
