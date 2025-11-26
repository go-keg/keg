package helpers

import "testing"

func TestFormatNumber(t *testing.T) {
	type args struct {
		v int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{1}, "1"},
		{"", args{9}, "9"},
		{"", args{999}, "999"},
		{"", args{1000}, "1.00K"},
		{"", args{1005}, "1.00K"},
		{"", args{1006}, "1.01K"},
		{"", args{1005000}, "1.00M"},
		{"", args{1006000}, "1.01M"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatNumber(tt.args.v); got != tt.want {
				t.Errorf("FormatNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatNumberRange(t *testing.T) {
	type args struct {
		v int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{0}, "0-1K"},
		{"", args{1000}, "0-1K"},
		{"", args{1001}, "1K-5K"},
		{"", args{10000001}, "10M+"},
		{"", args{1000000001}, "10M+"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatNumberRange(tt.args.v); got != tt.want {
				t.Errorf("FormatNumberRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	type args struct {
		x int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"", args{12345678}, 87654321},
		{"", args{3141592600}, 62951413},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reverse(tt.args.x); got != tt.want {
				t.Errorf("Reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFib(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"", args{8}, 21},
		{"", args{16}, 987},
		{"", args{32}, 2178309},
		{"", args{64}, 10610209857723},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Fib(tt.args.n); got != tt.want {
				t.Errorf("Fib() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoundToDecimal(t *testing.T) {
	type args struct {
		value     float64
		precision int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"",args{3.1415926,0},3},
		{"",args{3.1415926,1},3.1},
		{"",args{3.1415926,2},3.14},
		{"",args{3.1415926,3},3.142},
		{"",args{3.1415926,4},3.1416},
		{"",args{3.1415926,5},3.14159},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RoundToDecimal(tt.args.value, tt.args.precision); got != tt.want {
				t.Errorf("RoundToDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}