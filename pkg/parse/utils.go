package parse

import (
	"github.com/The-night-elves/sp/pb"
	"go/ast"
	"strings"
)

// input `json:"age" db:"age"`
//
//	retrun { "json": "age", "db": "age" }
func parseStructTags(row string) map[string]string {
	row = strings.Trim(row, "`")
	kvStr := strings.Split(row, " ")
	tags := make(map[string]string, len(kvStr))
	for i := 0; i < len(kvStr); i++ {
		key, val, _ := strings.Cut(kvStr[i], ":")
		if key == "" {
			continue
		}
		tags[key] = strings.Trim(val, "\"")
	}

	return tags
}

func parseImports(list []*ast.ImportSpec) (imports []string) {
	for _, item := range list {
		imports = append(imports, item.Path.Value)
	}
	return
}

// input struct field
func parseNode(n any) (name, kind string, fields []*pb.Field) {
	switch v := n.(type) {
	case *ast.StructType:
		fields = make([]*pb.Field, 0, len(v.Fields.List))
		for _, field := range v.Fields.List {
			var tags map[string]string
			if field.Tag != nil {
				tags = parseStructTags(field.Tag.Value)
			}

			subName, subKind, subFields := parseNode(field.Type)
			p := &pb.Field{Name: "", Kind: subKind, Tags: tags, Struct: nil}
			if len(field.Names) >= 1 {
				p.Name = field.Names[0].Name
			} else {
				p.Name = subName
			}
			if len(subFields) > 0 {
				p.Struct = &pb.Struct{Name: subName, Kind: subKind, Fields: subFields}
			}
			fields = append(fields, p)
		}
		return
	case *ast.StarExpr:
		name, kind, fields = parseNode(v.X)
		if kind != "" {
			kind = "*" + kind
		}
	case *ast.Ident:
		if v.Obj != nil {
			if v.Obj.Decl != nil {
				_, _, fields = parseNode(v.Obj.Decl)
			}
			name = v.Name
			if v.Obj.Kind == ast.Typ {
				kind = v.Obj.Name
			}
		} else {
			kind = v.Name
		}
	case *ast.TypeSpec:
		_, subKind, _ := parseNode(v.Type)
		fields = append(fields, &pb.Field{Name: v.Name.Name, Kind: subKind, Tags: nil, Struct: nil})
	case *ast.ArrayType:
		if v.Elt != nil {
			_, kind, _ = parseNode(v.Elt)
			kind = "[]" + kind
		}
	case *ast.SelectorExpr:
		_, kind, _ = parseNode(v.X)
		_, k2, _ := parseNode(v.Sel)
		if kind != "" && k2 != "" {
			kind = kind + "." + k2
		}
	}
	return
}
