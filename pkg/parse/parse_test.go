package parse

import (
	"github.com/stretchr/testify/assert"
	"go/parser"
	"go/token"
	"testing"
)

func TestBuilder_ParseFileByName(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		b := new(Builder)
		err := b.ParseFileByName("../../testdata/types.go")
		t.Log(err)
		t.Log(b)
	})

	t.Run("struct field pointer", func(t *testing.T) {
		const src = `
		package main
		
		type test struct {
			Name  *string
			Age   *uint8
			Email *string
		}`

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
		assert.Nilf(t, err, "err: %v", err)

		b := new(Builder)
		b.ParseByAstFile(file)
		const want = `{"Name":"main","Structs":[{"Name":"test","Fields":[{"Name":"Name","Kind":"*string"},{"Name":"Age","Kind":"*uint8"},{"Name":"Email","Kind":"*string"}]}]}`
		assert.Equal(t, want, b.String())
	})

}
