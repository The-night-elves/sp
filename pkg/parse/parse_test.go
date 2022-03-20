package parse

import (
	"testing"
)

func TestBuilder_ParseFileByName(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		b := new(Builder)
		err := b.ParseFileByName("/home/night/lang/go/sp/testdata/types.go")
		t.Log(err)
		t.Log(b)
	})

}
