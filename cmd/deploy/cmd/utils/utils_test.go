package utils

import "testing"

func TestGetBranch(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"", "dev", false},
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

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"", "v1.4.5-21-g1633e728", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetVersion()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}
