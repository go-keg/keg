package http

import (
	"context"
	"reflect"
	"testing"
)

func TestClient_Get(t *testing.T) {
	type fields struct {
		opts []ClientOptionFunc
	}
	type args struct {
		url  string
		opts []OptionFunc
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Response
		wantErr bool
	}{
		{"", fields{opts: []ClientOptionFunc{WithBaseURL("https://api.example.com/account")}}, args{
			url: "/users/{id}/info",
			opts: []OptionFunc{
				SetPathParams(map[string]string{
					"id": "123",
				}),
				SetQueryParams(map[string]string{
					"enable": "true",
					"rf":     "abc-def",
				}),
			},
		}, nil, false},
		{"", fields{opts: []ClientOptionFunc{WithBaseURL("https://api.example.com/account/")}}, args{
			url: "/users/{id}/info",
			opts: []OptionFunc{
				SetPathParams(map[string]string{
					"id": "123",
				}),
			},
		}, nil, false},
		{"", fields{opts: []ClientOptionFunc{WithBaseURL("https://api.example.com/account/")}}, args{
			url: "https://api.example-cover.com/users/{id}/info",
			opts: []OptionFunc{
				SetPathParams(map[string]string{
					"id": "123",
				}),
			},
		}, nil, false},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.fields.opts...)
			got, err := c.Get(ctx, tt.args.url, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Graphql(t *testing.T) {
	type fields struct {
		opts []ClientOptionFunc
	}
	var ctx = context.Background()
	type args struct {
		endpoint  string
		query     string
		variables map[string]any
		opts      []OptionFunc
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Response
		wantErr bool
	}{
		{"", fields{opts: nil}, args{
			endpoint: "https://v2-api-qa.tbanx.cn/resource/query",
			query:    "query ($first: Int){\n  emails(first:$first){\n    nodes{\n      id\n      sender\n      subject\n      bodyHTML\n    }\n  }\n}",
			variables: map[string]any{
				"first": 10,
			},
			opts: nil,
		}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.fields.opts...)
			got, err := c.Graphql(ctx, tt.args.endpoint, tt.args.query, tt.args.variables, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Graphql() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Graphql() got = %v, want %v", got, tt.want)
			}
		})
	}
}
