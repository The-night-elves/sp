package parse

import (
	"encoding/json"
	"github.com/The-night-elves/sp/pb"
	"go/ast"
	"go/parser"
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

	for name, obj := range file.Scope.Objects {
		if obj.Kind != ast.Typ {
			continue
		}

		spec, ok := obj.Decl.(*ast.TypeSpec)
		if !ok {
			continue
		}

		ast.Inspect(spec.Type, func(n ast.Node) bool {
			_, kind, fields := parseNode(n)
			if len(fields) == 0 {
				return false
			}

			b.Structs = append(b.Structs, &pb.Struct{
				Name:   name,
				Kind:   kind,
				Fields: fields,
			})
			return true
		})
	}
}

func (b *Builder) String() string {
	marshal, _ := json.Marshal(b.Pkg)
	return string(marshal)
}
