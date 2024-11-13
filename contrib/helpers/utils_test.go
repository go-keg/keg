package helpers

import "testing"

func TestHashWithSHA256(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{data: "HashWithSHA256"}, "0b0c8624f8585ea0f3cb1d8190f85536f7289776a7c318862fe10d1fcc59f5bb"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashWithSHA256(tt.args.data); got != tt.want {
				t.Errorf("HashWithSHA256() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashWithMD5(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{data: "HashWithMD5"}, "eba85330b2eed69d5be5cfe23376b08e"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashWithMD5(tt.args.data); got != tt.want {
				t.Errorf("HashWithMD5() = %v, want %v", got, tt.want)
			}
		})
	}
}
