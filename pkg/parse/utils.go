package parse

import (
	"go/ast"
	"strings"
)

// input `json:"age" db:"age"`
//  retrun { "json": "age", "db": "age" }
func parseStructTags(row string) map[string]string {
	if strings.HasPrefix(row, "`") {
		row = row[1:]
	}
	if strings.HasSuffix(row, "`") {
		row = row[:len(row)-1]
	}
	kvStr := strings.Split(row, " ")
	tags := make(map[string]string, len(kvStr))
	for i := 0; i < len(kvStr); i++ {
		key, val, _ := strings.Cut(kvStr[i], ":")
		if key == "" {
			continue
		}
		if strings.HasPrefix(val, "\"") {
			val = val[1:]
		}
		if strings.HasSuffix(val, "\"") {
			val = val[:len(val)-1]
		}
		tags[key] = val
	}

	return tags
}

func getFieldKind(raw any) (kind string) {
	switch x := raw.(type) {
	case *ast.StarExpr:
		// 如果嵌套的是 常规类型，类型为 *ast.Ident
		// 如果嵌套的是指针类型，类型为 *ast.SelectorExpr
		kind = getFieldKind(x.X)
	case *ast.SelectorExpr:
		kind = x.X.(*ast.Ident).Name + "." + x.Sel.Name
	case *ast.ArrayType:
		kind = "[]" + x.Elt.(*ast.Ident).Name
	case *ast.Ident:
		kind = x.Name
	}

	return
}

// 优先从 names 中获取，如果没有，尝试从 kind 中获取
//  kind pb.Struct
func getFieldName(names []*ast.Ident, kind string) (name string) {
	if len(names) > 0 {
		name = names[0].Name
	} else {
		_, name, _ = strings.Cut(kind, ".")
	}
	return
}

func parseImports(list []*ast.ImportSpec) (imports []string) {
	for _, item := range list {
		imports = append(imports, item.Path.Value)
	}
	return
}
