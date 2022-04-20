package parse

import (
	"bytes"
	"encoding/json"
	"github.com/The-night-elves/sp/pb"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
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

	b.ParseByAstFile(file)

	return nil
}

func (b *Builder) ParseByAstFile(file *ast.File) {
	b.Pkg = &pb.Pkg{
		Name:    file.Name.Name,
		Imports: parseImports(file.Imports),
	}

	fset := token.NewFileSet()
	buf := new(bytes.Buffer)

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

					var fieldKind string
					if err := printer.Fprint(buf, fset, field.Type); err == nil {
						fieldKind = buf.String()
						buf.Reset()
					} else {
						fieldKind = getFieldKind(field.Type)
					}

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
}

func (b *Builder) String() string {
	marshal, _ := json.Marshal(b.Pkg)
	return string(marshal)
}
