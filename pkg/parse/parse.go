package parse

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"sp/pb"
)

// this file ga ast to protocol buffer

type Builder struct {
	*pb.Pkg
}

func (b *Builder) ParseFileByName(filename string) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	b.Pkg = &pb.Pkg{
		Name:    file.Name.Name,
		Imports: parseImports(file.Imports),
	}

	for name, obj := range file.Scope.Objects {
		if obj.Kind != ast.Typ {
			continue
		}

		spec, ok := obj.Decl.(*ast.TypeSpec)
		if !ok {
			continue
		}

		ast.Inspect(spec.Type, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.StructType:
				fields := make([]*pb.Field, 0, len(x.Fields.List))
				for _, field := range x.Fields.List {
					var tags map[string]string
					if field.Tag != nil {
						tags = parseStructTags(field.Tag.Value)
					}

					fieldKind := getFieldKind(field.Type)
					fieldName := getFieldName(field.Names, fieldKind)

					fields = append(fields, &pb.Field{Name: fieldName, Kind: fieldKind, Tags: tags})
				}

				if len(fields) > 0 {
					b.Structs = append(b.Structs, &pb.Struct{Name: name, Fields: fields})
				}
				return false
			}
			return true
		})
	}
	return nil
}

func (b *Builder) String() string {
	marshal, _ := json.Marshal(b.Pkg)
	return string(marshal)
}
