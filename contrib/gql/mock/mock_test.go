package mock

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/go-keg/keg/contrib/helpers"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

func TestEngine_MockQuery(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{"", args{"query {\n  user {\n    id\n    name\n    email\n    age\n     }\n}"}, nil, false},
	}

	engine, err := NewEngine(gqlparser.MustLoadSchema(&ast.Source{
		Name: "Query",
		Input: `type Query {
  user: User!
}

type User {
  id: ID!
  name: String!
  email: String!
  age: Int!
}`}))
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := engine.MockQuery(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("MockQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MockQuery() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestName(t *testing.T) {
	var count int
	var index int
	for {
		s1 := helpers.MD5(faker.Email())
		s2 := helpers.MD5(s1)
		index++
		fmt.Printf("\rindex: %d%%", index)
		if s1[0] == s2[0] {
			if s1[5] == 'y' || s2[5] == 'y' {
				fmt.Println(s1[5], s2[5], 'y')
				fmt.Println(index, s1, s2)

				count++
			}
			if count > 3 {
				break
			}
		}
	}
}
