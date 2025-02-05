package gql

import (
	"testing"

	"github.com/samber/lo"
)

func TestOffsetLimit(t *testing.T) {
	type args struct {
		page *int
		size *int
		opts []PageOption
	}
	tests := []struct {
		name       string
		args       args
		wantOffset int
		wantLimit  int
	}{
		{"", args{nil, nil, nil}, 0, 10},
		{"", args{nil, nil, []PageOption{WithDefaultSize(20)}}, 0, 20},
		{"", args{nil, nil, []PageOption{WithMaxSize(100)}}, 0, 10},
		{"", args{nil, nil, []PageOption{WithMaxItems(5000)}}, 0, 10},

		{"", args{nil, lo.ToPtr(100), []PageOption{}}, 0, 100},
		{"", args{nil, lo.ToPtr(100), []PageOption{WithDefaultSize(50)}}, 0, 100},
		{"", args{nil, lo.ToPtr(100), []PageOption{WithMaxSize(50)}}, 0, 50},
		{"", args{nil, lo.ToPtr(100), []PageOption{WithMaxSize(500)}}, 0, 100},
		{"", args{nil, lo.ToPtr(100), []PageOption{WithMaxItems(50)}}, 0, 50},
		{"", args{nil, lo.ToPtr(100), []PageOption{WithMaxItems(500)}}, 0, 100},

		{"", args{lo.ToPtr(5), lo.ToPtr(100), []PageOption{}}, 400, 100},
		{"", args{lo.ToPtr(5), lo.ToPtr(100), []PageOption{WithDefaultSize(50)}}, 400, 100},
		{"", args{lo.ToPtr(5), lo.ToPtr(100), []PageOption{WithMaxSize(50)}}, 200, 50},
		{"", args{lo.ToPtr(5), lo.ToPtr(100), []PageOption{WithMaxItems(400)}}, 400, 0},
		{"", args{lo.ToPtr(5), lo.ToPtr(100), []PageOption{WithMaxItems(1000)}}, 400, 100},

		{"page zero", args{lo.ToPtr(0), lo.ToPtr(100), []PageOption{}}, 0, 100},
		{"page zero with max items", args{lo.ToPtr(0), lo.ToPtr(100), []PageOption{WithMaxItems(50)}}, 0, 50},
		{"size exceeds max size", args{lo.ToPtr(5), lo.ToPtr(200), []PageOption{WithMaxSize(100)}}, 400, 100},
		{"max items override size and page", args{lo.ToPtr(10), lo.ToPtr(100), []PageOption{WithMaxItems(50)}}, 900, 0},

		// 测试同时应用多个选项的情况
		{"combine max size and default size", args{nil, nil, []PageOption{WithDefaultSize(50), WithMaxSize(100)}}, 0, 50},
		{"combine max items with page", args{lo.ToPtr(2), lo.ToPtr(100), []PageOption{WithMaxItems(200)}}, 100, 100},

		// 测试page为零的情况
		{"page zero with size", args{lo.ToPtr(0), lo.ToPtr(50), []PageOption{}}, 0, 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOffset, gotLimit := OffsetLimit(tt.args.page, tt.args.size, tt.args.opts...)
			if gotOffset != tt.wantOffset {
				t.Errorf("OffsetLimit() gotOffset = %v, want %v", gotOffset, tt.wantOffset)
			}
			if gotLimit != tt.wantLimit {
				t.Errorf("OffsetLimit() gotLimit = %v, want %v", gotLimit, tt.wantLimit)
			}
		})
	}
}
