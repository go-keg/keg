package helpers

import "testing"

func TestSHA256(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{data: "SHA256"}, "0b0c8624f8585ea0f3cb1d8190f85536f7289776a7c318862fe10d1fcc59f5bb"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SHA256(tt.args.data); got != tt.want {
				t.Errorf("SHA256() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMD5(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{data: "MD5"}, "eba85330b2eed69d5be5cfe23376b08e"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MD5(tt.args.data); got != tt.want {
				t.Errorf("MD5() = %v, want %v", got, tt.want)
			}
		})
	}
}
