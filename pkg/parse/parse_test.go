package parse

import (
	"github.com/stretchr/testify/require"
	"go/parser"
	"go/token"
	"sort"
	"testing"
)

func TestBuilder_ParseFileByName(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		b := new(Builder)
		err := b.ParseFileByName("../../testdata/types.go")
		t.Log(err)
		t.Log(b)
	})
	tests := []struct {
		name string
		src  string
		want string
	}{
		{
			name: "struct field pointer",
			src: `
		package main
		
		type test struct {
			Name  *string
			Age   *uint8
			Email *string
		}`,
			want: `{"Name":"main","Structs":[{"Name":"test","Fields":[{"Name":"Name","Kind":"*string"},{"Name":"Age","Kind":"*uint8"},{"Name":"Email","Kind":"*string"}]}]}`,
		},
		{
			name: "type alias",
			src: `
		package main
		
		type Age uint8

		type Person struct {
			Name *string 
			Age  Age    
		}`,
			want: `{"Name":"main","Structs":[{"Name":"Person","Fields":[{"Name":"Name","Kind":"*string"},{"Name":"Age","Kind":"Age","struct":{"Name":"Age","Kind":"Age","Fields":[{"Name":"Age","Kind":"uint8"}]}}]}]}`,
		},
		{
			name: "anonymous pointer",
			src: `
		package main

		type sbuf struct {
			bufSlice []byte
			bufBytes bytes.Buffer
		}
		
		type obj struct {
			*sbuf
		}`,
			want: `{"Name":"main","Structs":[{"Name":"obj","Fields":[{"Name":"sbuf","Kind":"*sbuf","struct":{"Name":"sbuf","Kind":"*sbuf","Fields":[{"Name":"sbuf"}]}}]},{"Name":"sbuf","Fields":[{"Name":"bufSlice","Kind":"[]byte"},{"Name":"bufBytes","Kind":"bytes.Buffer"}]}]}`,
		},

		{
			name: "struct type alias pointer",
			src: `
		package main

		type sbuf struct {
			bufSlice []byte
			bufBytes bytes.Buffer
		}
		
		type obj struct {
			buf *sbuf
		}`,
			want: `{"Name":"main","Structs":[{"Name":"obj","Fields":[{"Name":"buf","Kind":"*sbuf","struct":{"Name":"sbuf","Kind":"*sbuf","Fields":[{"Name":"sbuf"}]}}]},{"Name":"sbuf","Fields":[{"Name":"bufSlice","Kind":"[]byte"},{"Name":"bufBytes","Kind":"bytes.Buffer"}]}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", tt.src, parser.ParseComments)
			require.Nil(t, err)
			b := new(Builder)
			b.ParseByAstFile(file)
			if len(b.Structs) > 1 {
				sort.Slice(b.Structs, func(i, j int) bool { return b.Structs[i].Name < b.Structs[j].Name })
			}
			require.Equal(t, tt.want, b.String())
		})
	}
}
