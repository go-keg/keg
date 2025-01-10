# Keg

[![Go Report Card](https://goreportcard.com/badge/github.com/go-keg/keg)](https://goreportcard.com/report/github.com/go-keg/keg)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-keg/keg.svg)](https://pkg.go.dev/github.com/go-keg/keg)
[![License](https://img.shields.io/github/license/go-keg/keg)](./LICENSE)

> A Tool Library

## 🚀 QuickStart

```bash
go get github.com/go-keg/keg
```

#### Install
```bash
go install github.com/go-keg/keg/cmd/keg@latest
```

## 🎯 Features

### cmd/keg
```shell
keg image tag #根据 keg.yaml 配置的分支策略生成 build tag
keg k8s deployment update-image -n $(SERVICE) #更新 deployment 镜像版本(重新部署)
keg k8s gen config -n ${NAMESPACE} # 根据模板生成对应namespace的k8s部署相关文件
```

### contrib/ent
* annotations 
  * EnumsGQLSchemaHook 将枚举值转换为GraphQL枚举类型，增加枚举值说明
* mixin
  * TimeMixin 为所有实体添加创建时间和更新时间字段
  * SoftDeleteMixin 为所有实体添加删除时间字段
* template 增加事务代码及软删除相关代码生成

### contrib/gql
* directive
  * @cache 为GraphQL字段增加Cache
  * @hashError 对未定义的错误信息进行Hash处理
  * @validate 对输入数据进行正则验证
* dataloader 优化N+1查询
```go
func (r AdminLoader) userRoleCount() gql.LoaderFunc {
	type item struct {
		ID    int64 `json:"id"`
		Count int   `json:"count"`
	}
	return func(ctx context.Context, keys dataloader.Keys) (map[dataloader.Key]any, error) {
		var items []item
		err := r.client.User.Query().Where(user.IDIn(gql.ToInts(keys)...)).Modify(func(s *sql.Selector) {
			t1 := sql.Table("user_roles").As("t1")
			s.Select(s.C(user.FieldID), sql.As(sql.Count(t1.C("role_id")), "count")).
				From(s).LeftJoin(t1).On(s.C(user.FieldID), t1.C("user_id")).
				GroupBy(s.C(user.FieldID))
		}).Scan(ctx, &items)
		if err != nil {
			return nil, err
		}
		return lo.SliceToMap(items, func(item item) (dataloader.Key, any) {
			return gql.ToStringKey(item.ID), item.Count
		}), nil
	}
}
```

## TODO
