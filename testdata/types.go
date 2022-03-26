package testdata

import (
	"bytes"
	"github.com/The-night-elves/sp/pb"
)

type Person struct {
	Name  string
	Age   uint8
	Email string
}

//go:generate ../cmd/sql_update_by_struct/sql_update_by_struct -type PersonUpdate

type PersonUpdate struct {
	Name  *string
	Age   *uint8
	Email *string
}

// 其他类型
type sbuf struct {
	bufSlice []byte
	bufBytes bytes.Buffer
}

type obj struct {
	*sbuf
}

type obj2 struct {
	*pb.Struct
}
