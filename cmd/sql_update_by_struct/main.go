package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/The-night-elves/sp/pkg/parse"
	"github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/packages"
	"log/slog"
	"path/filepath"
	"slices"
	"strings"
)

type Config struct {
	TypeNames []string
	Tag       string
	FuncName  string
}

func (c *Config) Parse() {
	flag.Func("type", "comma-separated list of type names; must be set", func(s string) error {
		c.TypeNames = strings.Split(s, ",")
		return nil
	})
	flag.StringVar(&c.Tag, "tag", "db", "extra struct tag name default 'db'")
	flag.StringVar(&c.FuncName, "func_name", "Update", "generate func name default 'Update'")
}

func (c *Config) Usage() {
	output := flag.CommandLine.Output()
	writer := bufio.NewWriter(output)
	_, _ = writer.WriteString("Usage of sql_update_by_struct:\n")
	_, _ = writer.WriteString("\tstringer [flags] -type 'struct_name1,struct_name2, ...'\n")
	_, _ = writer.WriteString("\tstringer [flags] -tag 'db' default 'db'\n")
	_, _ = writer.WriteString("\tstringer [flags] -func_name 'Update' \n")
	_, _ = writer.WriteString("Flags:\n")
	flag.PrintDefaults()
}

func (c *Config) String() string {
	return fmt.Sprintf("typeNames: [%s], tag: %s, funcName: %s",
		strings.Join(c.TypeNames, ","), c.Tag, c.FuncName)
}

func (c *Config) Generate(sp *parse.Builder, f *jen.File) bool {
	const agrsIdent = "args"
	const colsIdent = "cols"
	var find bool
	for _, obj := range sp.Structs {
		if !slices.Contains(c.TypeNames, obj.Name) {
			continue
		}
		find = true
		var ident = strings.ToLower(string(obj.Name[0]))
		funcHead := f.Func().
			Params(jen.Id(ident).Op("*").Id(obj.Name)).Id(c.FuncName).
			Params(jen.Id(agrsIdent).Index().Any()).
			Parens(jen.Index().String().Id(",").Index().Any())

		var body = new(jen.Statement)
		for _, field := range obj.Fields {
			name := field.Name
			if name == "" {
				name = field.Kind
			}
			dbField, ok := field.Tags[c.Tag]
			if !ok {
				dbField = strings.ToLower(name)
			}
			if field.Struct != nil {
				body.If(jen.Id(ident).Dot(name).Op("!=").Nil()).Block(
					jen.List(jen.Id("_cols"), jen.Id("_args")).Op(":=").
						Id(ident).Dot(name).Dot(c.FuncName).Params(jen.Id(agrsIdent)),

					jen.Id(agrsIdent).Op("=").Append(jen.Id(agrsIdent), jen.Id("_args...")),
					jen.Id(colsIdent).Op("=").Append(jen.Id(colsIdent), jen.Id("_cols...")),
				).Line()
			} else {
				body.If(jen.Id(ident).Dot(name).Op("!=").Nil()).Block(
					jen.Id(agrsIdent).Op("=").Append(jen.Id(agrsIdent), jen.Id("*").Id(ident).Dot(name)),

					jen.Id(colsIdent).Op("=").Append(jen.Id(colsIdent),
						jen.Id("\""+dbField).Id("=").Id(`$"+`).Qual("strconv", "Itoa").
							Call(jen.Len(jen.Id(agrsIdent)))),
				).Line()
			}
		}
		body.Return(jen.Id(colsIdent), jen.Id(agrsIdent))
		funcHead.Block(
			jen.Id(colsIdent).Op(":=").
				Make(jen.Index().String(), jen.Lit(0), jen.Lit(len(obj.Fields))),
			body,
		).Line()
	}
	return find
}

func main() {
	cmdCfg := new(Config)
	cmdCfg.Parse()
	flag.Usage = cmdCfg.Usage
	flag.Parse()

	if len(cmdCfg.TypeNames) == 0 {
		flag.Usage()
		return
	}

	var dir string
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
		dir = args[0]
	} else {
		dir = filepath.Dir(args[0])
	}

	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
		// TODO: Need to think about constants in test files. Maybe write type_string_test.go
		// in a separate pass? For later.
		Tests: false,
	}

	pkg, err := packages.Load(cfg, args...)
	if err != nil {
		slog.Error("load packages", slog.String("err", err.Error()))
		return
	} else if len(pkg) == 0 {
		slog.Error("no packages found")
		return
	}

	slog.Info("generate config", slog.String("config", cmdCfg.String()))
	sp := parse.Builder{}
	var f *jen.File

	var baseName string
	for i, file := range pkg[0].Syntax {
		sp.ParseByAstFile(file)
		if i == 0 {
			f = jen.NewFilePathName(dir, sp.Name)
			f.PackageComment("// Code generated by sql_update_by_struct DO NOT EDIT.\n")
		}
		find := cmdCfg.Generate(&sp, f)
		if baseName != "" {
			continue
		} else if find {
			// abst path /testdata/types.go
			_, fileName := filepath.Split(pkg[0].GoFiles[i])
			baseName, _ = strings.CutSuffix(fileName, ".go")
		}
	}

	// Write to file.
	outputName := filepath.Join(dir, baseName+".gen.go")
	err = f.Save(outputName)
	if err != nil {
		slog.Error("save file", slog.String("err", err.Error()))
		return
	}
	slog.Info("save file", slog.String("file", outputName))
}
